package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lesnoi-kot/karten-backend/src/api"
	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func main() {
	if os.Getenv("USE_DOTENV") == "true" {
		godotenv.Load()
	}

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
		DSN:    settings.AppConfig.StoreDSN,
		Logger: logger,
		Debug:  settings.AppConfig.Debug,
	})
	if err = storeService.Ping(); err != nil {
		logger.Fatalw("DB connection error", "error", err)
	}

	defer func() {
		if err := storeService.Close(); err != nil {
			logger.Errorw("DB connection close error", "error", err)
		} else {
			logger.Info("DB connection is closed")
		}
	}()

	apiService := api.NewAPI(api.APIConfig{
		Store:       storeService,
		Logger:      logger,
		FileStorage: fileStorage,
		ContextsContainer: entityservices.ContextsContainer{
			Store:       storeService,
			FileStorage: fileStorage,
		},
		APIPrefix:    settings.AppConfig.APIPrefix,
		CookieDomain: settings.AppConfig.CookieDomain,
		AllowOrigins: settings.AppConfig.AllowOrigins,
		Debug:        settings.AppConfig.Debug,
	})

	go handleSignals(apiService)

	if err := apiService.Start(settings.AppConfig.APIBindAddress); err != nil {
		logger.Info("API service is stopped")

		if errors.Is(err, http.ErrServerClosed) == false {
			logger.Errorw("Server stopped with an error", "error", err)
		}
	}
}

func prepareLogger(debug bool) *zap.SugaredLogger {
	var logger *zap.Logger

	if debug {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger = zap.Must(cfg.Build())
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
