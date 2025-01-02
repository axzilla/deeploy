package db

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "modernc.org/sqlite"
)

func Init() (*sql.DB, error) {
	// DB Connection
	db, err := sql.Open("sqlite", "deeploy.db")
	if err != nil {
		return nil, err
	}

	// Run Migrations
	m, err := migrate.New(
		"file://internal/db/migrations",
		"sqlite://deeploy.db",
	)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}
