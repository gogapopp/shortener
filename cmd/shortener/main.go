package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/logger"
	"github.com/gogapopp/shortener/internal/app/routes"
	"github.com/gogapopp/shortener/internal/app/storage"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	// инициализируем логер
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(err)
	}
	// парсим флаги
	config.InitializeServerConfig()
	// пытается подключится к базе данных, если не получается то пытаеимся записывать в файл, если не получается то в память
	if err := config.SetupDatabaseAndFilemanager(ctx); err != nil {
		log.Println("ошибка: ", err)
		os.Exit(1)
	}
	defer storage.DB().Close()
	// запускаем сервер
	logger.Log.Info("Running the server at", zap.String("addres", config.RunAddr))
	r := routes.Routes()
	log.Fatal(http.ListenAndServe(config.RunAddr, r))
	go func() {
		r := routes.Routes()
		log.Fatal(http.ListenAndServe(config.RunAddr, r))
	}()
}
