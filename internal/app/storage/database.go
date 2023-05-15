package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn
var err error

func InitializeDatabase(dsn string) {
	db, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func DB() *pgx.Conn {
	return db
}
