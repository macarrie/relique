package rsync_task

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/kennygrant/sanitize"

	rsync_lib "github.com/macarrie/relique/internal/rsync_task/lib"
)

type RsyncTask struct {
	Task         rsync_lib.Task
	LogFile      string
	LogErrorFile string
	BackupPath   string
}

func newBackup(source string, destination string, logsRootFolder string, backupPath string, options rsync_lib.RsyncOptions) RsyncTask {
	return RsyncTask{
		Task:         *rsync_lib.NewTask(source, destination, options),
		LogFile:      filepath.Clean(fmt.Sprintf("%s/rsync_log_%s.log", logsRootFolder, sanitize.Accents(sanitize.BaseName(backupPath)))),
		LogErrorFile: filepath.Clean(fmt.Sprintf("%s/rsync_log_error_%s.log", logsRootFolder, sanitize.Accents(sanitize.BaseName(backupPath)))),
		BackupPath:   backupPath,
	}
}

func NewRestore(source string, destination string, logsRootFolder string, backupPath string, exclude []string, excludeCVS bool, include []string) RsyncTask {
	rsyncOptions := rsync_lib.RsyncOptions{
		Archive:      true,
		DelayUpdates: true,
		Mkpath:       true,
		NumericIDs:   true,
		Perms:        true,
		Progress:     true,
		Quiet:        false,
		Recursive:    true,
		Relative:     false,
		Rsh:          "ssh",
		Stats:        true,
		Verbose:      true,
		Exclude:      exclude,
		CVSExclude:   excludeCVS,
		Include:      include,
	}

	return newBackup(source, destination, logsRootFolder, backupPath, rsyncOptions)
}

func NewFullBackup(source string, destination string, logsRootFolder string, backupPath string, exclude []string, excludeCVS bool, include []string) RsyncTask {
	rsyncOptions := rsync_lib.RsyncOptions{
		Archive:      true,
		DelayUpdates: true,
		NumericIDs:   true,
		Perms:        true,
		Progress:     true,
		Quiet:        false,
		Recursive:    true,
		Relative:     true,
		Rsh:          "ssh",
		Stats:        true,
		Verbose:      true,
		Exclude:      exclude,
		CVSExclude:   excludeCVS,
		Include:      include,
	}

	return newBackup(source, destination, logsRootFolder, backupPath, rsyncOptions)
}

func NewDiffBackup(source string, destination string, referencePath string, logsRootFolder string, backupPath string, exclude []string, excludeCVS bool, include []string) RsyncTask {
	rsyncOptions := rsync_lib.RsyncOptions{
		Archive:      true,
		DelayUpdates: true,
		LinkDest:     referencePath,
		NumericIDs:   true,
		Perms:        true,
		Progress:     true,
		Quiet:        false,
		Recursive:    true,
		Relative:     true,
		Rsh:          "ssh",
		Stats:        true,
		Verbose:      true,
		Exclude:      exclude,
		CVSExclude:   excludeCVS,
		Include:      include,
	}

	return newBackup(source, destination, logsRootFolder, backupPath, rsyncOptions)
}

func (t *RsyncTask) GetProgressLog() *slog.Logger {
	state := t.Task.State()
	return slog.With(
		slog.Float64("progress", state.Progress),
		slog.Int("elements_remaining", state.Remain),
		slog.Int("elements_count", state.Total),
		slog.String("transfer_speed", state.Speed),
		slog.String("backup_path", t.BackupPath),
	)
}

func MergeStats(tasks []RsyncTask) rsync_lib.Stats {
	mergedStats := rsync_lib.Stats{}
	for i, _ := range tasks {
		mergedStats.NumberOfFiles += tasks[i].Task.Stats.NumberOfFiles
		mergedStats.NumberOfRegularFiles += tasks[i].Task.Stats.NumberOfRegularFiles
		mergedStats.NumberOfDirectories += tasks[i].Task.Stats.NumberOfDirectories
		mergedStats.NumberOfDeletedFiles += tasks[i].Task.Stats.NumberOfDeletedFiles
		mergedStats.NumberOfCreatedFiles += tasks[i].Task.Stats.NumberOfCreatedFiles
		mergedStats.NumberOfCreatedRegularFiles += tasks[i].Task.Stats.NumberOfCreatedRegularFiles
		mergedStats.TotalFileSize += tasks[i].Task.Stats.TotalFileSize
		mergedStats.TotalTransferredFileSize += tasks[i].Task.Stats.TotalTransferredFileSize
		mergedStats.LiteralData += tasks[i].Task.Stats.LiteralData
		mergedStats.MatchedData += tasks[i].Task.Stats.MatchedData
		mergedStats.FileListSize += tasks[i].Task.Stats.FileListSize
		mergedStats.FileListGenerationTime += tasks[i].Task.Stats.FileListGenerationTime
		mergedStats.FileListTransferTime += tasks[i].Task.Stats.FileListTransferTime
		mergedStats.TotalBytesSent += tasks[i].Task.Stats.TotalBytesSent
		mergedStats.TotalBytesReceived += tasks[i].Task.Stats.TotalBytesReceived
	}

	return mergedStats
}
