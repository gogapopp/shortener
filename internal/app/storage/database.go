package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gogapopp/shortener/internal/app/models"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

var ErrConflict = errors.New("ErrConflict")
var ErrConnectToDatabase = errors.New("ErrConnectToDatabase")
var ErrBeginTx = errors.New("ErrBeginTx")
var ErrCreateIndex = errors.New("ErrCreateIndex")
var ErrInserIntoDB = errors.New("ErrInserIntoDB")
var ErrRowsAffected = errors.New("ErrRowsAffected")

// InitializeDatabase инициализирует базу данных если значение dsn не пустое
func InitializeDatabase(ctx context.Context, dsn string) error {
	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		return ErrConnectToDatabase
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return ErrConnectToDatabase
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return ErrBeginTx
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

	_, err = tx.ExecContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS long_url_id ON urls(long_url)")
	if err != nil {
		return ErrCreateIndex
	}

	tx.Commit()
	return nil
}

// InsertURL записывает ссылку в базу данных, если уже имеется то обрабатываем ошибку
func InsertURL(ctx context.Context, shortURL, longURL, correlationID string) error {
	result, err := db.ExecContext(ctx, "INSERT INTO urls (short_url, long_url, correlation_id) VALUES ($1, $2, $3) ON CONFLICT (long_url) DO NOTHING", shortURL, longURL, correlationID)
	if err != nil {
		return ErrInserIntoDB
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrRowsAffected
	}
	if rowsAffected == 0 {
		return ErrConflict
	}
	return nil
}

func BatchInsertURL(ctx context.Context, urls []models.BatchDatabaseResponse) error {
	// собираем запрос
	query := "INSERT INTO urls (short_url, long_url, correlation_id) VALUES "
	values := []interface{}{}

	for i, url := range urls {
		query += fmt.Sprintf("($%d, $%d, $%d),", i*3+1, i*3+2, i*3+3)
		values = append(values, url.ShortURL, url.OriginalURL, url.CorrelationID)
	}
	// удаляем последнюю запятую и обновляем поля
	query = query[:len(query)-1]

	// выполняем запрос
	_, err := db.ExecContext(ctx, query, values...)
	if err != nil {
		return ErrInserIntoDB
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
