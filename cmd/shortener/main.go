package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/routes"
)

func main() {
	ctx := context.Background()
	// парсим флаги
	config.InitializeServerConfig()
	// пытается подключится к базе данных, если не получается то пытаеимся записывать в файл, если не получается то в память
	config.SetupDatabaseAndFilemanager(ctx)
	// запускаем сервер
	fmt.Println("Running the server at", config.RunAddr)
	r := routes.Routes()
	log.Fatal(http.ListenAndServe(config.RunAddr, r))
}
