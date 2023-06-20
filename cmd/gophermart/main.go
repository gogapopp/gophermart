package main

import (
	"context"

	"github.com/gogapopp/gophermart/cmd/gophermart/server"
	"github.com/gogapopp/gophermart/config"
	"github.com/gogapopp/gophermart/internal/app/handler"
	"github.com/gogapopp/gophermart/internal/app/logger"
	"github.com/gogapopp/gophermart/internal/app/service"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

func main() {
	ctx := context.Background()
	// парсим env и флаги
	config.ParseConfig()
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
	defer db.Close()

	storage := storage.NewStorage(ctx, db)
	services := service.NewService(storage)
	handlers := handler.NewHandler(services, log)

	srv := server.NewServer(log)
	if err := srv.Run(handlers.InitRoutes()); err != nil {
		log.Fatal("error to start the server", err)
	}
}
