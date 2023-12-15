package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/lesnoi-kot/karten-backend/src/api"
	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func main() {
	logger := prepareLogger(os.Getenv("DEBUG") == "true")
	defer logger.Sync()

	if err := env.Parse(&settings.AppConfig); err != nil {
		logger.Fatalw("Config parsing error", "error", err)
	}

	printConfigs(logger)

	fileStorage, err := filestorage.NewFileSystemStorage(settings.AppConfig.FileStoragePath)
	if err != nil {
		logger.Fatalw("FileSystemStorage initialization error", "error", err)
	}

	storeService := store.NewStore(store.StoreConfig{
		DSN:         settings.AppConfig.StoreDSN,
		FileStorage: fileStorage,
		Logger:      logger,
		Debug:       settings.AppConfig.Debug,
	})
	if err = storeService.Ping(); err != nil {
		logger.Fatalw("DB connection error", "error", err)
	}

	apiService := api.NewAPI(api.APIConfig{
		Store:        storeService,
		Logger:       logger,
		FileStorage:  fileStorage,
		APIPrefix:    settings.AppConfig.APIPrefix,
		CookieDomain: settings.AppConfig.CookieDomain,
		AllowOrigins: settings.AppConfig.AllowOrigins,
		Debug:        settings.AppConfig.Debug,
	})

	go handleSignals(apiService)

	if err := apiService.Start(settings.AppConfig.APIBindAddress); err != nil {
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

func printConfigs(logger *zap.SugaredLogger) {
	logger.Infow("Karten config",
		"Debug", settings.AppConfig.Debug,
		"APIBindAddress", settings.AppConfig.APIBindAddress,
		"APIPrefix", settings.AppConfig.APIPrefix,
		"FrontendURL", settings.AppConfig.FrontendURL,
		"FileStoragePath", settings.AppConfig.FileStoragePath,
		"AllowOrigins", strings.Join(settings.AppConfig.AllowOrigins, ", "),
	)
}
