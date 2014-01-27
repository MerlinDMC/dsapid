package main

import (
	"flag"
	"fmt"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/converter/dsapi"
	"github.com/MerlinDMC/dsapid/converter/imgapi"
	"github.com/MerlinDMC/dsapid/server/logger"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/server/sync"
	"github.com/MerlinDMC/dsapid/storage"
	"github.com/codegangsta/martini"
	"net/http"
	"os"
	"runtime"
)

var (
	flagVersion      bool
	flagConfigFile   string
	flagMaxCpu       int
	flagMaxFetches   int
	flagLogDebug     bool
	flagLogDebugMore bool
	flagPrettifyJson bool
)

func init() {
	flag.BoolVar(&flagVersion, "V", false, "display version information and exit")
	flag.StringVar(&flagConfigFile, "config", "data/config.json", "configuration file")
	flag.IntVar(&flagMaxCpu, "max_cpu", 2, "number of processors to use")
	flag.IntVar(&flagMaxFetches, "max_fetches", 2, "number of parallel sync fetches")
	flag.BoolVar(&flagLogDebug, "debug", false, "display additional debug information")
	flag.BoolVar(&flagLogDebugMore, "debug_more", false, "display extra information while running subcommands")
	flag.BoolVar(&flagPrettifyJson, "prettify", false, "prettify json output")

	logger.SetName("dsapid")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if flagVersion {
		fmt.Printf("%s v%s\n", dsapid.AppName, dsapid.AppVersion)
		os.Exit(0)
	}

	var config Config = DefaultConfig()

	if err := config.Load(flagConfigFile); err != nil {
		os.Exit(1)
	}

	runtime.GOMAXPROCS(flagMaxCpu)

	router := martini.NewRouter()

	registerRoutes(router)

	handler := martini.New()
	handler.Action(router.Handle)

	user_storage := storage.NewUserStorage(config.UsersConfig)
	manifest_storage := storage.NewManifestStorage(config.DataDir)

	sync_manager := sync.NewManager(flagMaxFetches, user_storage, manifest_storage)
	sync_manager.Init()

	handler.MapTo(user_storage, (*storage.UserStorage)(nil))
	handler.MapTo(manifest_storage, (*storage.ManifestStorage)(nil))
	handler.MapTo(sync_manager, (*sync.SyncManager)(nil))

	handler.MapTo(dsapi.NewEncoder(config.BaseUrl, user_storage), (*converter.DsapiManifestEncoder)(nil))
	handler.MapTo(imgapi.NewEncoder(config.BaseUrl, user_storage), (*converter.ImgapiManifestEncoder)(nil))

	handler.Use(middleware.EncodeOutput(flagPrettifyJson))
	handler.Use(middleware.Auth(user_storage))

	if flagLogDebug {
		logger.SetLevel(logger.DEBUG)

		handler.Use(middleware.JsonLogger())
	}

	if config.MountUi != "" {
		router.Any("/", func(res http.ResponseWriter, req *http.Request) {
			http.Redirect(res, req, "/ui/", http.StatusMovedPermanently)
		})

		handler.Use(martini.Static(config.MountUi, martini.StaticOptions{
			Prefix:      "/ui",
			IndexFile:   "index.html",
			SkipLogging: true,
		}))
	}

	for _, source := range config.SyncSources {
		if source.Active {
			sync_manager.NewSyncer(source)
		}
	}

	if err := sync_manager.Run(); err != nil {
		logger.Fatalf("error starting sync manager: %s", err)

		os.Exit(2)
	}

	if config.Https.ListenAddress != "" {
		go http.ListenAndServeTLS(config.Https.ListenAddress, config.Https.Cert, config.Https.Key, handler)
	}

	http.ListenAndServe(config.Http.ListenAddress, handler)
}
