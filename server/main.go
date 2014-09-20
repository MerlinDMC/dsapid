package main

import (
	"flag"
	"fmt"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/converter/dsapi"
	"github.com/MerlinDMC/dsapid/converter/imgapi"
	"github.com/MerlinDMC/dsapid/server/middleware"
	dsapid_sync "github.com/MerlinDMC/dsapid/server/sync"
	"github.com/MerlinDMC/dsapid/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
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
	flag.StringVar(&flagLogLevel, "log_level", "error", "log level for console logs [debug,info,warn,error,fatal,panic]")
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

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	var config Config = DefaultConfig()

	if err := config.Load(flagConfigFile); err != nil {
		os.Exit(1)
	}

	switch config.LogLevel {
	case "trace", "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	}

	runtime.GOMAXPROCS(flagMaxCpu)

	router := martini.NewRouter()

	registerRoutes(router)

	handler := martini.New()
	handler.Action(router.Handle)

	log.WithFields(log.Fields{
		"config": config.UsersConfig,
	}).Debug("loading users")
	user_storage := storage.NewUserStorage(config.UsersConfig)

	log.WithFields(log.Fields{
		"directory": config.DataDir,
	}).Debug("loading datasets")
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

	switch config.LogLevel {
	case "trace", "debug", "info":
		handler.Use(middleware.LogrusLogger())
		break
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
		log.Fatalf("error starting sync manager: %s", err)

		os.Exit(2)
	}

	var wg sync.WaitGroup

	for server_name, server_config := range config.Listen {
		wg.Add(1)

		log.Infof("starting server %s...", server_name)

		go startServer(server_config, handler, &wg)
	}

	log.Infof("server running - waiting for shutdown")

	wg.Wait()
}

func startServer(config protoConfig, handler http.Handler, wg *sync.WaitGroup) (err error) {
	if config.UseSSL {
		log.Debugf("starting with ssl enabled at address %s", config.ListenAddress)

		err = http.ListenAndServeTLS(config.ListenAddress, config.Cert, config.Key, handler)
	} else {
		log.Debugf("starting at address %s", config.ListenAddress)

		err = http.ListenAndServe(config.ListenAddress, handler)
	}

	if err != nil {
		log.Errorf("server terminated: %s", err)
	}

	wg.Done()

	return err
}
