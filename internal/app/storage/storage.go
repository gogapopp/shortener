// package storage содержит описание методов хранилища
package storage

import (
	"database/sql"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/storage/files"
	"github.com/gogapopp/shortener/internal/app/storage/inmemory"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
)

// Storage определяет методы хранилища storage
//
//go:generate mockgen -source=storage.go -destination=storagemocks/mock.go
type Storage interface {
	BatchInsertURL(urls []models.BatchDatabaseResponse, userID string) error
	SaveURL(longURL, shortURL, correlationID string, userID string) error
	GetURL(shortURL, userID string) (bool, string, error)
	GetUserURLs(userID string) ([]models.UserURLs, error)
	SetDeleteFlag(IDs []string, userID string) error
	GetShortURL(longURL string) string
	Ping() (*sql.DB, error)
}

// NewRepo согласно конфигу определяет тип хранилища
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
