package client

import (
	"reflect"
	"testing"
	"time"

	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/macarrie/relique/internal/types/backup_type"

	"github.com/macarrie/relique/internal/types/module"

	"github.com/macarrie/relique/internal/db"

	log "github.com/macarrie/relique/internal/logging"
)

func SetupTest(t *testing.T) {
	log.IsTest = true
	log.Setup(true, log.TEST_LOG_PATH)
	if err := db.InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}
}

func TestClient_GetLog(t *testing.T) {
	SetupTest(t)

	tests := []struct {
		name    string
		client  Client
		wantNil bool
	}{
		{
			name: "get_log",
			client: Client{
				Name:    "test_client",
				Address: "test_client_addr",
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.GetLog(); (got == nil) != tt.wantNil {
				t.Errorf("GetLog() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func TestClient_String(t *testing.T) {
	SetupTest(t)

	tests := []struct {
		name   string
		client Client
		want   string
	}{
		{
			name: "client_string",
			client: Client{
				Name:    "client",
				Address: "addr",
			},
			want: "client (addr)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadFromPath(t *testing.T) {
	SetupTest(t)
	module.MODULES_INSTALL_PATH = "../../../test/config/modules"

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []Client
		wantErr bool
	}{
		{
			name: "example",
			args: args{path: "../../../test/config/clients/"},
			want: []Client{Client{
				Name:    "example",
				Address: "example",
				Port:    1234,
				Modules: []module.Module{
					module.Module{
						ModuleType:        "default_test",
						Name:              "example-diff",
						BackupType:        backup_type.BackupType{Type: backup_type.Diff},
						ScheduleNames:     []string{"daily"},
						BackupPaths:       []string{"/tmp/example"},
						PreBackupScript:   "/tmp/prebackup.sh",
						PostBackupScript:  "/tmp/postbackup.sh",
						PreRestoreScript:  "/tmp/prerestore.sh",
						PostRestoreScript: "/tmp/postrestore.sh",
					},
					module.Module{
						ModuleType:        "default_test",
						Name:              "example-full",
						BackupType:        backup_type.BackupType{Type: backup_type.Full},
						ScheduleNames:     []string{"daily"},
						BackupPaths:       []string{"/tmp/example"},
						PreBackupScript:   "/tmp/prebackup.sh",
						PostBackupScript:  "/tmp/postbackup.sh",
						PreRestoreScript:  "/tmp/prerestore.sh",
						PostRestoreScript: "/tmp/postrestore.sh",
					},
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFromPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadFromPath() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_loadFromFile(t *testing.T) {
	SetupTest(t)
	module.MODULES_INSTALL_PATH = "../../../test/config/modules"

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
		{
			name: "example",
			args: args{file: "../../../test/config/clients/example.toml"},
			want: Client{
				Name:    "example",
				Address: "example",
				Port:    1234,
				Modules: []module.Module{
					module.Module{
						ModuleType:        "default_test",
						Name:              "example-diff",
						BackupType:        backup_type.BackupType{Type: backup_type.Diff},
						ScheduleNames:     []string{"daily"},
						BackupPaths:       []string{"/tmp/example"},
						PreBackupScript:   "/tmp/prebackup.sh",
						PostBackupScript:  "/tmp/postbackup.sh",
						PreRestoreScript:  "/tmp/prerestore.sh",
						PostRestoreScript: "/tmp/postrestore.sh",
					},
					module.Module{
						ModuleType:        "default_test",
						Name:              "example-full",
						BackupType:        backup_type.BackupType{Type: backup_type.Full},
						ScheduleNames:     []string{"daily"},
						BackupPaths:       []string{"/tmp/example"},
						PreBackupScript:   "/tmp/prebackup.sh",
						PostBackupScript:  "/tmp/postbackup.sh",
						PreRestoreScript:  "/tmp/prerestore.sh",
						PostRestoreScript: "/tmp/postrestore.sh",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "default_values",
			args:    args{file: "../../../test/config/clients/empty.toml"},
			want:    Client{},
			wantErr: false,
		},
		{
			name:    "unreadable",
			args:    args{file: "../../../test/config/clients/unreadable.toml.test"},
			want:    Client{},
			wantErr: true,
		},
		{
			name:    "not_found",
			args:    args{file: "../../../test/config/clients/not_found.toml"},
			want:    Client{},
			wantErr: true,
		},
		{
			name:    "invalid_toml",
			args:    args{file: "../../../test/config/modules/parse_error.toml"},
			want:    Client{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFromFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Valid(t *testing.T) {
	tests := []struct {
		name   string
		client Client
		want   bool
	}{
		{
			name: "valid",
			client: Client{
				Name:    "valid_client",
				Address: "valid_client_addr",
			},
			want: true,
		},
		{
			name: "missing_name",
			client: Client{
				Address: "missing_name_addr",
			},
			want: false,
		},
		{
			name: "missing_address",
			client: Client{
				Name: "missing_addr",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillSchedulesStruct(t *testing.T) {
	client1 := Client{
		Name:    "client1",
		Address: "client1",
		Port:    1234,
		Modules: []module.Module{
			{
				ModuleType:    "foo",
				Name:          "foo",
				BackupType:    backup_type.BackupType{Type: backup_type.Full},
				ScheduleNames: []string{"sched1"},
			},
		},
	}
	client2 := Client{
		Name:    "client2",
		Address: "client2",
		Port:    1234,
		Modules: []module.Module{
			{
				ModuleType:    "foo",
				Name:          "foo",
				BackupType:    backup_type.BackupType{Type: backup_type.Full},
				ScheduleNames: []string{"sched1", "sched2"},
			},
		},
	}

	sched1 := schedule.Schedule{
		Name: "sched1",
		Monday: schedule.Timeranges{
			Ranges: []schedule.Timerange{
				{
					Start: time.Time{}.Add(1 * time.Hour),
					End:   time.Time{}.Add(2 * time.Hour),
				},
			},
		},
	}
	sched2 := schedule.Schedule{
		Name: "sched2",
		Tuesday: schedule.Timeranges{
			Ranges: []schedule.Timerange{
				{
					Start: time.Time{}.Add(1 * time.Hour),
					End:   time.Time{}.Add(2 * time.Hour),
				},
			},
		},
	}

	type args struct {
		client    Client
		schedules []schedule.Schedule
	}
	tests := []struct {
		name               string
		args               args
		wantSchedulesNamed []string
		wantErr            bool
	}{
		{
			name: "single_schedule",
			args: args{
				client:    client1,
				schedules: []schedule.Schedule{sched1},
			},
			wantSchedulesNamed: []string{
				"sched1",
			},
			wantErr: false,
		},
		{
			name: "multiple_schedules",
			args: args{
				client:    client2,
				schedules: []schedule.Schedule{sched1, sched2},
			},
			wantSchedulesNamed: []string{
				"sched1",
				"sched2",
			},
			wantErr: false,
		},
		{
			name: "unknown_schedules",
			args: args{
				client:    client2,
				schedules: []schedule.Schedule{sched2},
			},
			wantSchedulesNamed: []string{},
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FillSchedulesStruct([]Client{tt.args.client}, tt.args.schedules)
			if (err != nil) != tt.wantErr {
				t.Errorf("FillSchedulesStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Do not analyze content on error return
			if tt.wantErr {
				return
			}

			var gotSchedules []string
			for _, mod := range got[0].Modules {
				for _, sched := range mod.Schedules {
					gotSchedules = append(gotSchedules, sched.Name)
				}
			}
			if !reflect.DeepEqual(gotSchedules, tt.wantSchedulesNamed) {
				t.Errorf("FillSchedulesStruct() got = %v, want %v", gotSchedules, tt.wantSchedulesNamed)
			}
		})
	}
}
