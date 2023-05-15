package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/logger"
	"github.com/gogapopp/shortener/internal/app/middlewares"
	"github.com/gogapopp/shortener/internal/app/storage"
)

var BaseAddr string
var RunAddr string
var StoragePath string
var DatabaseDSN string

func main() {
	// парсим флаги и env
	initializeServerConfig()
	if DatabaseDSN != "" {
		// инициализируем базу данных и передаём значение запуска базы данных в пакет storage
		fmt.Println(DatabaseDSN)
		storage.InitializeDatabase(DatabaseDSN)
	}
	// передаём в filemanager адрес сохранения файла
	storage.GetStoragePath(StoragePath)
	if StoragePath == "" {
		handlers.WriteToFile(false)
	} else {
		handlers.WriteToFile(true)
		storage.CreateFile()
		storage.Load()
		storage.RestoreURL()
	}

	// запускаем сервер
	fmt.Println("Running the server at", RunAddr)
	runServer()
}

// RunServer() запускает сервер и инициализирует логер
func runServer() {
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(middlewares.GzipMiddleware(handlers.PostShortURL)))
		r.Get("/{id}", logger.ResponseLogger(middlewares.GzipMiddleware(handlers.GetHandleURL)))
		r.Get("/ping", logger.ResponseLogger(handlers.GetPingDatabase))
		r.Post("/api/shorten", logger.RequestJSONLogger(middlewares.GzipMiddleware(handlers.PostJSONHandler)))
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
	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		flags.FlagStoragePath = envStoragePath
	}
	if envDatabasePath := os.Getenv("DATABASE_DSN"); envDatabasePath != "" {
		flags.FlagDatabasePath = envDatabasePath
	}
	// передаём FlagBaseAddr в handlers.go (функция записывает значение в переменную которая находится в пакете handlers)
	BaseAddr := flags.FlagBaseAddr
	// передаём в encryptor адрес
	encryptor.GetBaseAddr(BaseAddr)
	RunAddr = flags.FlagRunAddr
	StoragePath = flags.FlagStoragePath
	DatabaseDSN = flags.FlagDatabasePath
}
