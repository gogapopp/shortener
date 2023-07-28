package storage

import (
	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/storage/files"
	"github.com/gogapopp/shortener/internal/app/storage/inmemory"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
)

type Storage interface {
	SaveURL(baseURL, longURL, shortURL, correlationID string) error
	GetURL(longURL string) (string, error)
	Ping() error
}

func NewRepo(cfg *config.Config) (Storage, error) {
	switch {
	case cfg.DatabasePath != "":
		return postgres.NewStorage(cfg.DatabasePath)
	case cfg.FileStoragePath != "":
		return files.NewStorage(cfg.FileStoragePath)
	default:
		return inmemory.NewStorage(), nil
	}

}
