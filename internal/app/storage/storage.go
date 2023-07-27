package storage

import (
	"context"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/storage/files"
	"github.com/gogapopp/shortener/internal/app/storage/inmemory"
	"go.uber.org/zap"
)

type Storage interface {
	SaveURL(baseURL, longURL, shortURL string) error
	GetURL(longURL string) (string, error)
}

func NewRepo(ctx context.Context, cfg *config.Config, log *zap.SugaredLogger) (Storage, error) {
	switch {
	// case cfg.DatabasePath != "":
	// 	return postgres.NewStorage(ctx, cfg.DatabasePath, log)
	case cfg.FileStoragePath != "":
		return files.NewStorage(cfg.FileStoragePath, log)
	default:
		return inmemory.NewStorage(log), nil
	}

}
