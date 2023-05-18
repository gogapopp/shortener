package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB
var err error

// InitializeDatabase инициализирует базу данных если значение dsn не пустое
func InitializeDatabase(dsn string) {
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id serial PRIMARY KEY,
			short_url TEXT,
			long_url TEXT,
			correlation_id TEXT
		)
	`)

	tx.Commit()
}

var ErrConflict = errors.New("ErrConflict")

// InsertURL записывает ссылку в базу данных, если уже имеется то обрабатываем ошибку
func InsertURL(ctx context.Context, shortURL, longURL, correlationID string) error {
	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS long_url_id ON urls(long_url)")
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.ExecContext(ctx, "INSERT INTO urls (short_url, long_url, correlation_id) VALUES ($1, $2, $3) ON CONFLICT (long_url) DO NOTHING", shortURL, longURL, correlationID)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rowsAffected == 0 {
		return ErrConflict
	}
	return nil
}

// FindShortURL получаем из базы данных short_url которая соответсвует longURL
func FindShortURL(ctx context.Context, longURL string) string {
	var shortURL string
	row := db.QueryRowContext(ctx, "SELECT short_url FROM urls WHERE long_url = $1", longURL)
	err := row.Scan(&shortURL)
	if err != nil {
		log.Fatal(err)
	}
	return shortURL
}

// DB возвращает значение *sql.DB
func DB() *sql.DB {
	return db
}
