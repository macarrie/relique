package api

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestJobLogsInvalidUuid(t *testing.T) {
	setupRemoveTestEnv(t)

	tests := []string{"", "abc", "../etc", "..%2F..%2Fetc", "../../" + uuid.New().String(), uuid.New().String() + "/../x"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, err := JobLogs(input); !errors.Is(err, ErrInvalidUuid) {
				t.Errorf("JobLogs(%q) error = %v, want ErrInvalidUuid", input, err)
			}
		})
	}
}

func TestJobLogsUnknownJob(t *testing.T) {
	setupRemoveTestEnv(t)

	_, err := JobLogs(uuid.New().String())
	if err == nil || errors.Is(err, ErrInvalidUuid) {
		t.Errorf("JobLogs() on unknown job error = %v, want a not found error", err)
	}
}

func TestJobLogs(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, 1, true, false)
	if err := os.WriteFile(filepath.Join(a.storage, "_logs", "rsync_log_-etc.log"), []byte("etc out"), 0644); err != nil {
		t.Fatal(err)
	}

	logs, err := JobLogs(a.uuid)
	if err != nil {
		t.Fatalf("JobLogs() error = %v", err)
	}
	if len(logs) != 1 || logs[0].BackupPath != "/etc" || logs[0].Stdout != "etc out" {
		t.Errorf("JobLogs() = %+v", logs)
	}
}
