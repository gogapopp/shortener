package storage

import (
	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/storage/files"
	"github.com/gogapopp/shortener/internal/app/storage/inmemory"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
)

type Storage interface {
	SaveURL(longURL, shortURL, correlationID string, userID string) error
	GetURL(shortURL, userID string) (string, error)
	GetShortURL(longURL string) string
	BatchInsertURL(urls []models.BatchDatabaseResponse) error
	Ping() error
	GetUserURLs(userID string) ([]models.UserURLs, error)
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
