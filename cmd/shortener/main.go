// package main реализует вызов всех компонентов нужных для работы сервера и запускает сервер
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/http-server/server"
	"github.com/gogapopp/shortener/internal/app/lib/logger"
	"github.com/gogapopp/shortener/internal/app/storage"
)

// go build flags
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main реализует вызов всех компонентов нужных для работы сервера и запускает сервер
func main() {
	ctx := context.Background()
	// build stdout
	buildStdout()
	// парсим конфиг
	cfg := config.ParseConfig()
	// инициализируем логер
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	// подключаем хранилище
	storage, err := storage.NewRepo(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := server.RunServer(cfg, storage, log)

	// реализация graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint

	// проверяем установленно ли подключение к бд
	db, err := storage.Ping()
	if err == nil {
		if err := db.Close(); err != nil {
			log.Info("error close db conn:", err)
		}
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Info("error shutdown the server:", err)
	}
}

// buildStdout выводит в строку терминала build version, build date, build commit
func buildStdout() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
