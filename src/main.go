package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/lesnoi-kot/karten-backend/src/api"
	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type AppConfig struct {
	StoreDSN        string   `env:"STORE_DSN,notEmpty,unset"`
	APIBindAddress  string   `env:"API_HOST,notEmpty"`
	APIPrefix       string   `env:"API_PREFIX"`
	CookieDomain    string   `env:"COOKIE_DOMAIN,notEmpty"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH,notEmpty"`
	AllowOrigins    []string `env:"ALLOW_ORIGINS,notEmpty" envSeparator:","`
	Debug           bool     `env:"DEBUG"`
}

func main() {
	logger := prepareLogger(os.Getenv("DEBUG") != "")
	defer logger.Sync()

	var cfg AppConfig
	if err := env.Parse(&cfg); err != nil {
		logger.Fatalw("Config parsing error", "error", err)
	}

	fileStorage, err := filestorage.NewFileSystemStorage(cfg.FileStoragePath)
	if err != nil {
		logger.Fatalw("FileSystemStorage initialization error", "error", err)
	}

	storeService, err := store.NewStore(store.StoreConfig{
		DSN:         cfg.StoreDSN,
		FileStorage: fileStorage,
		Logger:      logger,
		Debug:       cfg.Debug,
	})
	if err != nil {
		logger.Fatalw("DB connection error", "error", err)
	}

	apiService := api.NewAPI(api.APIConfig{
		Store:        storeService,
		Logger:       logger,
		FileStorage:  fileStorage,
		APIPrefix:    cfg.APIPrefix,
		CookieDomain: cfg.CookieDomain,
		AllowOrigins: cfg.AllowOrigins,
		Debug:        cfg.Debug,
	})

	go handleSignals(apiService)

	if err := apiService.Start(cfg.APIBindAddress); err != nil {
		logger.Info("API service is stopped")

		if err := storeService.Close(); err != nil {
			logger.Errorw("Store connection close error", "error", err)
		} else {
			logger.Info("Store connection is closed")
		}

		if err != http.ErrServerClosed {
			logger.Errorw("Server stopped with an error", "error", err)
		}
	}
}

func prepareLogger(debug bool) *zap.SugaredLogger {
	var logger *zap.Logger

	if debug {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	zap.RedirectStdLog(logger)
	return logger.Sugar()
}

func handleSignals(apiService *api.APIService) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	apiService.Shutdown()
}
