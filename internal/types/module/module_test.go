package module

import (
	"reflect"
	"testing"

	"github.com/macarrie/relique/internal/types/custom_errors"

	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/logging"

	"github.com/macarrie/relique/internal/types/backup_type"
)

func SetupTest(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	if err := db.InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}
}

func TestGetByID(t *testing.T) {
	SetupTest(t)

	testModule := Module{
		ModuleType:        "test",
		Name:              "test",
		BackupType:        backup_type.BackupType{Type: backup_type.Full},
		Schedules:         nil,
		BackupPaths:       []string{"path1", "path2"},
		PreBackupScript:   "",
		PostBackupScript:  "",
		PreRestoreScript:  "",
		PostRestoreScript: "",
	}
	if _, err := testModule.Save(); err != nil {
		t.Errorf("cannot save module: '%s'", err)
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name         string
		args         args
		wantName     string
		wantErr      bool
		wantNotFound bool
	}{
		{
			name:         "get existing module",
			args:         args{id: testModule.ID},
			wantName:     "test",
			wantErr:      false,
			wantNotFound: false,
		},
		{
			name:         "get unknown module",
			args:         args{id: 1234},
			wantName:     "",
			wantErr:      true,
			wantNotFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNotFound && !custom_errors.IsDBNotFoundError(err) {
				t.Errorf("GetByID() error = %v, wantNotFound %v", err, tt.wantNotFound)
				return
			}
			if got.Name != tt.wantName {
				t.Errorf("GetByID() got = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

func TestGetID(t *testing.T) {
	SetupTest(t)

	testModule := Module{
		ModuleType:        "test",
		Name:              "test",
		BackupType:        backup_type.BackupType{Type: backup_type.Full},
		Schedules:         nil,
		BackupPaths:       []string{"path1", "path2"},
		PreBackupScript:   "",
		PostBackupScript:  "",
		PreRestoreScript:  "",
		PostRestoreScript: "",
	}
	if _, err := testModule.Save(); err != nil {
		t.Errorf("cannot save module: '%s'", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name         string
		args         args
		want         int64
		wantErr      bool
		wantNotFound bool
	}{
		{
			name:         "normal",
			args:         args{name: "test"},
			want:         testModule.ID,
			wantErr:      false,
			wantNotFound: false,
		},
		{
			name:         "not_found",
			args:         args{name: "not found"},
			want:         0,
			wantErr:      true,
			wantNotFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetID(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNotFound && !custom_errors.IsDBNotFoundError(err) {
				t.Errorf("GetID() error = %v, wantNotFound %v", err, tt.wantNotFound)
				return
			}
			if got != tt.want {
				t.Errorf("GetID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	SetupTest(t)

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    Module
		wantErr bool
	}{
		{
			name: "example",
			args: args{file: "../../../test/config/modules/example.toml"},
			want: Module{
				ID:         0,
				ModuleType: "example",
				Name:       "example_module",
				BackupType: backup_type.BackupType{
					Type: backup_type.Full,
				},
				Schedules:         nil,
				BackupPaths:       []string{"/tmp/example"},
				PreBackupScript:   "/tmp/prebackup.sh",
				PostBackupScript:  "/tmp/postbackup.sh",
				PreRestoreScript:  "/tmp/prerestore.sh",
				PostRestoreScript: "/tmp/postrestore.sh",
			},
			wantErr: false,
		},
		{
			name:    "default_values",
			args:    args{file: "../../../test/config/modules/empty.toml"},
			want:    Module{},
			wantErr: false,
		},
		{
			name:    "unreadable",
			args:    args{file: "../../../test/config/modules/unreadable.toml.test"},
			want:    Module{},
			wantErr: true,
		},
		{
			name:    "not_found",
			args:    args{file: "../../../test/config/modules/not_found.toml"},
			want:    Module{},
			wantErr: true,
		},
		{
			name:    "invalid_toml",
			args:    args{file: "../../../test/config/modules/parse_error.toml"},
			want:    Module{},
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

func TestModule_GetLog(t *testing.T) {
	SetupTest(t)

	type fields struct {
		ID                int64
		ModuleType        string
		Name              string
		BackupType        backup_type.BackupType
		Schedules         []string
		BackupPaths       []string
		PreBackupScript   string
		PostBackupScript  string
		PreRestoreScript  string
		PostRestoreScript string
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
	}{
		{
			name: "module log",
			fields: fields{
				ModuleType: "test_module",
				Name:       "test_module",
				BackupType: backup_type.BackupType{Type: backup_type.Full},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				ID:                tt.fields.ID,
				ModuleType:        tt.fields.ModuleType,
				Name:              tt.fields.Name,
				BackupType:        tt.fields.BackupType,
				Schedules:         tt.fields.Schedules,
				BackupPaths:       tt.fields.BackupPaths,
				PreBackupScript:   tt.fields.PreBackupScript,
				PostBackupScript:  tt.fields.PostBackupScript,
				PreRestoreScript:  tt.fields.PreRestoreScript,
				PostRestoreScript: tt.fields.PostRestoreScript,
			}
			if got := m.GetLog(); (got == nil) != tt.wantNil {
				t.Errorf("GetLog() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func TestModule_LoadDefaultConfiguration(t *testing.T) {
	SetupTest(t)

	MODULES_INSTALL_PATH = "../../../test/config/modules"

	tests := []struct {
		name         string
		mod          Module
		wantedValues map[string]interface{}
		wantErr      bool
	}{
		{
			name: "backup_paths",
			mod: Module{
				ModuleType:        "default_test",
				Name:              "backup_paths",
				BackupPaths:       nil,
				PreBackupScript:   "not_empty",
				PostBackupScript:  "not_empty",
				PreRestoreScript:  "not_empty",
				PostRestoreScript: "not_empty",
			},
			wantedValues: map[string]interface{}{
				"BackupPaths":       []string{"/tmp/example"},
				"PreBackupScript":   "not_empty",
				"PostBackupScript":  "not_empty",
				"PreRestoreScript":  "not_empty",
				"PostRestoreScript": "not_empty",
			},
			wantErr: false,
		},
		{
			name: "prebackup",
			mod: Module{
				ModuleType:        "default_test",
				Name:              "prebackup",
				BackupPaths:       []string{"not_empty"},
				PreBackupScript:   "",
				PostBackupScript:  "not_empty",
				PreRestoreScript:  "not_empty",
				PostRestoreScript: "not_empty",
			},
			wantedValues: map[string]interface{}{
				"BackupPaths":       []string{"not_empty"},
				"PreBackupScript":   "/tmp/prebackup.sh",
				"PostBackupScript":  "not_empty",
				"PreRestoreScript":  "not_empty",
				"PostRestoreScript": "not_empty",
			},
			wantErr: false,
		},
		{
			name: "postbackup",
			mod: Module{
				ModuleType:        "default_test",
				Name:              "backup_paths",
				BackupPaths:       []string{"not_empty"},
				PreBackupScript:   "not_empty",
				PostBackupScript:  "",
				PreRestoreScript:  "not_empty",
				PostRestoreScript: "not_empty",
			},
			wantedValues: map[string]interface{}{
				"BackupPaths":       []string{"not_empty"},
				"PreBackupScript":   "not_empty",
				"PostBackupScript":  "/tmp/postbackup.sh",
				"PreRestoreScript":  "not_empty",
				"PostRestoreScript": "not_empty",
			},
			wantErr: false,
		},
		{
			name: "prerestore",
			mod: Module{
				ModuleType:        "default_test",
				Name:              "prerestore",
				BackupPaths:       []string{"not_empty"},
				PreBackupScript:   "not_empty",
				PostBackupScript:  "not_empty",
				PreRestoreScript:  "",
				PostRestoreScript: "not_empty",
			},
			wantedValues: map[string]interface{}{
				"BackupPaths":       []string{"not_empty"},
				"PreBackupScript":   "not_empty",
				"PostBackupScript":  "not_empty",
				"PreRestoreScript":  "/tmp/prerestore.sh",
				"PostRestoreScript": "not_empty",
			},
			wantErr: false,
		},
		{
			name: "postrestore",
			mod: Module{
				ModuleType:        "default_test",
				Name:              "postrestore",
				BackupPaths:       []string{"not_empty"},
				PreBackupScript:   "not_empty",
				PostBackupScript:  "not_empty",
				PreRestoreScript:  "not_empty",
				PostRestoreScript: "",
			},
			wantedValues: map[string]interface{}{
				"BackupPaths":       []string{"not_empty"},
				"PreBackupScript":   "not_empty",
				"PostBackupScript":  "not_empty",
				"PreRestoreScript":  "not_empty",
				"PostRestoreScript": "/tmp/postrestore.sh",
			},
			wantErr: false,
		},
		{
			name: "not_installed_module",
			mod: Module{
				ModuleType: "not_installed_module",
				Name:       "not_installed_module",
			},
			wantedValues: map[string]interface{}{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mod.LoadDefaultConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("LoadDefaultConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
			r := reflect.ValueOf(tt.mod)
			for key, val := range tt.wantedValues {
				got := reflect.Indirect(r).FieldByName(key)
				wanted := reflect.ValueOf(val)
				if got.String() != wanted.String() {
					t.Errorf("LoadDefaultConfiguration() key = %v, wanted = %v, got %v", key, wanted, got)
				}
			}
		})
	}
}

func TestModule_Save(t *testing.T) {
	SetupTest(t)

	savedModule := Module{
		ID:                0,
		ModuleType:        "test",
		Name:              "testModule",
		BackupType:        backup_type.BackupType{Type: backup_type.Full},
		Schedules:         nil,
		BackupPaths:       nil,
		PreBackupScript:   "",
		PostBackupScript:  "",
		PreRestoreScript:  "",
		PostRestoreScript: "",
	}
	_, err := savedModule.Save()
	if err != nil {
		t.Errorf("Cannot save module for save test: '%s'", err)
	}

	tests := []struct {
		name    string
		mod     Module
		wantID  bool
		wantErr bool
	}{
		{
			name: "save_new_module",
			mod: Module{
				ModuleType:        "new_module_type",
				Name:              "new_module_to_save",
				BackupType:        backup_type.BackupType{Type: backup_type.Full},
				Schedules:         nil,
				BackupPaths:       nil,
				PreBackupScript:   "not_empty",
				PostBackupScript:  "not_empty",
				PreRestoreScript:  "not_empty",
				PostRestoreScript: "not_empty",
			},
			wantID:  true,
			wantErr: false,
		},
		{
			name:    "existing_item",
			mod:     savedModule,
			wantID:  true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mod.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantID && (got == 0) {
				t.Errorf("Save() got = %v, wanted not null ID", got)
			}

			modFromDB, err := GetByID(got)
			if err != nil {
				t.Errorf("Save() cannot get module from DB, err = '%s'", err)
			}
			if !reflect.DeepEqual(modFromDB, tt.mod) {
				t.Errorf("Save() mod = %v, from_db = %v", tt.mod, modFromDB)
			}
		})
	}
}

func TestModule_String(t *testing.T) {
	SetupTest(t)

	tests := []struct {
		name string
		mod  Module
		want string
	}{
		{
			name: "module",
			mod:  Module{Name: "test"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mod.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_Update(t *testing.T) {
	SetupTest(t)

	savedModule := Module{
		ID:                0,
		ModuleType:        "test",
		Name:              "testModule",
		BackupType:        backup_type.BackupType{Type: backup_type.Full},
		Schedules:         nil,
		BackupPaths:       nil,
		PreBackupScript:   "",
		PostBackupScript:  "",
		PreRestoreScript:  "",
		PostRestoreScript: "",
	}
	savedId, err := savedModule.Save()
	if err != nil {
		t.Errorf("Cannot save module for update test: '%s'", err)
	}

	tests := []struct {
		name    string
		fields  Module
		want    int64
		wantErr bool
	}{
		{
			name:    "not_saved_module",
			fields:  Module{},
			want:    0,
			wantErr: true,
		},
		{
			name:    "savedModule",
			fields:  savedModule,
			want:    savedId,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				ID:                tt.fields.ID,
				ModuleType:        tt.fields.ModuleType,
				Name:              tt.fields.Name,
				BackupType:        tt.fields.BackupType,
				Schedules:         tt.fields.Schedules,
				BackupPaths:       tt.fields.BackupPaths,
				PreBackupScript:   tt.fields.PreBackupScript,
				PostBackupScript:  tt.fields.PostBackupScript,
				PreRestoreScript:  tt.fields.PreRestoreScript,
				PostRestoreScript: tt.fields.PostRestoreScript,
			}
			got, err := m.Update()
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_Valid(t *testing.T) {
	tests := []struct {
		name string
		mod  Module
		want bool
	}{
		{
			name: "valid_module",
			mod: Module{
				ModuleType: "not_empty",
				Name:       "valid_module",
				BackupType: backup_type.BackupType{Type: backup_type.Diff},
			},
			want: true,
		},
		{
			name: "unknown_backup_type",
			mod: Module{
				ModuleType: "not_empty",
				Name:       "unknown_backup_type",
				BackupType: backup_type.BackupType{Type: backup_type.Unknown},
			},
			want: false,
		},
		{
			name: "missing_name",
			mod: Module{
				ModuleType: "not_empty",
				BackupType: backup_type.BackupType{Type: backup_type.Full},
			},
			want: false,
		},
		{
			name: "missing_module_type",
			mod: Module{
				Name:       "missing_module_type",
				BackupType: backup_type.BackupType{Type: backup_type.Full},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mod.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}