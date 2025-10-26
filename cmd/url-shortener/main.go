package main

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/api"
	logging "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
	"url-shortener/pkg/lib/logger/sl"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatal("Ошибка загрузки env файла:", err)
	}
	//cfg := config.MustLoad()
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	newLog, err := logging.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}

	repository, err := storage.NewRepository(context.Background(), cfg.Repository)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	serviceInstance := service.New(repository, newLog)

	app := api.NewRouters(&api.Routers{Service: serviceInstance}, cfg.Rest.Token)

	go func() {
		newLog.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	newLog.Info("Shutting down gracefully...")
}
