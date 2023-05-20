package config

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/storage"
)

var BaseAddr string
var RunAddr string
var StoragePath string
var DatabaseDSN string

// InitializeFilemamager инициализируе функции создания файла, считывает его и заполняет значения из файла в память
func InitializeFilemamager(storagePath string) error {
	// инициализируем функции
	if err := storage.CreateFile(); err != nil {
		return err
	}
	if err := storage.Load(); err != nil {
		return err
	}
	storage.RestoreURL()
	// разрешаем запись в файл
	handlers.WriteToFile(true)
	return nil
}

// InitializeServerConfig парсим флаги и записывает в переменную окружения
func InitializeServerConfig() {
	ParseFlags()
	// передаём FlagBaseAddr в handlers.go (функция записывает значение в переменную которая находится в пакете handlers)
	BaseAddr := FlagBaseAddr
	// передаём в encryptor адрес
	encryptor.GetBaseAddr(BaseAddr)
	RunAddr = FlagRunAddr
	StoragePath = FlagStoragePath
	// передаём в filemanager адрес сохранения файла
	storage.GetStoragePath(StoragePath)
	DatabaseDSN = FlagDatabasePath
}

// SetupDatabaseAndFilemanager пытается запустить базу данных, если приходит ошибка то пытается запустить файл менеджер, если опять ошибка начинает запись в память
func SetupDatabaseAndFilemanager(ctx context.Context) {
	if err := storage.InitializeDatabase(ctx, DatabaseDSN); err != nil {
		if errors.Is(err, storage.ErrConnectToDatabase) {
			fmt.Println("не удалось инициализировать базу данных:", err)
			handlers.WriteToDatabase(false)
		} else {
			fmt.Println("не удалось начать запись в базу данных:", err)
			handlers.WriteToDatabase(false)
			log.Fatal(err)
		}
	}
	if err := InitializeFilemamager(StoragePath); err != nil {
		if errors.Is(err, storage.ErrCreateFile) {
			fmt.Println("не удалось создать файл:", err)
		} else {
			fmt.Println("не удалось начать запись в файл:", err)
			log.Fatal(err)
		}
	}
}
