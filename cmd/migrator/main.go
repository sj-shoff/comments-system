package main

import (
	"comments-system/pkg/logger/sl"
	"flag"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := flag.String("dsn", "", "PostgreSQL DSN")
	migrationsPath := flag.String("migrations-path", "", "Path to migrations")
	flag.Parse()

	if *dsn == "" || *migrationsPath == "" {
		slog.Error("Both dsn and migrations-path must be provided")
		os.Exit(1)
	}

	m, err := migrate.New("file://"+*migrationsPath, "postgres://"+*dsn)
	if err != nil {
		slog.Error("Migration initialization failed", "error", sl.Err(err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Migration failed", sl.Err(err))
		os.Exit(1)
	}

	slog.Info("Migrations applied successfully")
}
