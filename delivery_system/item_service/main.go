package main

import (
	"context"
	"delivery_system/item_service/config"
	"delivery_system/item_service/internal/api"
	"delivery_system/item_service/internal/repository"
	"delivery_system/item_service/internal/service"
	"delivery_system/pkg/common_models"
	"delivery_system/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	itemRepo := repository.NewItemMemory(&cfg.Repository, nil)
	userRepo := repository.NewUserMock(&cfg.Repository, []common_models.User{
		{UserID: 42, Usernmame: "test"},
	})

	// Бизнес-логика
	serv := service.New(&cfg.Service, itemRepo, userRepo)

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
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

	if err := publicServer.ShutdownWithContext(ctx); err != nil {
		zap.L().Error("failed shutdown public server", zap.Error(err))
	}

	cancel()
	zap.L().Info("shutdown")
}
