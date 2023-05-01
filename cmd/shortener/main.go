package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/logger"
	"github.com/gogapopp/shortener/internal/app/middlewares"
)

var BaseAddr string
var RunAddr string
var StoragePath string

func main() {
	initializeServerConfig()
	// передаём в config адрес сохранения файла
	config.GetStoragePath(StoragePath)
	if StoragePath == "" {
		handlers.WriteToFile(false)
	} else {
		if !strings.HasPrefix(StoragePath, ".") {
			StoragePath = "." + StoragePath
		}
		handlers.WriteToFile(true)
		config.CreateFile()
		config.Load()
		config.RestoreURL()
	}

	// запускаем сервер
	fmt.Println("Running the server at", RunAddr)
	runServer()
}

// RunServer запускает сервер
func runServer() {
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(middlewares.GzipMiddleware(handlers.PostShortURL)))
		r.Get("/{id}", logger.ResponseLogger(middlewares.GzipMiddleware(handlers.GetHandleURL)))
		r.Post("/api/shorten", logger.RequestLogger(middlewares.GzipMiddleware(handlers.PostJSONHandler)))
	})

	log.Fatal(http.ListenAndServe(RunAddr, r))
}

func initializeServerConfig() {
	flags := config.ParseFlags()
	// проверяем есть ли переменные окружения
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flags.FlagRunAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		flags.FlagBaseAddr = envBaseAddr
	}
	if envStoragePath := os.Getenv("FILE_STORAGE_PATHL"); envStoragePath != "" {
		flags.FlagStoragePath = envStoragePath
	}
	// передаём FlagBaseAddr в handlers.go (функция записывает значение в переменную которая находится в пакете handlers)
	BaseAddr := flags.FlagBaseAddr
	// передаём в encryptor адрес
	encryptor.GetBaseAddr(BaseAddr)
	RunAddr = flags.FlagRunAddr
	StoragePath = flags.FlagStoragePath
}
