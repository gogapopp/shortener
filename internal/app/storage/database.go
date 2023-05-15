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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
}

// DB возвращает значение *sql.DB
func DB() *sql.DB {
	return db
}
