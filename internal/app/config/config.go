// package config предназначен для парсинга флагов и env в структуру Config
package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config хранит информацию о конфиге
type Config struct {
	// адрес для запуска сервера
	RunAddr string `env:"SERVER_ADDRESS"`
	// адрес хоста для запуска сервера
	BaseAddr string `env:"BASE_URL"`
	// название для файла записи
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// данные для подключения к БД
	DatabasePath string `env:"DATABASE_DSN"`
}

// ParseConfig парсит флаги и переменные окружения при запуске программы
func ParseConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.BaseAddr, "b", "http://localhost:8080/", "base url")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file storage path")
	flag.StringVar(&cfg.DatabasePath, "d", "", "database path")
	flag.Parse()
	cleanenv.ReadEnv(&cfg)
	return &cfg
}
