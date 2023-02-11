package db

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"sync"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"

	_ "github.com/mattn/go-sqlite3"
)

var pool *sql.DB
var lock sync.RWMutex

const TEST_DB_PATH = "/tmp/unittests.db"

var IsTest = false

var DbPath string
var DbPathReadInConfig bool

// Set default value for dbPath according to OS if not already set in configuration file
func setDbPathDefaultValue() {
	if IsTest {
		// Let test suite define DBPath
		return
	}

	if DbPathReadInConfig {
		return
	}

	switch runtime.GOOS {
	case "freebsd":
		DbPath = "/usr/local/relique/db/server.db"
	default:
		DbPath = "/var/lib/relique/db/server.db"
	}

	log.WithFields(log.Fields{
		"path": DbPath,
	}).Debug("Set default path for db")
}

func Init() error {
	if err := Open(true); err != nil {
		return errors.Wrap(err, "cannot Open database connection")
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
	IsTest = true
	if err := ResetTestDB(); err != nil {
		return errors.Wrap(err, "cannot reset test DB")
	}

	return Init()
}

// Used for unit tests
func ResetTestDB() error {
	pool = nil
	DbPath = TEST_DB_PATH
	if _, err := os.Lstat(DbPath); os.IsNotExist(err) {
		// DB not created, nothing to do
		return nil
	}

	if err := os.Remove(DbPath); err != nil {
		return errors.Wrap(err, "cannot delete db file")
	}

	return nil
}

func Open(RW bool) error {
	setDbPathDefaultValue()

	log.WithFields(log.Fields{
		"path": DbPath,
	}).Info("Opening database connection")

	connection, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared&mode=rwc", DbPath))
	if err != nil {
		return errors.Wrap(err, "cannot open sqlite connection")
	}
	pool = connection

	if _, err := os.Lstat(DbPath); os.IsNotExist(err) {
		// Do not check RW access to DB when it just has been created and the underlying file does not exist yet
		return nil
	}

	if RW {
		if err := unix.Access(DbPath, unix.W_OK); err != nil {
			return errors.Wrap(err, "cannot get RW access to DB")
		}
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

CREATE TABLE IF NOT EXISTS jobs (
	id INTEGER PRIMARY KEY,
	uuid TEXT NOT NULL UNIQUE,
	status INTEGER NOT NULL,
	backup_type INTEGER NOT NULL,
	job_type INTEGER NOT NULL,
	done INTEGER NOT NULL,
	start_time TIMESTAMP,
	end_time TIMESTAMP,
    restore_job_uuid TEXT,
    restore_destination TEXT,
    storage_root TEXT,
    module_type TEXT,
    client_name TEXT
);`

	_, err := Write().Exec(schema)
	defer Unlock()
	if err != nil {
		return errors.Wrap(err, "cannot perform schema init request")
	}

	return nil
}

func Migrate() error {
	log.Info("Performing database migrations")

	//if err := migration_migrationName(); err != nil {
	//log.WithFields(log.Fields{
	//"err": err,
	//}).Fatal("Cannot perform DB migration. Exiting")
	//}

	return nil
}
