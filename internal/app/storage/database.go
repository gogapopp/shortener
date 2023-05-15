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

func DB() *sql.DB {
	return db
}
