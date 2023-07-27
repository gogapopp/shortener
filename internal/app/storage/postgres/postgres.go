package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

type storage struct {
	db *sql.DB
}

func NewStorage(ctx context.Context, databaseDSN string, log *zap.SugaredLogger) (*storage, error) {
	const op = "storage.postgres.NewStorage"

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		log.Info(fmt.Sprintf("%s: %s", op, err))
		return nil, err
	}
	return &storage{
		db: db,
	}, nil
}

func (s *storage) SaveURL(baseURL, longURL, shortURL string) error {
	return nil
}

func (s *storage) GetURL(longURL string) (string, error) {
	return "", nil
}

func (s *storage) Ping() error {
	err := s.db.Ping()
	return err
}
