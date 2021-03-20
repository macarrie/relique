package rsync

import (
	"fmt"
	"path/filepath"

	"github.com/macarrie/relique/internal/types/backup_type"

	"github.com/macarrie/relique/internal/types/sync_task"

	"github.com/macarrie/relique/internal/types/config/server_daemon_config"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/macarrie/relique/internal/types/relique_job"

	"github.com/zloylos/grsync"
)

func getClientStorageDest(j *relique_job.ReliqueJob, path string) string {
	return filepath.Clean(fmt.Sprintf("%s:%s", j.Client.Address, path))
}

func GetRestoreSyncTask(j *relique_job.ReliqueJob, path string) sync_task.SyncTask {
	fmt.Printf("JOB: %+v\n", j)
	localRestorePath := fmt.Sprintf("%s/", filepath.Clean(fmt.Sprintf("%s/%s", getJobStorageRoot(&relique_job.ReliqueJob{Uuid: j.RestoreJobUuid}), path)))

	var restoreDestination string
	if j.RestoreDestination == "" {
		restoreDestination = path
	} else {
		restoreDestination = filepath.Clean(fmt.Sprintf("%s/%s", j.RestoreDestination, path))
	}
	remotePath := fmt.Sprintf("%s:%s", j.Client.Address, restoreDestination)

	rsyncOptions := grsync.RsyncOptions{
		Relative:     false,
		Verbose:      true,
		Archive:      true,
		Recursive:    true,
		Perms:        true,
		Rsh:          "ssh",
		DelayUpdates: true,
		NumericIDs:   true,
		Stats:        true,
		Progress:     true,
		Delete:       true,
		DeleteAfter:  true,
	}

	// TODO: Remove
	j.GetLog().WithFields(log.Fields{
		"src":  localRestorePath,
		"dest": remotePath,
	}).Warning("RSYNC PARAMS")

	return sync_task.SyncTask{
		Path:  path,
		Task:  grsync.NewTask(localRestorePath, remotePath, rsyncOptions),
		Error: nil,
	}
}

func getJobStorageRoot(j *relique_job.ReliqueJob) string {
	return filepath.Clean(fmt.Sprintf("%s/%s", server_daemon_config.Config.BackupStoragePath, j.Uuid))
}

func GetBackupSyncTask(j *relique_job.ReliqueJob, path string) sync_task.SyncTask {
	linkDest := ""
	if j.BackupType.Type == backup_type.Diff && j.PreviousJobUuid != "" {
		linkDest = getJobStorageRoot(&relique_job.ReliqueJob{Uuid: j.PreviousJobUuid})
	}

	rsyncOptions := grsync.RsyncOptions{
		Relative:     true,
		Verbose:      true,
		Archive:      true,
		Recursive:    true,
		Perms:        true,
		Rsh:          "ssh",
		DelayUpdates: true,
		NumericIDs:   true,
		Stats:        true,
		Progress:     true,
		LinkDest:     linkDest,
	}

	remoteClientStorageDest := getClientStorageDest(j, path)
	jobStorageRoot := getJobStorageRoot(j)
	return sync_task.SyncTask{
		Path:  path,
		Task:  grsync.NewTask(remoteClientStorageDest, jobStorageRoot, rsyncOptions),
		Error: nil,
	}
}
