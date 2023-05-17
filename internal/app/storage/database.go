package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB
var err error

// InitializeDatabase инициализирует базу данных если значение dsn не пустое
func InitializeDatabase(dsn string) {
	db, err = sql.Open("pgx", dsn)
	fmt.Println(err)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	fmt.Println(err)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id serial PRIMARY KEY,
			short_url TEXT,
			long_url TEXT
		)
	`)
	fmt.Println(err)

	tx.Commit()
}

// func CreateTable() error {
// 	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
// 		id SERIAL PRIMARY KEY,
// 		short_url TEXT NOT NULL,
// 		long_url TEXT NOT NULL
// 	)`)
// 	return err
// }

func InsertURL(ctx context.Context, shortURL, longURL string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO urls (short_url, long_url) VALUES ($1, $2)", shortURL, longURL)
	fmt.Println(err)
	return err
}

// DB возвращает значение *sql.DB
func DB() *sql.DB {
	return db
}
