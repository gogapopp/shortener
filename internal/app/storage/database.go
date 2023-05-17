package storage

import (
	"context"
	"database/sql"
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

func InsertURL(ctx context.Context, shortURL, longURL, correlationID string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO urls (short_url, long_url, correlation_id) VALUES ($1, $2, $3)", shortURL, longURL, correlationID)
	return err
}

// DB возвращает значение *sql.DB
func DB() *sql.DB {
	return db
}
