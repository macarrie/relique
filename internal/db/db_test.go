package db

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/macarrie/relique/internal/logging"
)

func TestInit(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)

	tests := []struct {
		name         string
		databasePath string
		wantErr      bool
	}{
		{
			name:         "successful init",
			databasePath: TEST_DB_PATH,
			wantErr:      false,
		},
		{
			name:         "unwritable database file",
			databasePath: "../../test/db/unwritable.db",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbPath = tt.databasePath
			if err := Init(); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMigrate(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	dbPath = TEST_DB_PATH

	if err := open(); err != nil {
		t.Errorf("cannot init db: '%s'", err)
	}
	if err := SetupSchema(); err != nil {
		t.Errorf("cannot setup db schema: '%s'", err)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "migrate",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Migrate(); (err != nil) != tt.wantErr {
				t.Errorf("Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRUnlock(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	if err := InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}

	tests := []struct {
		name string
		run  func()
	}{
		{
			name: "normal runlock",
			run: func() {
				Read()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
			RUnlock()
		})
	}
}

func TestRead(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	if err := InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}

	tests := []struct {
		name string
		want *sql.DB
	}{
		// TODO: Add test cases.
		{
			name: "read",
			want: pool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer RUnlock()
			if got := Read(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetupSchema(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	dbPath = TEST_DB_PATH

	if err := open(); err != nil {
		t.Errorf("cannot init db: '%s'", err)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "init schema",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetupSchema(); (err != nil) != tt.wantErr {
				t.Errorf("SetupSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnlock(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	if err := InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}

	tests := []struct {
		name string
		run  func()
	}{
		{
			name: "normal unlock",
			run: func() {
				Write()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
			Unlock()
		})
	}
}

func TestWrite(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	if err := InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}

	tests := []struct {
		name string
		want *sql.DB
	}{
		{
			name: "write",
			want: pool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Unlock()
			if got := Write(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_open(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	dbPath = TEST_DB_PATH

	tests := []struct {
		name         string
		databasePath string
		wantErr      bool
	}{
		{
			name:         "open_regular_db_file",
			databasePath: TEST_DB_PATH,
			wantErr:      false,
		},
		{
			name:         "unwritable_db_file",
			databasePath: "../../test/db/unwritable.db",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbPath = tt.databasePath
			if err := open(); (err != nil) != tt.wantErr {
				t.Errorf("open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
