package config

import (
	"flag"
	"os"
)

var (
	FlagRunAddr      string
	FlagBaseAddr     string
	FlagStoragePath  string
	FlagDatabasePath string
)

// ParseFlags парсит флаги при запуске программы
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8080/", "base url")
	flag.StringVar(&FlagStoragePath, "f", "", "file storage path")
	flag.StringVar(&FlagDatabasePath, "d", "", "database path")
	flag.Parse()

	// проверяем есть ли переменные окружения
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		FlagBaseAddr = envBaseAddr
	}
	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		FlagStoragePath = envStoragePath
	}
	if envDatabasePath := os.Getenv("DATABASE_DSN"); envDatabasePath != "" {
		FlagDatabasePath = envDatabasePath
	}
}
