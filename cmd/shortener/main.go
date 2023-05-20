package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/logger"
	"github.com/gogapopp/shortener/internal/app/routes"
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
	config.SetupDatabaseAndFilemanager(ctx)
	// запускаем сервер
	logger.Log.Info("Running the server at", zap.String("addres", config.RunAddr))
	r := routes.Routes()
	log.Fatal(http.ListenAndServe(config.RunAddr, r))
}
