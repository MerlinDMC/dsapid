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
	dsapid_sync "github.com/MerlinDMC/dsapid/server/sync"
	"github.com/MerlinDMC/dsapid/storage"
	"github.com/codegangsta/martini"
	"net/http"
	"os"
	"runtime"
	"sync"
)

var (
	flagVersion      bool
	flagConfigFile   string
	flagMaxCpu       int
	flagMaxFetches   int
	flagLogLevel     string
	flagPrettifyJson bool
)

func init() {
	flag.BoolVar(&flagVersion, "V", false, "display version information and exit")
	flag.StringVar(&flagConfigFile, "config", "data/config.json", "configuration file")
	flag.IntVar(&flagMaxCpu, "max_cpu", 2, "number of processors to use")
	flag.IntVar(&flagMaxFetches, "max_fetches", 1, "number of parallel sync fetches")
	flag.StringVar(&flagLogLevel, "log_level", "error", "log level for console logs [trace,debug,info,warn,error,fatal]")
	flag.BoolVar(&flagPrettifyJson, "prettify", false, "prettify json output")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if flagVersion {
		fmt.Printf("%s v%s\n", dsapid.AppName, dsapid.AppVersion)
		os.Exit(0)
	}

	logger.SetName("dsapid")

	switch flagLogLevel {
	case "trace":
		logger.SetLevel(logger.TRACE)
		break
	case "debug":
		logger.SetLevel(logger.DEBUG)
		break
	case "info":
		logger.SetLevel(logger.INFO)
		break
	case "warn":
		logger.SetLevel(logger.WARN)
		break
	case "error":
		logger.SetLevel(logger.ERROR)
		break
	case "fatal":
		logger.SetLevel(logger.FATAL)
		break
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

	logger.Debugf("loading users from %s", config.UsersConfig)
	user_storage := storage.NewUserStorage(config.UsersConfig)

	logger.Debugf("loading datasets from %s", config.DataDir)
	manifest_storage := storage.NewManifestStorage(config.DataDir)

	sync_manager := dsapid_sync.NewManager(flagMaxFetches, user_storage, manifest_storage)
	sync_manager.Init()

	handler.MapTo(user_storage, (*storage.UserStorage)(nil))
	handler.MapTo(manifest_storage, (*storage.ManifestStorage)(nil))
	handler.MapTo(sync_manager, (*dsapid_sync.SyncManager)(nil))

	handler.MapTo(dsapi.NewEncoder(config.BaseUrl, user_storage), (*converter.DsapiManifestEncoder)(nil))
	handler.MapTo(imgapi.NewEncoder(config.BaseUrl, user_storage), (*converter.ImgapiManifestEncoder)(nil))

	handler.Use(middleware.EncodeOutput(flagPrettifyJson))
	handler.Use(middleware.Auth(user_storage))

	if logger.GetLevel() <= logger.INFO {
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

	var wg sync.WaitGroup

	for server_name, server_config := range config.Listen {
		wg.Add(1)

		logger.Infof("starting server %s...", server_name)

		go startServer(server_config, handler, &wg)
	}

	logger.Infof("server running - waiting for shutdown")

	wg.Wait()
}

func startServer(config protoConfig, handler http.Handler, wg *sync.WaitGroup) (err error) {
	if config.UseSSL {
		logger.Debugf("starting with ssl enabled at address %s", config.ListenAddress)

		err = http.ListenAndServeTLS(config.ListenAddress, config.Cert, config.Key, handler)
	} else {
		logger.Debugf("starting at address %s", config.ListenAddress)

		err = http.ListenAndServe(config.ListenAddress, handler)
	}

	if err != nil {
		logger.Errorf("server terminated: %s", err)
	}

	wg.Done()

	return err
}
