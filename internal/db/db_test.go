package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	tmp := t.TempDir()

	if err := Init(tmp); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, dbFileName)); err != nil {
		t.Errorf("Init() did not create DB file: %v", err)
	}

	// Schema must have been applied
	for _, table := range []string{"jobs", "images"} {
		var name string
		err := Handler().QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %s not found after migrations: %v", table, err)
		}
	}

	// Running Init a second time must be a no-op for migrations
	if err := Init(tmp); err != nil {
		t.Errorf("second Init() error = %v", err)
	}
}

func TestInitInvalidPath(t *testing.T) {
	if err := Init("/does/not/exist/relique"); err == nil {
		t.Errorf("Init() on invalid path should return an error")
	}
}

func TestOpenReadOnlyFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root bypasses file permissions")
	}

	tmp := t.TempDir()
	if err := Init(tmp); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	dbPath := filepath.Join(tmp, dbFileName)
	if err := os.Chmod(dbPath, 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dbPath, 0644)

	if err := Open(dbPath, true); err == nil {
		t.Errorf("Open() in RW mode on read only file should return an error")
	}
	if err := Open(dbPath, false); err != nil {
		t.Errorf("Open() in RO mode on read only file error = %v", err)
	}
}
