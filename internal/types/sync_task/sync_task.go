package sync_task

import "github.com/zloylos/grsync"

type SyncTask struct {
	Path  string
	Task  *grsync.Task
	Error error
	Done  bool
}

type SyncTaskProgress struct {
	Path     string
	Progress grsync.State
	Error    string
	Done     bool
}

func ProgressFromSyncTask(s SyncTask) SyncTaskProgress {
	errStr := ""
	if s.Error != nil {
		errStr = s.Error.Error()
	}

	return SyncTaskProgress{
		Path:     s.Path,
		Progress: s.Task.State(),
		Error:    errStr,
		Done:     s.Done,
	}
}
