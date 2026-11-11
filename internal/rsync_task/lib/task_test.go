package rsync_lib

import (
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestMatcher(t *testing.T) {
	m := newMatcher(`id=(\d+)`)

	if !m.Match("user id=42") {
		t.Errorf("Match() = false, want true")
	}
	if m.Match("no match here") {
		t.Errorf("Match() = true, want false")
	}
	if got := m.Extract("user id=42"); got != "42" {
		t.Errorf("Extract() = %v, want 42", got)
	}
	if got := m.Extract("no match here"); got != "" {
		t.Errorf("Extract() on no match = %v, want empty string", got)
	}
	if got := m.ExtractAllStringSubmatch("id=1 id=2 id=3", 2); len(got) != 2 || got[1][1] != "2" {
		t.Errorf("ExtractAllStringSubmatch() = %v", got)
	}
}

func TestGetTaskProgress(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantRemain int
		wantTotal  int
	}{
		{name: "valid", input: "12/100", wantRemain: 12, wantTotal: 100},
		{name: "no_separator", input: "12", wantRemain: 0, wantTotal: 0},
		{name: "empty", input: "", wantRemain: 0, wantTotal: 0},
		{name: "not_numbers", input: "a/b", wantRemain: 0, wantTotal: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remain, total := getTaskProgress(tt.input)
			if remain != tt.wantRemain || total != tt.wantTotal {
				t.Errorf("getTaskProgress() = %v, %v, want %v, %v", remain, total, tt.wantRemain, tt.wantTotal)
			}
		})
	}
}

func TestGetTaskSpeed(t *testing.T) {
	tests := []struct {
		name string
		data [][]string
		want string
	}{
		{name: "nil", data: nil, want: ""},
		{name: "single_match", data: [][]string{{"1.00kB/s", "1.00kB/s"}}, want: ""},
		{name: "two_matches", data: [][]string{{"1.00kB/s", "1.00kB/s"}, {"2.00MB/s", "2.00MB/s"}}, want: "2.00MB/s"},
		{name: "short_second_match", data: [][]string{{"a"}, {"b"}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTaskSpeed(tt.data); got != tt.want {
				t.Errorf("getTaskSpeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessStdout(t *testing.T) {
	task := &Task{state: &State{}, log: &Log{}}
	input := strings.Join([]string{
		"sending incremental file list",
		"         32,768 100%   31.25MB/s    0:00:00 (xfr#1, to-chk=25/100)",
		"",
	}, "\n")

	var wg sync.WaitGroup
	wg.Add(1)
	processStdout(&wg, task, strings.NewReader(input))
	wg.Wait()

	state := task.State()
	if state.Remain != 25 || state.Total != 100 {
		t.Errorf("processStdout() remain/total = %v/%v, want 25/100", state.Remain, state.Total)
	}
	if math.Abs(state.Progress-75) > 0.001 {
		t.Errorf("processStdout() progress = %v, want 75", state.Progress)
	}
	if !strings.Contains(task.Log().Stdout, "sending incremental file list") {
		t.Errorf("processStdout() did not record stdout: %q", task.Log().Stdout)
	}
}

func TestProcessStderr(t *testing.T) {
	task := &Task{state: &State{}, log: &Log{}}

	var wg sync.WaitGroup
	wg.Add(1)
	processStderr(&wg, task, strings.NewReader("error line 1\nerror line 2\n"))
	wg.Wait()

	stderr := task.Log().Stderr
	if !strings.Contains(stderr, "error line 1") || !strings.Contains(stderr, "error line 2") {
		t.Errorf("processStderr() stderr = %q", stderr)
	}
}

func TestNewTask(t *testing.T) {
	task := NewTask("/src", "/dest", RsyncOptions{HumanReadable: true})

	args := strings.Join(task.Rsync.Cmd.Args, " ")
	for _, forced := range []string{"--archive", "--partial", "--progress"} {
		if !strings.Contains(args, forced) {
			t.Errorf("NewTask() args %q missing forced option %s", args, forced)
		}
	}
	if strings.Contains(args, "--human-readable") {
		t.Errorf("NewTask() args %q should not contain --human-readable", args)
	}
	if (task.State() != State{}) {
		t.Errorf("NewTask() initial state = %+v, want empty", task.State())
	}
}

func TestTaskRun(t *testing.T) {
	if _, err := exec.LookPath("rsync"); err != nil {
		t.Skip("rsync binary not available")
	}

	src := t.TempDir()
	dest := filepath.Join(t.TempDir(), "dest")
	for _, f := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(src, f), []byte(f), 0644); err != nil {
			t.Fatal(err)
		}
	}

	task := NewTask(src+"/", dest, RsyncOptions{Stats: true})
	if err := task.Run(); err != nil {
		t.Fatalf("Run() error = %v, stderr = %s", err, task.Log().Stderr)
	}
	if !strings.Contains(task.Log().Stdout, "Number of files") {
		t.Errorf("Run() stdout does not contain rsync stats: %q", task.Log().Stdout)
	}
	if _, err := os.Stat(filepath.Join(dest, "a.txt")); err != nil {
		t.Errorf("Run() did not sync files: %v", err)
	}
}

func TestTaskRunMissingBinary(t *testing.T) {
	task := NewTask("/src", filepath.Join(t.TempDir(), "dest"), RsyncOptions{RsyncBinaryPath: "/does/not/exist/rsync"})
	if err := task.Run(); err == nil {
		t.Errorf("Run() with missing binary should return an error")
	}
}
