package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	ErrURLExists = errors.New("url exists")
)

type storage struct {
	db *sql.DB
}

func NewStorage(databaseDSN string) (*storage, error) {
	const op = "storage.postgres.NewStorage"

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS urls (
        id serial PRIMARY KEY,
        short_url TEXT,
        long_url TEXT,
        correlation_id TEXT
    );
    CREATE UNIQUE INDEX IF NOT EXISTS long_url_id ON urls(long_url);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &storage{
		db: db,
	}, nil
}

func (s *storage) SaveURL(baseURL, longURL, shortURL, correlationID string) error {
	const op = "storage.postgres.SaveURL"
	result, err := s.db.Exec("INSERT INTO urls (short_url, long_url, correlation_id) VALUES ($1, $2, $3) ON CONFLICT (long_url) DO NOTHING", shortURL, longURL, correlationID)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	if rowsAffected == 0 {
		return ErrURLExists
	}
	return nil
}

func (s *storage) GetURL(longURL string) (string, error) {
	const op = "storage.postgres.GetURL"
	var shortURL string
	row := s.db.QueryRow("SELECT short_url FROM urls WHERE long_url = $1", longURL)
	err := row.Scan(&shortURL)
	if err != nil {
		return "", fmt.Errorf("%s: %s", op, err)
	}
	return shortURL, nil
}

func (s *storage) Ping() error {
	err := s.db.Ping()
	return err
}
