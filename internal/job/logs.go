package job

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/macarrie/relique/internal/rsync_task"
)

// JobLog holds the rsync output of one sync task of a job
type JobLog struct {
	// Sanitized name used in the log file names
	Name string `json:"name" toml:"name"`
	// Backup path the logs belong to. Empty if it cannot be matched with a path of the job module
	BackupPath      string `json:"backup_path" toml:"backup_path"`
	Stdout          string `json:"stdout" toml:"stdout"`
	StdoutTruncated bool   `json:"stdout_truncated" toml:"stdout_truncated"`
	Stderr          string `json:"stderr" toml:"stderr"`
	StderrTruncated bool   `json:"stderr_truncated" toml:"stderr_truncated"`
}

// GetLogsFolderPath returns the folder where the rsync logs of the job are written.
// The path is built only from the job uuid and its repository, and must end with <uuid>/_logs
func (j *Job) GetLogsFolderPath() (string, error) {
	if j.Uuid == "" {
		return "", fmt.Errorf("cannot get logs folder of a job without uuid")
	}

	storagePath, err := j.GetStorageFolderPath()
	if err != nil {
		return "", fmt.Errorf("cannot determine job storage folder: %w", err)
	}

	logsPath := filepath.Clean(filepath.Join(storagePath, "_logs"))
	if filepath.Base(logsPath) != "_logs" || filepath.Base(filepath.Dir(logsPath)) != j.Uuid {
		return "", fmt.Errorf("cannot use logs folder '%s': path is not inside a folder named after uuid '%s'", logsPath, j.Uuid)
	}

	return logsPath, nil
}

// GetLogs reads the rsync logs of the job. Only the last maxBytes of each file are returned.
// File names are taken from the job logs folder listing, never from user input.
// An empty list is returned if no log has been written yet (logs are written when each sync ends)
func (j *Job) GetLogs(maxBytes int64) ([]JobLog, error) {
	logsPath, err := j.GetLogsFolderPath()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(logsPath)
	if errors.Is(err, fs.ErrNotExist) {
		return []JobLog{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot list job logs folder: %w", err)
	}

	backupPaths := make(map[string]string)
	for _, p := range j.Module.BackupPaths {
		backupPaths[rsync_task.LogName(p)] = p
	}

	logs := make(map[string]*JobLog)
	for _, entry := range entries {
		// Symlinks and folders are skipped to only serve files written by relique
		if !entry.Type().IsRegular() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, rsync_task.LOG_FILE_SUFFIX) {
			continue
		}

		isError := false
		var name string
		switch {
		case strings.HasPrefix(fileName, rsync_task.LOG_ERROR_FILE_PREFIX):
			isError = true
			name = strings.TrimPrefix(fileName, rsync_task.LOG_ERROR_FILE_PREFIX)
		case strings.HasPrefix(fileName, rsync_task.LOG_FILE_PREFIX):
			name = strings.TrimPrefix(fileName, rsync_task.LOG_FILE_PREFIX)
		default:
			continue
		}
		name = strings.TrimSuffix(name, rsync_task.LOG_FILE_SUFFIX)

		content, truncated, err := readTail(filepath.Join(logsPath, fileName), maxBytes)
		if err != nil {
			return nil, fmt.Errorf("cannot read job log file '%s': %w", fileName, err)
		}

		log, ok := logs[name]
		if !ok {
			log = &JobLog{Name: name, BackupPath: backupPaths[name]}
			logs[name] = log
		}
		if isError {
			log.Stderr = content
			log.StderrTruncated = truncated
		} else {
			log.Stdout = content
			log.StdoutTruncated = truncated
		}
	}

	result := make([]JobLog, 0, len(logs))
	for _, log := range logs {
		result = append(result, *log)
	}
	sort.Slice(result, func(a, b int) bool { return result[a].Name < result[b].Name })

	return result, nil
}

// readTail returns the last maxBytes of a file, and whether the file was longer than that
func readTail(path string, maxBytes int64) (string, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", false, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return "", false, fmt.Errorf("cannot stat file: %w", err)
	}

	truncated := false
	if maxBytes >= 0 && info.Size() > maxBytes {
		truncated = true
		if _, err := f.Seek(-maxBytes, io.SeekEnd); err != nil {
			return "", false, fmt.Errorf("cannot seek file: %w", err)
		}
	}

	var reader io.Reader = f
	if maxBytes >= 0 {
		reader = io.LimitReader(f, maxBytes)
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", false, fmt.Errorf("cannot read file: %w", err)
	}

	return string(content), truncated, nil
}
