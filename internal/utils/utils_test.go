package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/consts"
	"github.com/macarrie/relique/internal/repo"
)

func TestGetStoragePath(t *testing.T) {
	localRepo := repo.RepoLocalNew("local_repo", "/tmp/relique_storage", false)
	genericRepo := &repo.GenericRepository{Name: "generic_repo", Type: "unknown"}

	config.Current.Repositories = []repo.Repository{&localRepo, genericRepo}
	defer func() { config.Current.Repositories = nil }()

	tests := []struct {
		name     string
		repo     repo.Repository
		repoName string
		uuid     string
		want     string
		wantErr  bool
	}{
		{
			name:    "local_repo",
			repo:    &localRepo,
			uuid:    "1234",
			want:    "/tmp/relique_storage/1234",
			wantErr: false,
		},
		{
			name:     "repo_from_config",
			repo:     nil,
			repoName: "local_repo",
			uuid:     "5678",
			want:     "/tmp/relique_storage/5678",
			wantErr:  false,
		},
		{
			name:    "nil_repo_no_name",
			repo:    nil,
			uuid:    "1234",
			wantErr: true,
		},
		{
			name:     "nil_repo_unknown_name",
			repo:     nil,
			repoName: "does_not_exist",
			uuid:     "1234",
			wantErr:  true,
		},
		{
			name:    "unimplemented_repo_type",
			repo:    genericRepo,
			uuid:    "1234",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStoragePath(tt.repo, tt.repoName, tt.uuid)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetStoragePath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("GetStoragePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCatalogPath(t *testing.T) {
	tmp := t.TempDir()
	viper.SetConfigFile(filepath.Join(tmp, "relique.toml"))
	config.Current.CatalogPath = "catalog"
	defer func() { config.Current.CatalogPath = "" }()

	want := filepath.Join(tmp, "catalog", "1234")
	if got := GetCatalogPath("1234"); got != want {
		t.Errorf("GetCatalogPath() = %v, want %v", got, want)
	}
}

func TestFormatDatetime(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{name: "zero_time", t: time.Time{}, want: "---"},
		{name: "valid_time", t: time.Date(2024, 3, 5, 7, 8, 9, 0, time.UTC), want: "2024/03/05 07:08:09"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDatetime(tt.t); got != tt.want {
				t.Errorf("FormatDatetime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	if got := FormatDuration(90 * time.Second); got != "1m30s" {
		t.Errorf("FormatDuration() = %v, want %v", got, "1m30s")
	}
}

func TestSerializeToFile(t *testing.T) {
	type sample struct {
		Name  string `toml:"name"`
		Value int    `toml:"value"`
	}

	tmp := t.TempDir()
	tests := []struct {
		name    string
		dest    string
		wantErr bool
	}{
		{name: "valid_destination", dest: filepath.Join(tmp, "out.toml"), wantErr: false},
		{name: "missing_parent_folder", dest: filepath.Join(tmp, "does", "not", "exist.toml"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SerializeToFile(sample{Name: "test", Value: 42}, tt.dest)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SerializeToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			content, readErr := os.ReadFile(tt.dest)
			if readErr != nil {
				t.Fatalf("cannot read serialized file: %v", readErr)
			}
			if !strings.Contains(string(content), `name = "test"`) || !strings.Contains(string(content), "value = 42") {
				t.Errorf("SerializeToFile() unexpected file contents: %s", content)
			}
		})
	}
}

func TestRenderTemplateToMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		tpl          string
		data         interface{}
		wantContains string
		wantErr      bool
	}{
		{
			name:         "simple_template",
			tpl:          "# {{ .Title }}\n\n{{ join .Items \", \" }} - {{ file_size .Size }}",
			data:         map[string]interface{}{"Title": "Hello", "Items": []string{"a", "b"}, "Size": uint64(2048)},
			wantContains: "a, b - 2.0 kB",
			wantErr:      false,
		},
		{
			name:    "execution_error",
			tpl:     "{{ .Missing.Field }}",
			data:    struct{}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderTemplateToMarkdown(tt.name, tt.tpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RenderTemplateToMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("RenderTemplateToMarkdown() = %q, want it to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestGenerateCustomRestorePaths(t *testing.T) {
	tests := []struct {
		name           string
		raw            []string
		modBackupPaths []string
		want           map[string]string
	}{
		{
			name: "no_custom_paths",
			raw:  []string{},
			want: map[string]string{},
		},
		{
			name: "source_and_dest",
			raw:  []string{"/src:/dest"},
			want: map[string]string{"/src": "/dest"},
		},
		{
			name: "source_only",
			raw:  []string{"/src"},
			want: map[string]string{"/src": "/src"},
		},
		{
			name: "multiple_paths",
			raw:  []string{"/a:/b", "/c"},
			want: map[string]string{"/a": "/b", "/c": "/c"},
		},
		{
			name:           "dest_only_restores_all_module_paths",
			raw:            []string{"/a:/b", ":/restore"},
			modBackupPaths: []string{"/etc", "/var/lib"},
			want:           map[string]string{"/etc": "/restore/etc", "/var/lib": "/restore/var/lib"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateCustomRestorePaths(tt.raw, tt.modBackupPaths); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateCustomRestorePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddPaginationParams(t *testing.T) {
	var limit int
	cmd := &cobra.Command{Use: "test"}
	AddPaginationParams(cmd, &limit)

	if limit != consts.DEFAULT_PAGE_SIZE {
		t.Errorf("default limit = %v, want %v", limit, consts.DEFAULT_PAGE_SIZE)
	}
	if err := cmd.Flags().Parse([]string{"-l", "7"}); err != nil {
		t.Fatalf("cannot parse flags: %v", err)
	}
	if limit != 7 {
		t.Errorf("parsed limit = %v, want 7", limit)
	}
}
