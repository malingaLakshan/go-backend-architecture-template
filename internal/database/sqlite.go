// Package database provides SQLite connection management.
package database

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // SQLite driver (pure Go, no CGO required).
)

// Connect opens a connection to the SQLite database file.
// It creates the data directory if it does not exist.
func Connect(dbPath string, logger *slog.Logger) (*sql.DB, error) {
	// Ensure the directory for the database file exists.
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Error("failed to create database directory", "path", dir, "error", err)
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Error("failed to open database", "path", dbPath, "error", err)
		return nil, err
	}

	// Verify the connection is working.
	if err := db.Ping(); err != nil {
		logger.Error("failed to ping database", "error", err)
		return nil, err
	}

	// Enable WAL mode for better concurrent read performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		logger.Warn("failed to enable WAL mode", "error", err)
	}

	// Enable foreign keys.
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		logger.Warn("failed to enable foreign keys", "error", err)
	}

	logger.Info("database connected successfully", "path", dbPath)
	return db, nil
}
