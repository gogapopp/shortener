package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gogapopp/shortener/internal/app/lib/models"
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
        correlation_id TEXT,
		user_id TEXT
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

func (s *storage) SaveURL(longURL, shortURL, correlationID string, userID string) error {
	const op = "storage.postgres.SaveURL"
	result, err := s.db.Exec("INSERT INTO urls (short_url, long_url, correlation_id, user_id) VALUES ($1, $2, $3, $4) ON CONFLICT (long_url) DO NOTHING", shortURL, longURL, correlationID, userID)
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

func (s *storage) GetURL(shortURL, userID string) (string, error) {
	const op = "storage.postgres.GetURL"
	var longURL string
	row := s.db.QueryRow("SELECT long_url FROM urls WHERE short_url = $1", shortURL)
	err := row.Scan(&longURL)
	if err != nil {
		return "", fmt.Errorf("%s: %s", op, err)
	}
	return longURL, nil
}

func (s *storage) Ping() (*sql.DB, error) {
	err := s.db.Ping()
	return s.db, err
}

func (s *storage) BatchInsertURL(urls []models.BatchDatabaseResponse, userID string) error {
	const op = "storage.postgres.BatchInsertURL"
	// собираем запрос
	query := "INSERT INTO urls (short_url, long_url, correlation_id, user_id) VALUES "
	values := []interface{}{}

	for i, url := range urls {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d),", i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, url.ShortURL, url.OriginalURL, url.CorrelationID, userID)
	}
	// удаляем последнюю запятую и обновляем поля
	query = query[:len(query)-1]
	query = fmt.Sprintf("%sON CONFLICT (long_url) DO NOTHING", query)

	// выполняем запрос
	_, err := s.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

func (s *storage) GetShortURL(longURL string) string {
	const op = "storage.postgres.GetURL"
	var shortURL string
	row := s.db.QueryRow("SELECT short_url FROM urls WHERE long_url = $1", longURL)
	row.Scan(&shortURL)
	return shortURL
}

func (s *storage) GetUserURLs(userID string) ([]models.UserURLs, error) {
	const op = "storage.postgres.GetUserURLs"
	rows, err := s.db.Query("SELECT long_url, short_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}
	defer rows.Close()

	var urls []models.UserURLs
	for rows.Next() {
		var url models.UserURLs
		if err := rows.Scan(&url.OriginalURL, &url.ShortURL); err != nil {
			return nil, fmt.Errorf("%s: %s", op, err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return urls, nil
}
