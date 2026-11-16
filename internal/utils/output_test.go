package utils

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

type outputTestData struct {
	Name     string   `json:"name" toml:"name"`
	SSHPort  int      `json:"ssh_port" toml:"ssh_port"`
	Includes []string `json:"includes" toml:"includes"`
}

func TestOutputFormat_WriteRaw(t *testing.T) {
	data := outputTestData{Name: "client", SSHPort: 22, Includes: []string{"/etc"}}
	tests := []struct {
		name    string
		format  OutputFormat
		want    string
		wantErr bool
	}{
		{
			name:   "json",
			format: OutputFormat{JSON: true},
			want:   "{\n  \"name\": \"client\",\n  \"ssh_port\": 22,\n  \"includes\": [\n    \"/etc\"\n  ]\n}\n",
		},
		{
			name:   "yaml",
			format: OutputFormat{YAML: true},
			want:   "includes:\n    - /etc\nname: client\nssh_port: 22\n",
		},
		{
			name:   "toml",
			format: OutputFormat{TOML: true},
			want:   "name = 'client'\nssh_port = 22\nincludes = ['/etc']\n",
		},
		{
			name:    "no_format",
			format:  OutputFormat{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.format.WriteRaw(&buf, data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WriteRaw() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("WriteRaw() got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOutputFormat_IsRaw(t *testing.T) {
	if (OutputFormat{}).IsRaw() {
		t.Errorf("IsRaw() on empty format = true, want false")
	}
	if !(OutputFormat{YAML: true}).IsRaw() {
		t.Errorf("IsRaw() with yaml = false, want true")
	}
}

func TestAddOutputFormatFlagsExclusive(t *testing.T) {
	var f OutputFormat
	cmd := &cobra.Command{Use: "test", Run: func(cmd *cobra.Command, args []string) {}}
	AddOutputFormatFlags(cmd, &f)
	cmd.SetArgs([]string{"--json", "--yaml"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	if err := cmd.Execute(); err == nil {
		t.Errorf("Execute() with --json and --yaml should fail")
	}
}
