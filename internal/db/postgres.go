package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}
}
