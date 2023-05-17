package main

import (
	"fmt"
	"log"
	"net/http"

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
	// парсим флаги
	InitializeServerConfig()
	// проверяем есть ли значения подключения к базе данных
	if DatabaseDSN != "" {
		// инициализируем базу данных и передаём значение запуска базы данных в пакет storage
		storage.InitializeDatabase(DatabaseDSN)
		handlers.WriteToDatabase(true)
	} else if StoragePath != "" {
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

func InitializeServerConfig() {
	config.ParseFlags()
	// передаём FlagBaseAddr в handlers.go (функция записывает значение в переменную которая находится в пакете handlers)
	BaseAddr := config.FlagBaseAddr
	// передаём в encryptor адрес
	encryptor.GetBaseAddr(BaseAddr)
	RunAddr = config.FlagRunAddr
	StoragePath = config.FlagStoragePath
	// передаём в filemanager адрес сохранения файла
	storage.GetStoragePath(StoragePath)
	DatabaseDSN = config.FlagDatabasePath
}
