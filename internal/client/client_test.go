package client

import (
	"reflect"
	"testing"

	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/module"
)

func TestClient_GetLog(t *testing.T) {
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
			want: "client",
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
	module.MODULES_INSTALL_PATH = "../../test/modules"

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
			args: args{path: "../../test/clients/"},
			want: []Client{{
				Name:    "example",
				Address: "example",
				SSHUser: "example",
				SSHPort: 1234,
				Modules: []module.Module{
					{
						ModuleType:  "default_test",
						Name:        "example-diff",
						BackupType:  backup_type.BackupType{Type: backup_type.Diff},
						BackupPaths: []string{"/tmp/example"},
					},
					{
						ModuleType:  "default_test",
						Name:        "example-full",
						BackupType:  backup_type.BackupType{Type: backup_type.Full},
						BackupPaths: []string{"/tmp/example"},
					},
				},
			},
			New("invalid_module", "invalid_module"),
			New("unknown_module", "unknown_module"),
		},
			wantErr: false,
		},
		{
			name:    "path_does_not_exist",
			args:    args{path: "/does_not_exist"},
			want:    nil,
			wantErr: true,
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

func TestLoadFromFile(t *testing.T) {
	module.MODULES_INSTALL_PATH = "../../test/modules"

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
			args: args{file: "../../test/clients/example.toml"},
			want: Client{
				Name:    "example",
				Address: "example",
				SSHUser: "example",
				SSHPort: 1234,
				Modules: []module.Module{
					{
						ModuleType:  "default_test",
						Name:        "example-diff",
						BackupType:  backup_type.BackupType{Type: backup_type.Diff},
						BackupPaths: []string{"/tmp/example"},
					},
					{
						ModuleType:  "default_test",
						Name:        "example-full",
						BackupType:  backup_type.BackupType{Type: backup_type.Full},
						BackupPaths: []string{"/tmp/example"},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "default_values",
			args:    args{file: "../../test/clients/empty.toml"},
			want:    Client{},
			wantErr: false,
		},
		{
			name:    "not_found",
			args:    args{file: "../../test/clients/not_found.toml"},
			want:    Client{},
			wantErr: true,
		},
		{
			name:    "invalid_toml",
			args:    args{file: "../../test/modules/parse_error.toml"},
			want:    Client{},
			wantErr: true,
		},
		{
			name: "module_not_installed",
			args: args{file: "../../test/clients/client_module_unknown.toml"},
			want: New("unknown_module", "unknown_module"),
			wantErr: false,
		},
		{
			name: "module_invalid",
			args: args{file: "../../test/clients/client_module_invalid.toml"},
			want: New("invalid_module", "invalid_module"),
			wantErr: false,
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
