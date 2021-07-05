package sync_task

import (
	"github.com/macarrie/relique/internal/lib/rsync"
)

type SyncTask struct {
	Task *rsync.Rsync
}

func New() *SyncTask {
	return &SyncTask{}
}
