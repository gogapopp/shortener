package config

import (
	"context"
	"errors"
	"fmt"

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
	// разрешаем запись в файл
	handlers.WriteToFile(true)
	// инициализируем функции
	if err := storage.CreateFile(); err != nil {
		return err
	}
	storage.Load()
	storage.RestoreURL()
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
			fmt.Println("не удалось начать запись в базу данных:", err)
			handlers.WriteToDatabase(false)
			if err := InitializeFilemamager(StoragePath); err != nil {
				if errors.Is(err, storage.ErrCreateFile) {
					fmt.Println("не удалось начать запись в файл:", err)
					handlers.WriteToFile(false)
				}
			}
		}
	}
}
