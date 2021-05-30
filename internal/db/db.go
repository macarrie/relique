package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"

	_ "github.com/mattn/go-sqlite3"
)

var pool *sql.DB
var lock sync.RWMutex

const TEST_DB_PATH = "/var/lib/relique/db/unittests.db"

var dbPath = "/var/lib/relique/db/server.db"

func Init() error {
	if err := open(); err != nil {
		return errors.Wrap(err, "cannot open database connection")
	}
	if err := SetupSchema(); err != nil {
		return errors.Wrap(err, "cannot init database schema")
	}
	if err := Migrate(); err != nil {
		return errors.Wrap(err, "cannot perform DB migrations")
	}

	return nil
}

// Used for unit tests
func InitTestDB() error {
	dbPath = TEST_DB_PATH
	if err := ResetTestDB(); err != nil {
		return errors.Wrap(err, "cannot reset test DB")
	}

	return Init()
}

// Used for unit tests
func ResetTestDB() error {
	pool = nil
	dbPath = TEST_DB_PATH
	if _, err := os.Lstat(dbPath); os.IsNotExist(err) {
		// DB not created, nothing to do
		return nil
	}

	if err := os.Remove(dbPath); err != nil {
		return errors.Wrap(err, "cannot delete db file")
	}

	return nil
}

func open() error {
	log.WithFields(log.Fields{
		"path": dbPath,
	}).Info("Opening database connection")

	connection, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared&mode=rwc", dbPath))
	if err != nil {
		return errors.Wrap(err, "cannot open sqlite connection")
	}
	pool = connection

	if _, err := os.Lstat(dbPath); os.IsNotExist(err) {
		// Do not check RW access to DB when it just has been created and the underlying file does not exist yet
		return nil
	}

	if err := unix.Access(dbPath, unix.W_OK); err != nil {
		return errors.Wrap(err, "cannot check RW access to DB")
	}

	if err := connection.Ping(); err != nil {
		return errors.Wrap(err, "cannot ping DB")
	}

	pool.SetMaxOpenConns(1)
	return nil
}

func checkNilPool() {
	if pool == nil {
		log.Fatal("Found empty database connexion handler. This should not have happened")
	}
}

func Write() *sql.DB {
	checkNilPool()
	lock.Lock()
	return pool
}

func Read() *sql.DB {
	checkNilPool()
	lock.RLock()
	return pool
}

func RUnlock() {
	lock.RUnlock()
}

func Unlock() {
	lock.Unlock()
}

func SetupSchema() error {
	log.Info("Setting up database schema")
	schema := `
CREATE TABLE IF NOT EXISTS modules (
	id INTEGER PRIMARY KEY,
	module_type TEXT NOT NULL,
	name TEXT NOT NULL UNIQUE,
	backup_type INTEGER NOT NULL,
	backup_paths TEXT,
	pre_backup_script TEXT,
	post_backup_script TEXT,
	pre_restore_script TEXT,
	post_restore_script TEXT
);
CREATE TABLE IF NOT EXISTS clients (
	 id INTEGER PRIMARY KEY,
	 config_version TEXT,
	 name TEXT NOT NULL UNIQUE,
	 address TEXT NOT NULL,
	 port INTEGER NOT NULL,
	 server_address INTEGER NOT NULL,
	 server_port INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS schedules (
	 id INTEGER PRIMARY KEY,
	 name TEXT NOT NULL UNIQUE,
	 monday TEXT,
	 tuesday TEXT,
	 wednesday TEXT,
	 thursday TEXT,
	 friday TEXT,
	 saturday TEXT,
	 sunday TEXT
);
CREATE TABLE IF NOT EXISTS modules_schedules (
	schedule_id INTEGER,
	module_id INTEGER,
	FOREIGN KEY(schedule_id) REFERENCES schedules(id),
	FOREIGN KEY(module_id) REFERENCES modules(id)
);
CREATE TABLE IF NOT EXISTS jobs (
	id INTEGER PRIMARY KEY,
	uuid TEXT NOT NULL UNIQUE,
	status INTEGER NOT NULL,
	backup_type INTEGER NOT NULL,
	job_type INTEGER NOT NULL,
	done INTEGER NOT NULL,
	module_id INTEGER NOT NULL,
	client_id INTEGER NOT NULL,
	start_time TIMESTAMP,
	end_time TIMESTAMP,
    restore_job_uuid TEXT,
    restore_destination TEXT,
	FOREIGN KEY(module_id) REFERENCES modules(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY(client_id) REFERENCES clients(id) ON DELETE CASCADE ON UPDATE CASCADE
); `

	_, err := Write().Exec(schema)
	defer Unlock()
	if err != nil {
		return errors.Wrap(err, "cannot perform schema init request")
	}

	return nil
}

func Migrate() error {
	log.Info("Performing database migrations")
	// TODO: Perform database migrations
	return nil
}
