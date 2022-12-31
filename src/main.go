package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/lesnoi-kot/karten-backend/src/api"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type AppConfig struct {
	StoreDSN       string `env:"STORE_DSN,notEmpty,unset"`
	APIBindAddress string `env:"API_HOST,notEmpty"`
}

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	var cfg AppConfig

	if err := env.Parse(&cfg); err != nil {
		logger.Fatal(err.Error())
	}

	storeService, err := store.NewDB(cfg.StoreDSN)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	apiService := api.NewAPI(api.APIConfig{
		DB:     storeService,
		Logger: logger,
	})

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)

		<-quit
		logger.Info("Shutting down the server...")
		apiService.Shutdown()
	}()

	if err := apiService.Start(cfg.APIBindAddress); err != nil {
		if err := storeService.Close(); err != nil {
			logger.Info("Store connections close error")
		} else {
			logger.Info("Store connections closed")
		}

		if err != http.ErrServerClosed {
			logger.Error(err.Error())
		}
	}
}
