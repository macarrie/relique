package common

import (
	"reflect"
	"testing"

	"github.com/macarrie/relique/internal/logging"
)

func TestLoad(t *testing.T) {
	logging.Setup(true, "/tmp/relique_unit_tests.log")

	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr bool
	}{
		{
			name:    "file_not_found",
			args:    args{"../../../../test/config/not_found.toml"},
			want:    Configuration{},
			wantErr: true,
		},
		{
			name:    "parse error",
			args:    args{"../../../../test/config/parse_error.toml"},
			want:    Configuration{},
			wantErr: true,
		},
		{
			name: "default values",
			args: args{"../../../../test/config/empty.toml"},
			want: Configuration{
				Version:                   "",
				BindAddr:                  "0.0.0.0",
				PublicAddress:             "localhost",
				Port:                      8433,
				SSLCert:                   "/etc/relique/certs/cert.pem",
				SSLKey:                    "/etc/relique/certs/key.pem",
				StrictSSLCertificateCheck: true,
				ClientCfgPath:             "clients",
				SchedulesCfgPath:          "schedules",
				BackupStoragePath:         "/opt/relique",
				RetentionPath:             "/var/lib/relique/retention.dat",
				ModuleInstallPath:         "/var/lib/relique/modules",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UseFile(tt.args.fileName)
			got, err := Load("")
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr bool
	}{
		{
			name:    "default values",
			args:    args{"../../../../test/config/empty.toml"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UseFile(tt.args.fileName)
			Load("")

			if err := Check(); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
