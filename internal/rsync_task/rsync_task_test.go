package rsync_task

import (
	"reflect"
	"strings"
	"testing"

	rsync_lib "github.com/macarrie/relique/internal/rsync_task/lib"
)

func argsOf(task *RsyncTask) string {
	return strings.Join(task.Task.Rsync.Cmd.Args, " ")
}

func TestNewFullBackup(t *testing.T) {
	task := NewFullBackup("user@host:/etc", "/storage/_data", "/storage/_logs", "/etc", []string{"*.tmp"}, true, []string{"*.conf"})

	if task.BackupPath != "/etc" {
		t.Errorf("BackupPath = %v, want /etc", task.BackupPath)
	}
	if task.LogFile != "/storage/_logs/rsync_log_-etc.log" {
		t.Errorf("LogFile = %v", task.LogFile)
	}
	if task.LogErrorFile != "/storage/_logs/rsync_log_error_-etc.log" {
		t.Errorf("LogErrorFile = %v", task.LogErrorFile)
	}

	args := argsOf(&task)
	for _, want := range []string{"--relative", "--rsh ssh", "--stats", "--exclude=*.tmp", "--include=*.conf", "--cvs-exclude", "user@host:/etc /storage/_data"} {
		if !strings.Contains(args, want) {
			t.Errorf("full backup args %q missing %q", args, want)
		}
	}
	if strings.Contains(args, "--link-dest") {
		t.Errorf("full backup args %q must not contain --link-dest", args)
	}
}

func TestNewDiffBackup(t *testing.T) {
	task := NewDiffBackup("user@host:/etc", "/storage/_data", "/previous/_data", "/storage/_logs", "/etc", nil, false, nil)

	args := argsOf(&task)
	for _, want := range []string{"--relative", "--link-dest /previous/_data"} {
		if !strings.Contains(args, want) {
			t.Errorf("diff backup args %q missing %q", args, want)
		}
	}
	if strings.Contains(args, "--cvs-exclude") {
		t.Errorf("diff backup args %q must not contain --cvs-exclude", args)
	}
}

func TestNewRestore(t *testing.T) {
	task := NewRestore("/storage/_data/etc", "user@host:/etc", "/storage/_logs", "/etc", nil, false, nil)

	args := argsOf(&task)
	if !strings.Contains(args, "--mkpath") {
		t.Errorf("restore args %q missing --mkpath", args)
	}
	if strings.Contains(args, "--relative") {
		t.Errorf("restore args %q must not contain --relative", args)
	}
}

func TestLogFileNameSanitized(t *testing.T) {
	task := NewFullBackup("src", "dest", "/logs", "/home/user/My Documents/été", nil, false, nil)
	if strings.ContainsAny(strings.TrimPrefix(task.LogFile, "/logs/"), " /é") {
		t.Errorf("LogFile name is not sanitized: %v", task.LogFile)
	}
}

func TestGetProgressLog(t *testing.T) {
	task := NewFullBackup("src", "dest", "/logs", "/etc", nil, false, nil)
	if task.GetProgressLog() == nil {
		t.Errorf("GetProgressLog() returned nil")
	}
}

func TestMergeStats(t *testing.T) {
	newTasksWithStats := func(stats ...rsync_lib.Stats) []RsyncTask {
		tasks := make([]RsyncTask, len(stats))
		for i := range stats {
			tasks[i].Task.Stats = stats[i]
		}
		return tasks
	}

	tests := []struct {
		name  string
		tasks []RsyncTask
		want  rsync_lib.Stats
	}{
		{
			name:  "no_tasks",
			tasks: nil,
			want:  rsync_lib.Stats{},
		},
		{
			name: "multiple_tasks",
			tasks: newTasksWithStats(
				rsync_lib.Stats{NumberOfFiles: 1, NumberOfRegularFiles: 2, TotalFileSize: 100, FileListGenerationTime: 0.5, TotalBytesSent: 10, TransferSpeed: 99},
				rsync_lib.Stats{NumberOfFiles: 3, NumberOfRegularFiles: 4, TotalFileSize: 200, FileListGenerationTime: 0.25, TotalBytesSent: 20, TransferSpeed: 1},
			),
			want: rsync_lib.Stats{
				NumberOfFiles:          4,
				NumberOfRegularFiles:   6,
				TotalFileSize:          300,
				FileListGenerationTime: 0.75,
				TotalBytesSent:         30,
				// TransferSpeed and TransferSpeedup are not merged
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeStats(tt.tasks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeStats() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
