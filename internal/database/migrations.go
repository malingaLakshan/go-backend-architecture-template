// Package database – migrations.go runs database schema migrations.
package database

import (
	"database/sql"
	"log/slog"
)

// migration represents a single database migration.
type migration struct {
	Name string
	SQL  string
}

// allMigrations returns the ordered list of migrations to apply.
// Add new migrations to the end of this slice — never modify existing ones.
func allMigrations() []migration {
	return []migration{
		{
			Name: "001_create_customers_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS customers (
					id         TEXT PRIMARY KEY,
					name       TEXT NOT NULL,
					email      TEXT NOT NULL UNIQUE,
					phone      TEXT,
					status     TEXT NOT NULL,
					created_at TEXT NOT NULL,
					updated_at TEXT NOT NULL
				);
			`,
		},
	}
}

// RunMigrations applies all pending migrations.
// It uses a simple migrations table to track which migrations have already run.
func RunMigrations(db *sql.DB, logger *slog.Logger) error {
	// Create the migrations tracking table if it does not exist.
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name       TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT (datetime('now'))
		);
	`)
	if err != nil {
		logger.Error("failed to create schema_migrations table", "error", err)
		return err
	}

	for _, m := range allMigrations() {
		// Check if this migration has already been applied.
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE name = ?", m.Name).Scan(&count)
		if err != nil {
			logger.Error("failed to check migration status", "migration", m.Name, "error", err)
			return err
		}
		if count > 0 {
			logger.Debug("migration already applied, skipping", "migration", m.Name)
			continue
		}

		// Apply the migration.
		if _, err := db.Exec(m.SQL); err != nil {
			logger.Error("migration failed", "migration", m.Name, "error", err)
			return err
		}

		// Record the migration.
		if _, err := db.Exec("INSERT INTO schema_migrations (name) VALUES (?)", m.Name); err != nil {
			logger.Error("failed to record migration", "migration", m.Name, "error", err)
			return err
		}

		logger.Info("migration applied successfully", "migration", m.Name)
	}

	logger.Info("all migrations completed")
	return nil
}
