package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/routes"
	"github.com/gogapopp/shortener/internal/app/storage"
)

var BaseAddr string
var RunAddr string
var StoragePath string
var DatabaseDSN string

func main() {
	// парсим флаги и env
	initializeServerConfig()
	// проверяем есть ли значения подключения к базе данных
	if DatabaseDSN != "" {
		// инициализируем базу данных и передаём значение запуска базы данных в пакет storage
		storage.InitializeDatabase(DatabaseDSN)
	}
	// передаём в filemanager адрес сохранения файла
	storage.GetStoragePath(StoragePath)
	// проверяем есть ли путь сохранения файлов записи
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
	r := routes.Routes()
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
