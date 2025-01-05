package db

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Init() (*sql.DB, error) {
	// Create data directory if it doesn't exist
	err := os.MkdirAll("internal/app/data", 0755)
	if err != nil {
		return nil, err
	}
	// Open database connection - this will create the DB file if it doesn't exist
	db, err := sql.Open("sqlite", "internal/app/data/deeploy.db")
	if err != nil {
		return nil, err
	}

	// Enable Write-Ahead Logging (WAL) mode for better performance and concurrency
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	// Set timeout for busy connections to 5 seconds to handle concurrent access
	_, err = db.Exec("PRAGMA busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	// Adjust synchronous setting for better performance while maintaining safety
	_, err = db.Exec("PRAGMA synchronous=NORMAL")
	if err != nil {
		return nil, err
	}

	// Create Migration-Source
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatalf("Loading migrations failed: %v", err)
	}

	// Create Migrations-Object
	m, err := migrate.NewWithSourceInstance("iofs", source, "sqlite://internal/app/data/deeploy.db")
	if err != nil {
		log.Fatalf("Initialize migrations failed: %v", err)
	}

	// Run all pending migrations
	// ErrNoChange is ok - it just means we're up to date
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}
