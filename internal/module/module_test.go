package module

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/macarrie/relique/internal/backup_type"
)

func TestLoadFromFile(t *testing.T) {
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
			args: args{file: "../../test/modules/example.toml"},
			want: Module{
				ModuleType: "example",
				Name:       "example_module",
				BackupType: backup_type.BackupType{
					Type: backup_type.Full,
				},
				BackupPaths: []string{"/tmp/example"},
			},
			wantErr: false,
		},
		{
			name:    "default_values",
			args:    args{file: "../../test/modules/empty.toml"},
			want:    Module{},
			wantErr: true,
		},
		{
			name:    "unreadable",
			args:    args{file: "../../test/modules/unreadable.toml.test"},
			want:    Module{},
			wantErr: true,
		},
		{
			name:    "not_found",
			args:    args{file: "../../test/modules/not_found.toml"},
			want:    Module{},
			wantErr: true,
		},
		{
			name:    "invalid_toml",
			args:    args{file: "../../test/modules/parse_error.toml"},
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
	type fields struct {
		ID          int64
		ModuleType  string
		Name        string
		BackupType  backup_type.BackupType
		BackupPaths []string
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
				ModuleType:  tt.fields.ModuleType,
				Name:        tt.fields.Name,
				BackupType:  tt.fields.BackupType,
				BackupPaths: tt.fields.BackupPaths,
			}
			if got := m.GetLog(); (got == nil) != tt.wantNil {
				t.Errorf("GetLog() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func TestModule_LoadDefaultConfiguration(t *testing.T) {
	MODULES_INSTALL_PATH = "../../test/modules"

	tests := []struct {
		name         string
		mod          Module
		wantedValues map[string]interface{}
		wantErr      bool
	}{
		{
			name: "backup_paths",
			mod: Module{
				ModuleType:  "default_test",
				Name:        "backup_paths",
				BackupPaths: nil,
			},
			wantedValues: map[string]interface{}{
				"BackupPaths": []string{"/tmp/example"},
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

func TestModule_String(t *testing.T) {
	tests := []struct {
		name string
		mod  Module
		want string
	}{
		{
			name: "module",
			mod:  Module{Name: "test"},
			want: "test/default",
		},
		{
			name: "module variant",
			mod:  Module{Name: "test", Variant: "variant1"},
			want: "test/variant1",
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

func TestModule_Valid(t *testing.T) {
	tests := []struct {
		name    string
		mod     Module
		wantErr bool
	}{
		{
			name: "valid_module",
			mod: Module{
				ModuleType: "not_empty",
				Name:       "valid_module",
				BackupType: backup_type.BackupType{Type: backup_type.Diff},
			},
			wantErr: false,
		},
		{
			name: "unknown_backup_type",
			mod: Module{
				ModuleType: "not_empty",
				Name:       "unknown_backup_type",
				BackupType: backup_type.BackupType{Type: backup_type.Unknown},
			},
			wantErr: true,
		},
		{
			name: "missing_name",
			mod: Module{
				ModuleType: "not_empty",
				BackupType: backup_type.BackupType{Type: backup_type.Full},
			},
			wantErr: true,
		},
		{
			name: "missing_module_type",
			mod: Module{
				Name:       "missing_module_type",
				BackupType: backup_type.BackupType{Type: backup_type.Full},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mod.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Valid() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extractArchive(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name:    "unfound_file",
			source:  "/tmp/relique-module-archive-does-not-exist.tar.gz",
			wantErr: true,
		},
		{
			name:    "correct_archive",
			source:  "../../test/modules/relique-module-generic.tar.gz",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpTestFolder, err := os.MkdirTemp("", "relique-test-module-extract-*")
			defer os.RemoveAll(tmpTestFolder)
			if err != nil {
				t.Errorf("extractArchive() cannot create test folder, error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			extractErr := extractArchive(tt.source, tmpTestFolder)
			if (extractErr != nil) != tt.wantErr {
				t.Errorf("extractArchive() error = %v, wantErr %v", extractErr, tt.wantErr)
				return
			}
			if extractErr == nil {
				defaultToml := fmt.Sprintf("%s/default.toml", tmpTestFolder)
				if _, err := os.Lstat(defaultToml); os.IsNotExist(err) {
					t.Errorf("extractArchive() cannot find default.toml from extract archive: %v", defaultToml)
				}
				return
			}
		})
	}
}

func Test_gitClone(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name:    "unfound_git_repo",
			source:  "not_found_relique_repo",
			wantErr: true,
		},
		{
			name:    "correct_repo",
			source:  "github.com/macarrie/relique-module-generic",
			wantErr: false,
		},
		{
			name:    "correct_repo_with_https",
			source:  "http://github.com/macarrie/relique-module-generic",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpTestFolder, err := os.MkdirTemp("", "relique-test-module-git-clone-*")
			defer os.RemoveAll(tmpTestFolder)
			if err != nil {
				t.Errorf("gitClone() cannot create test folder, error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cloneErr := gitClone(tt.source, tmpTestFolder)
			if (cloneErr != nil) != tt.wantErr {
				t.Errorf("gitClone() error = %v, wantErr %v", cloneErr, tt.wantErr)
				return
			}
			if cloneErr == nil {
				defaultToml := fmt.Sprintf("%s/default.toml", tmpTestFolder)
				if _, err := os.Lstat(defaultToml); os.IsNotExist(err) {
					t.Errorf("gitClone() cannot find default.toml from extract archive: %v", defaultToml)
				}
				return
			}
		})
	}
}

func Test_downloadArchive(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "not_found",
			url:     "not found",
			wantErr: true,
		},
		{
			name:    "generic",
			url:     "https://github.com/macarrie/relique-module-generic/releases/download/0.0.1/relique-module-generic.tar.gz",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpTestFolder, err := os.MkdirTemp("", "relique-test-module-download-archive-*")
			defer os.RemoveAll(tmpTestFolder)
			if err != nil {
				t.Errorf("downloadArchive() cannot create test folder, error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			dest := fmt.Sprintf("%s/module.tar.gz", tmpTestFolder)
			if err := downloadArchive(dest, tt.url); (err != nil) != tt.wantErr {
				t.Errorf("downloadArchive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstall(t *testing.T) {
	type args struct {
		path      string
		local     bool
		archive   bool
		force     bool
		skipChown bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remote_git_install",
			args: args{
				path:      "https://github.com/macarrie/relique-module-generic",
				local:     false,
				archive:   false,
				force:     false,
				skipChown: false,
			},
			wantErr: false,
		},
		{
			name: "remote_git_install_404",
			args: args{
				path:      "https://github.com/macarrie/relique-module-doesnotexist",
				local:     false,
				archive:   false,
				force:     false,
				skipChown: false,
			},
			wantErr: true,
		},
		{
			name: "remote_archive_install",
			args: args{
				path:      "https://github.com/macarrie/relique-module-generic/releases/download/0.0.1/relique-module-generic.tar.gz",
				local:     false,
				archive:   true,
				force:     false,
				skipChown: false,
			},
			wantErr: false,
		},
		{
			name: "remote_archive_install_404",
			args: args{
				path:      "https://localhost:8433/archive-does-not-exist.tar.gz",
				local:     false,
				archive:   true,
				force:     false,
				skipChown: false,
			},
			wantErr: true,
		},
		{
			name: "local_archive_install",
			args: args{
				path:      "../../test/modules/relique-module-generic.tar.gz",
				local:     true,
				archive:   true,
				force:     false,
				skipChown: false,
			},
			wantErr: false,
		},
		{
			name: "local_archive_install_404",
			args: args{
				path:      "/tmp/relique-module-archive-does-not-exist.tar.gz",
				local:     true,
				archive:   true,
				force:     false,
				skipChown: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Set custom installed modules
			testInstallFolder, err := os.MkdirTemp("", "relique-module-unittest-install-*")
			if err != nil {
				t.Errorf("Install() error = %v, cannot create temporary install folder", err)
			}
			defer os.RemoveAll(testInstallFolder)
			MODULES_INSTALL_PATH = testInstallFolder

			installErr := Install(MODULES_INSTALL_PATH, tt.args.path, tt.args.local, tt.args.archive, tt.args.force, tt.args.skipChown)
			if (installErr != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", installErr, tt.wantErr)
			}
			if tt.wantErr && installErr != nil {
				// Stop further checks if expected error happened
				return
			}

			installedModules, err := GetLocallyInstalled(MODULES_INSTALL_PATH)
			if err != nil {
				t.Errorf("Install() error = %v", err)
			}
			if len(installedModules) != 1 {
				t.Errorf("Post Install() check. %d modules installed", len(installedModules))
				return
			}
			if installedModules[0].Name != "generic" {
				t.Errorf("Post Install() check. Installed modules name = %s", installedModules[0].Name)
			}
		})
	}
}

func TestGetByName(t *testing.T) {
	type args struct {
		list             []Module
		searchModuleName string
	}
	modList := []Module{
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
		{
			Name: "test3",
		},
	}
	tests := []struct {
		name    string
		args    args
		want    Module
		wantErr bool
	}{
		{
			name: "module_found",
			args: args{
				list:             modList,
				searchModuleName: "test1",
			},
			want:    Module{Name: "test1"},
			wantErr: false,
		},
		{
			name: "module_not_found",
			args: args{
				list:             modList,
				searchModuleName: "totodesbois",
			},
			want:    Module{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetByName(tt.args.list, tt.args.searchModuleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLocallyInstalled(t *testing.T) {
	type args struct {
		modulesInstallPath string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "get_installed",
			args: args{
				modulesInstallPath: "../../test/modules/",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "mod_install_path_does_not_exist",
			args: args{
				modulesInstallPath: "/does_not_exist",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MODULES_INSTALL_PATH = tt.args.modulesInstallPath
			got, err := GetLocallyInstalled(tt.args.modulesInstallPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocallyInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GetLocallyInstalled() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInstalled(t *testing.T) {
	type args struct {
		modulesInstallPath string
		moduleName         string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "installed_module",
			args: args{
				modulesInstallPath: "../../test/modules",
				moduleName:         "default_test",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "mod_not_installed",
			args: args{
				modulesInstallPath: "../../test/modules",
				moduleName:         "not_installed",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "mod_install_path_does_not_exist",
			args: args{
				modulesInstallPath: "/does_not_exist",
				moduleName:         "default_test",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MODULES_INSTALL_PATH = tt.args.modulesInstallPath
			got, err := IsInstalled(tt.args.modulesInstallPath, tt.args.moduleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsInstalled() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		moduleName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove_successful",
			args: args{
				moduleName: "generic",
			},
			wantErr: false,
		},
		{
			name: "remove_not_installed_module",
			args: args{
				moduleName: "not_installed",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpTestFolder, err := os.MkdirTemp("", "relique-test-remove-*")
			defer os.RemoveAll(tmpTestFolder)
			if err != nil {
				t.Errorf("Remove() cannot create test folder, error = %v", err)
				return
			}
			if err := Install(tmpTestFolder, "../../test/modules/relique-module-generic.tar.gz", true, true, false, true); err != nil {
				t.Errorf("Remove() cannot install test module, error = %v", err)
				return
			}

			MODULES_INSTALL_PATH = tmpTestFolder
			removeErr := Remove(tmpTestFolder, tt.args.moduleName)
			if (removeErr != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", removeErr, tt.wantErr)
				return
			}
		})
	}
}
