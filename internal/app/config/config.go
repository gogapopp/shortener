// package config is intended for parsing flags and env
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config contains the configuration parameters for the application
type Config struct {
	// address to start the server
	RunAddr string `env:"SERVER_ADDRESS" json:"server_address"`
	// host address to start the server
	BaseAddr string `env:"BASE_URL" json:"base_url"`
	// name for the record file
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	// data for connecting to the database
	DatabasePath string `env:"DATABASE_DSN" json:"database_dsn"`
	// contains a variable for enabling the TLS certificate for the server
	HTTPSEnable bool `env:"ENABLE_HTTPS" json:"enable_https"`
	// allowed addresses for /api/internal/stats
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// ParseConfig parses flags and environment variables at program startup
func ParseConfig() *Config {
	var cfg Config
	var fileConfig Config
	// get the path to the config.json file
	configpath := os.Getenv("CONFIG")
	if configpath == "" {
		flag.StringVar(&configpath, "c", "", "config.json file path")
	}
	// writing the flags to the config
	flag.StringVar(&cfg.RunAddr, "a", "0.0.0.0:8080", "address and port to run server")
	flag.StringVar(&cfg.BaseAddr, "b", "http://localhost:8080/", "base url")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file storage path")
	flag.StringVar(&cfg.DatabasePath, "d", "", "database path")
	flag.BoolVar(&cfg.HTTPSEnable, "s", false, "https enable")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "subnet")
	flag.Parse()
	// parse config.json
	if configpath != "" {
		fileConfig = ParseConfigFile(configpath)
	}
	// writing end in the config
	cleanenv.ReadEnv(&cfg)
	// writing config.json
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
	if cfg.TrustedSubnet == "" {
		cfg.TrustedSubnet = fileConfig.TrustedSubnet
	}
	// if cfg.HTTPSEnable false
	if !cfg.HTTPSEnable {
		cfg.HTTPSEnable = fileConfig.HTTPSEnable
	}
	return &cfg
}

// ParseConfigFile parses config from config.json file
func ParseConfigFile(configpath string) Config {
	var fileConfig Config
	// writing the config from config.json
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
