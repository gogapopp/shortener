package storage

import (
	"context"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/storage/inmemory"
	"go.uber.org/zap"
)

type Storage interface {
	SaveURL(baseURL, longURL, shortURL string)
}

func NewRepo(ctx context.Context, cfg *config.Config, log *zap.SugaredLogger) Storage {
	return inmemory.NewStorage(log)
	// if cfg.DatabasePath != "" {
	// 	return postgres.NewStorage(ctx, cfg.DatabasePath, log)
	// } else if cfg.FileStoragePath != "" {
	// 	return files.NewStorage(cfg.FileStoragePath, log)
	// } else {
	// 	return inmemory.NewStorage(log)
	// }
}
