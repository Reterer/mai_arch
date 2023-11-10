package main

import (
	"delivery_system/client_service/config"
	"delivery_system/client_service/internal/api"
	"delivery_system/client_service/internal/repository"
	"delivery_system/client_service/internal/service"
	"delivery_system/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const path = "config/config.json"

func main() {
	// Инициализация логгера
	logger.Initialization()
	defer logger.Close()

	// Парсинг конфигов
	cfg, err := config.ParseConfig(path)
	if err != nil {
		zap.L().Panic("parse config", zap.Error(err))
	}

	// Репозиторий
	repo := repository.New(&cfg.Repository)

	// Бизнес-логика
	serv := service.New(&cfg.Service, repo)

	// Api
	public := api.New(&cfg.Api, serv)
	publicServer := &fasthttp.Server{
		Handler: public.Handler,
	}

	go func() {
		if err := publicServer.ListenAndServe(cfg.Port); err != nil {
			zap.L().Panic("failed public server", zap.Error(err))
		}
	}()

	zap.L().Info("client_service successfully launched")

	// Ожидаем прерываний
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	if err := publicServer.Shutdown(); err != nil {
		zap.L().Error("failed shutdown public server", zap.Error(err))
	}

	zap.L().Info("shutdown")
}
