// package config предназначен для парсинга флагов и env в структуру Config
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит параметры конфигурации для приложения
type Config struct {
	// адрес для запуска сервера
	RunAddr string `env:"SERVER_ADDRESS" json:"server_address"`
	// адрес хоста для запуска сервера
	BaseAddr string `env:"BASE_URL" json:"base_url"`
	// название для файла записи
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	// данные для подключения к БД
	DatabasePath string `env:"DATABASE_DSN" json:"database_dsn"`
	// содержит переменную для включения TLS сертификата для сервера
	HTTPSEnable bool `env:"ENABLE_HTTPS" json:"enable_https"`
}

// ParseConfig парсит флаги и переменные окружения при запуске программы
func ParseConfig() *Config {
	var cfg Config
	var fileConfig Config
	// получаем путь к файлу config.json
	configpath := os.Getenv("CONFIG")
	if configpath == "" {
		flag.StringVar(&configpath, "c", "", "config.json file path")
	}
	// записываем флаги в конфиг
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.BaseAddr, "b", "http://localhost:8080/", "base url")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file storage path")
	flag.StringVar(&cfg.DatabasePath, "d", "", "database path")
	flag.BoolVar(&cfg.HTTPSEnable, "s", false, "https enable")
	flag.Parse()
	// парсим config.json
	if configpath != "" {
		fileConfig = ParseConfigFile(configpath)
	}
	// записываем env в конфиг
	cleanenv.ReadEnv(&cfg)
	// записываем config.json
	if cfg.BaseAddr == "" {
		cfg.BaseAddr = fileConfig.BaseAddr
	}
	if cfg.DatabasePath == "" {
		cfg.DatabasePath = fileConfig.DatabasePath
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = fileConfig.FileStoragePath
	}
	if cfg.RunAddr == "" {
		cfg.RunAddr = fileConfig.RunAddr
	}
	// если cfg.HTTPSEnable false
	if !cfg.HTTPSEnable {
		cfg.HTTPSEnable = fileConfig.HTTPSEnable
	}
	return &cfg
}

// ParseConfigFile парсит конфиг из файла config.json
func ParseConfigFile(configpath string) Config {
	var fileConfig Config
	// записываем конфиг из config.json
	file, err := os.Open(configpath)
	if err != nil {
		log.Fatalf("error read config.json file: %s", err)
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&fileConfig); err != nil {
		log.Fatalf("error decode config.json file: %s", err)
	}
	return fileConfig
}
