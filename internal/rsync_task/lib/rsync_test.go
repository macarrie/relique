package rsync_lib

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetArguments(t *testing.T) {
	tests := []struct {
		name    string
		options RsyncOptions
		want    []string
	}{
		{
			name:    "no_options",
			options: RsyncOptions{},
			want:    []string{},
		},
		{
			name: "boolean_options",
			options: RsyncOptions{
				Verbose:   true,
				Archive:   true,
				Recursive: true,
				Relative:  true,
				Perms:     true,
				Stats:     true,
			},
			want: []string{"--verbose", "--archive", "--recursive", "--relative", "--perms", "--stats"},
		},
		{
			name: "value_options",
			options: RsyncOptions{
				RsyncPath:      "sudo rsync",
				Rsh:            "ssh",
				Timeout:        10,
				LinkDest:       "/previous",
				CompressLevel:  3,
				SkipCompress:   []string{"gz", "zip"},
				BandwidthLimit: 100,
			},
			want: []string{
				"--rsync-path", "sudo rsync",
				"--rsh", "ssh",
				"--timeout", "10",
				"--link-dest", "/previous",
				"--compress-level", "3",
				"--skip-compress", "gz,zip",
				"--bwlimit", "100",
			},
		},
		{
			name: "include_exclude_filter",
			options: RsyncOptions{
				Include: []string{"*.go"},
				Exclude: []string{"*.tmp", "cache/"},
				Filter:  "- *.bak",
				Chown:   "user:group",
			},
			want: []string{
				"--include=*.go",
				"--exclude=*.tmp",
				"--exclude=cache/",
				"--filter=- *.bak",
				"--chown=user:group",
			},
		},
		{
			name:    "zero_numeric_values_are_ignored",
			options: RsyncOptions{Timeout: 0, BandwidthLimit: 0, CompressLevel: 0},
			want:    []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArguments(tt.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRsync(t *testing.T) {
	tests := []struct {
		name       string
		options    RsyncOptions
		wantBinary string
		wantArgs   []string
	}{
		{
			name:       "default_binary",
			options:    RsyncOptions{Verbose: true},
			wantBinary: "rsync",
			wantArgs:   []string{"rsync", "--verbose", "/src", "/dest"},
		},
		{
			name:       "custom_binary",
			options:    RsyncOptions{RsyncBinaryPath: "/opt/bin/rsync"},
			wantBinary: "/opt/bin/rsync",
			wantArgs:   []string{"/opt/bin/rsync", "/src", "/dest"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRsync("/src", "/dest", tt.options)
			if r.Source != "/src" || r.Destination != "/dest" {
				t.Errorf("NewRsync() source/dest = %v/%v", r.Source, r.Destination)
			}
			if filepath.Base(r.Cmd.Path) != filepath.Base(tt.wantBinary) {
				t.Errorf("NewRsync() binary = %v, want %v", r.Cmd.Path, tt.wantBinary)
			}
			if !reflect.DeepEqual(r.Cmd.Args, tt.wantArgs) {
				t.Errorf("NewRsync() args = %v, want %v", r.Cmd.Args, tt.wantArgs)
			}
		})
	}
}

func TestCreateDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "c")
	if err := createDir(dir); err != nil {
		t.Fatalf("createDir() error = %v", err)
	}
	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		t.Errorf("createDir() did not create folder: %v", err)
	}
}

func TestIsExist(t *testing.T) {
	tmp := t.TempDir()
	if isExist(filepath.Join(tmp, "does_not_exist")) {
		t.Errorf("isExist() = true for missing path")
	}
	if !isExist(tmp) {
		t.Errorf("isExist() = false for existing folder")
	}
	file := filepath.Join(tmp, "file")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	if isExist(file) {
		t.Errorf("isExist() = true for regular file")
	}
}

func TestRsyncRunLocal(t *testing.T) {
	if _, err := exec.LookPath("rsync"); err != nil {
		t.Skip("rsync binary not available")
	}

	src := t.TempDir()
	dest := filepath.Join(t.TempDir(), "dest")
	if err := os.WriteFile(filepath.Join(src, "file.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRsync(src+"/", dest, RsyncOptions{Archive: true})
	if err := r.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	content, err := os.ReadFile(filepath.Join(dest, "file.txt"))
	if err != nil || string(content) != "hello" {
		t.Errorf("Run() did not sync file, content = %q, err = %v", content, err)
	}
}
