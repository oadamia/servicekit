package servicekit

import (
	"database/sql"
	"log/slog"

	"github.com/pressly/goose"
)

type DBConfig interface {
	DBConnString() string
	DBMigrationPath() string
}

// Connect connects to the database and returns the connection
func Connect(config DBConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DBConnString())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping database", "error", err)
		return nil, err
	}
	return db, nil
}

// Migrate runs the database migrations
func Migrate(config DBConfig, db *sql.DB) error {
	err := goose.Up(db, config.DBMigrationPath())
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		return err
	}
	return nil
}
