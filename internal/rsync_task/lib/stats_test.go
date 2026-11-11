package rsync_lib

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pelletier/go-toml"
)

const sampleRsyncLog = `sending incremental file list

Number of files: 234 (reg: 200, dir: 34)
Number of created files: 56 (reg: 50)
Number of deleted files: 3
Number of regular files transferred: 50
Total file size: 1,234,567 bytes
Total transferred file size: 34,567 bytes
Literal data: 30,000 bytes
Matched data: 4,567 bytes
File list size: 12,345
File list generation time: 0.012 seconds
File list transfer time: 0.003 seconds
Total bytes sent: 45,678
Total bytes received: 1,234

sent 45,678 bytes  received 1,234 bytes  9,382.40 bytes/sec
total size is 1,234,567  speedup is 26.83
`

func TestGetFromRsyncLog(t *testing.T) {
	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "rsync.log")
	if err := os.WriteFile(logPath, []byte(sampleRsyncLog), 0644); err != nil {
		t.Fatal(err)
	}
	emptyLogPath := filepath.Join(tmp, "empty.log")
	if err := os.WriteFile(emptyLogPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		want    Stats
		wantErr bool
	}{
		{
			name: "full_log",
			path: logPath,
			want: Stats{
				NumberOfFiles:               234,
				NumberOfRegularFiles:        200,
				NumberOfDirectories:         34,
				NumberOfCreatedFiles:        56,
				NumberOfCreatedRegularFiles: 50,
				NumberOfDeletedFiles:        3,
				TotalFileSize:               1234567,
				TotalTransferredFileSize:    34567,
				LiteralData:                 30000,
				MatchedData:                 4567,
				FileListSize:                12345,
				FileListGenerationTime:      0.012,
				FileListTransferTime:        0.003,
				TotalBytesSent:              45678,
				TotalBytesReceived:          1234,
				TransferSpeed:               9382.40,
				TransferSpeedup:             26.83,
			},
		},
		{
			name: "empty_log",
			path: emptyLogPath,
			want: Stats{},
		},
		{
			name:    "missing_file",
			path:    filepath.Join(tmp, "does_not_exist.log"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Stats
			err := s.GetFromRsyncLog(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFromRsyncLog() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(s, tt.want) {
				t.Errorf("GetFromRsyncLog() = %+v, want %+v", s, tt.want)
			}
		})
	}
}

// Regression test: rsync prints file counts with thousands separators ("1,234")
func TestGetFromRsyncLogThousandsSeparators(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "rsync.log")
	log := "Number of files: 1,234 (reg: 1,000, dir: 234)\nNumber of created files: 1,050 (reg: 1,001)\nNumber of deleted files: 2,000\n"
	if err := os.WriteFile(logPath, []byte(log), 0644); err != nil {
		t.Fatal(err)
	}

	var s Stats
	if err := s.GetFromRsyncLog(logPath); err != nil {
		t.Fatalf("GetFromRsyncLog() error = %v", err)
	}
	want := Stats{
		NumberOfFiles:               1234,
		NumberOfRegularFiles:        1000,
		NumberOfDirectories:         234,
		NumberOfCreatedFiles:        1050,
		NumberOfCreatedRegularFiles: 1001,
		NumberOfDeletedFiles:        2000,
	}
	if !reflect.DeepEqual(s, want) {
		t.Errorf("GetFromRsyncLog() = %+v, want %+v", s, want)
	}
}

func TestGetFromRsyncLogWithoutDirs(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "rsync.log")
	if err := os.WriteFile(logPath, []byte("Number of files: 3 (reg: 3)\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var s Stats
	if err := s.GetFromRsyncLog(logPath); err != nil {
		t.Fatalf("GetFromRsyncLog() error = %v", err)
	}
	if s.NumberOfFiles != 3 || s.NumberOfRegularFiles != 3 || s.NumberOfDirectories != 0 {
		t.Errorf("GetFromRsyncLog() = %+v", s)
	}
}

func TestStripCommas(t *testing.T) {
	if got := stripCommas("1,234,567"); got != "1234567" {
		t.Errorf("stripCommas() = %v, want 1234567", got)
	}
}

func TestLoadStatsFromFile(t *testing.T) {
	tmp := t.TempDir()

	want := Stats{NumberOfFiles: 10, NumberOfRegularFiles: 8, NumberOfDirectories: 2, TotalFileSize: 4096}
	content, err := toml.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	validPath := filepath.Join(tmp, "stats.toml")
	if err := os.WriteFile(validPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	invalidPath := filepath.Join(tmp, "invalid.toml")
	if err := os.WriteFile(invalidPath, []byte("not = = toml"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		want    Stats
		wantErr bool
	}{
		{name: "valid_file", path: validPath, want: want},
		{name: "parse_error", path: invalidPath, wantErr: true},
		{name: "missing_file", path: filepath.Join(tmp, "missing.toml"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadStatsFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadStatsFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadStatsFromFile() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
