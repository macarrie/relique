package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/sys/unix"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

var dbFileName string = "relique.sqlite"
var pool *sql.DB

//go:embed migrations/*.sql
var fs embed.FS

func Init(dbPath string) error {
	fullPath := fmt.Sprintf("%s/%s", dbPath, dbFileName)
	if err := Open(fullPath, true); err != nil {
		return fmt.Errorf("cannot open database connection: %w", err)
	}
	if err := Migrate(fullPath); err != nil {
		return fmt.Errorf("cannot perform DB migrations: %w", err)
	}

	return nil
}

func Open(dbPath string, RW bool) error {
	slog.With(
		slog.String("path", dbPath),
	).Debug("Opening database connection")

	connection, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared&mode=rwc", dbPath))
	if err != nil {
		return fmt.Errorf("cannot open sqlite connection: %w", err)
	}
	pool = connection

	if _, err := os.Lstat(dbPath); os.IsNotExist(err) {
		// Do not check RW access to DB when it just has been created and the underlying file does not exist yet
		return nil
	}

	if RW {
		if err := unix.Access(dbPath, unix.W_OK); err != nil {
			return fmt.Errorf("cannot get RW access to DB: %w", err)
		}
	}

	if err := connection.Ping(); err != nil {
		return fmt.Errorf("cannot ping DB: %w", err)
	}

	pool.SetMaxOpenConns(1)
	return nil
}

func Migrate(dbPath string) error {
	slog.Debug("Performing database migrations")

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("cannot load database migrations: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, fmt.Sprintf("sqlite3://%s?cache=shared&mode=rwc", dbPath))
	if err != nil {
		return fmt.Errorf("cannot initialize migration assistant: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Debug("No migrations to perform")
		} else {
			return fmt.Errorf("error encountered when performing migrations: %w", err)
		}
	}

	slog.Debug("Database migrations finished")
	return nil
}

func Handler() *sql.DB {
	if pool == nil {
		slog.Error("Encountered empty DB handler, this should not have happened")
		os.Exit(1)
	}

	return pool
}
