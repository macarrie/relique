package job

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeLogFile(t *testing.T, dir string, name string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestGetLogsFolderPath(t *testing.T) {
	r := setupTestEnv(t)
	j := NewBackup(testClient(), testModule("mod", 1), r)

	got, err := j.GetLogsFolderPath()
	if err != nil {
		t.Fatalf("GetLogsFolderPath() error = %v", err)
	}
	storage, _ := j.GetStorageFolderPath()
	if want := filepath.Join(storage, "_logs"); got != want {
		t.Errorf("GetLogsFolderPath() = %v, want %v", got, want)
	}

	empty := Job{Repository: r}
	if _, err := empty.GetLogsFolderPath(); err == nil {
		t.Errorf("GetLogsFolderPath() on job without uuid should return an error")
	}
}

func TestGetLogs(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		maxBytes int64
		want     []JobLog
	}{
		{
			name:     "no logs folder",
			files:    nil,
			maxBytes: 1024,
			want:     []JobLog{},
		},
		{
			name: "stdout and stderr paired and mapped to backup paths",
			files: map[string]string{
				"rsync_log_-etc.log":         "etc out",
				"rsync_log_error_-etc.log":   "etc err",
				"rsync_log_-var-log.log":     "var out",
				"rsync_log_error_-other.log": "other err",
			},
			maxBytes: 1024,
			want: []JobLog{
				{Name: "-etc", BackupPath: "/etc", Stdout: "etc out", Stderr: "etc err"},
				{Name: "-other", Stderr: "other err"},
				{Name: "-var-log", BackupPath: "/var/log", Stdout: "var out"},
			},
		},
		{
			name: "unrelated files ignored",
			files: map[string]string{
				"rsync_log_-etc.log": "etc out",
				"notes.txt":          "ignored",
				"rsync_log_-etc.txt": "ignored",
				"other.log":          "ignored",
			},
			maxBytes: 1024,
			want: []JobLog{
				{Name: "-etc", BackupPath: "/etc", Stdout: "etc out"},
			},
		},
		{
			name: "large logs truncated to their end",
			files: map[string]string{
				"rsync_log_-etc.log":       "0123456789",
				"rsync_log_error_-etc.log": "short",
			},
			maxBytes: 5,
			want: []JobLog{
				{Name: "-etc", BackupPath: "/etc", Stdout: "56789", StdoutTruncated: true, Stderr: "short"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestEnv(t)
			j := NewBackup(testClient(), testModule("mod", 1, "/etc", "/var/log"), r)
			logsPath, err := j.GetLogsFolderPath()
			if err != nil {
				t.Fatal(err)
			}
			for name, content := range tt.files {
				writeLogFile(t, logsPath, name, content)
			}

			got, err := j.GetLogs(tt.maxBytes)
			if err != nil {
				t.Fatalf("GetLogs() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogs() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetLogsSkipsSymlinks(t *testing.T) {
	r := setupTestEnv(t)
	j := NewBackup(testClient(), testModule("mod", 1, "/etc"), r)
	logsPath, err := j.GetLogsFolderPath()
	if err != nil {
		t.Fatal(err)
	}

	secret := filepath.Join(t.TempDir(), "secret")
	writeLogFile(t, filepath.Dir(secret), "secret", "do not serve")
	writeLogFile(t, logsPath, "rsync_log_-etc.log", "etc out")
	if err := os.Symlink(secret, filepath.Join(logsPath, "rsync_log_error_-etc.log")); err != nil {
		t.Fatal(err)
	}

	got, err := j.GetLogs(1024)
	if err != nil {
		t.Fatalf("GetLogs() error = %v", err)
	}
	for _, l := range got {
		if strings.Contains(l.Stdout+l.Stderr, "do not serve") {
			t.Errorf("GetLogs() served the content of a symlinked file")
		}
	}
}
