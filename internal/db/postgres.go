package db

import (
	"database/sql"
	"embed"
	"io/fs" // Added
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

// IMPORTANT: This comment MUST be right above the variable!
//go:embed migrations/*.sql
var migrationFiles embed.FS

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
	RunMigrations(DB)
}

func RunMigrations(sqlDB *sql.DB) {
	// 1. Strip the "migrations" prefix
	strippedFS, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		log.Fatalf("Could not sub-directory migrations: %v", err)
	}

	// 2. Use "." for the path
	d, err := iofs.New(strippedFS, ".")
	if err != nil {
		log.Fatalf("Could not create iofs instance: %v", err)
	}

	// 3. Use NewWithInstance
	driver, _ := postgres.WithInstance(sqlDB, &postgres.Config{})
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied successfully!")
}
