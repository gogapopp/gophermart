package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogapopp/gophermart/internal/app/config"
	"github.com/gogapopp/gophermart/internal/app/handler"
	"github.com/gogapopp/gophermart/internal/app/logger"
	"github.com/gogapopp/gophermart/internal/app/server"
	"github.com/gogapopp/gophermart/internal/app/service"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// парсим env и флаги
	config := config.ParseConfig()
	// инициализируем конфиг
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	// устанавливаем подключение к бд
	db, err := storage.NewDB(ctx, config.DatabaseURI)
	if err != nil {
		log.Fatal("error initialize database", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("error closing database", err)
		}
	}()

	storage := storage.NewStorage(ctx, db)
	services := service.NewService(storage)
	handlers := handler.NewHandler(config, services, log)

	srv := server.NewServer(config, log)
	go func() {
		if err := srv.Run(handlers.InitRoutes()); err != nil {
			log.Fatalln("error to start the server", err)
		}
	}()
	// реализация graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("error shutting down server", err)
	}
}
