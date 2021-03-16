package rsync

import (
	"fmt"
	"path/filepath"

	"github.com/macarrie/relique/internal/types/backup_job"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/macarrie/relique/internal/types/config/client_daemon_config"
	"github.com/zloylos/grsync"
)

func getRemoteStorageDest(j *backup_job.BackupJob, path string) string {
	return filepath.Clean(fmt.Sprintf("%s:%s", client_daemon_config.BackupConfig.ServerAddress, j.StorageDestination))
}

func GetSyncTask(j *backup_job.BackupJob, path string) *grsync.Task {
	linkDest := ""
	if j.BackupType.Type == backup_type.Diff {
		linkDest = j.PreviousJobStorageDestination
	}

	rsyncOptions := grsync.RsyncOptions{
		Relative:     true,
		Verbose:      true,
		Archive:      true,
		Recursive:    true,
		Perms:        true,
		WholeFile:    true,
		Rsh:          "ssh",
		DelayUpdates: true,
		NumericIDs:   true,
		Stats:        true,
		Progress:     true,
		LinkDest:     linkDest,
	}

	remoteStorageDest := getRemoteStorageDest(j, path)

	return grsync.NewTask(path, remoteStorageDest, rsyncOptions)
}
