package backup_item

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"syscall"

	"github.com/macarrie/relique/internal/types/rsync"

	"github.com/macarrie/relique/internal/types/backup_job"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/macarrie/relique/internal/types/config/server_daemon_config"

	"github.com/pkg/errors"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/monmohan/xferspdy"
)

type BackupItem struct {
	Job            backup_job.BackupJob `json:"job,omitempty"`
	JobUuid        string               `json:"job_uuid"`
	Path           string               `json:"path"`
	Signature      *rsync.Signature     `json:"signature,omitempty"`
	Diff           rsync.Diff           `json:"diff,omitempty"`
	UID            int
	GID            int
	Permissions    os.FileMode
	IsDir          bool
	IsSymlink      bool
	Exists         bool
	CreateHardlink bool
	SymlinkTarget  string
	Checksum       []byte
}

type BackupItemFile struct {
	Item BackupItem            `form:"item" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func New(jobUuid string, path string, info os.FileInfo) *BackupItem {
	var uid, gid int
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)
	}

	isSymlink := info.Mode()&os.ModeSymlink != 0
	var symlinkTarget string
	if isSymlink {
		// TODO: handle error
		symlinkTarget, _ = os.Readlink(path)
	}
	return &BackupItem{
		JobUuid:       jobUuid,
		Path:          path,
		Permissions:   info.Mode(),
		UID:           uid,
		GID:           gid,
		IsDir:         info.IsDir(),
		IsSymlink:     isSymlink,
		SymlinkTarget: symlinkTarget,
		Exists:        true,
	}
}

func (b *BackupItem) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"job_uuid": b.JobUuid,
		"path":     b.Path,
	})
}

func (b *BackupItem) GetSignature() error {
	b.GetLog().Debug("Getting item signature")
	job, err := backup_job.GetByUuid(b.JobUuid)
	if err != nil {
		return errors.Wrap(err, "cannot find backup item associated job")
	}
	if job.Uuid == "" {
		return fmt.Errorf("empty job uuid loaded from db")
	}
	b.Job = job
	sigFilePath, err := getSourceDiffPath(b)
	if err != nil {
		b.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get signature source path. Using /dev/null instead")
		sigFilePath = "/dev/null"
	}

	var fingerprint *xferspdy.Fingerprint
	if _, err := os.Lstat(sigFilePath); os.IsNotExist(err) {
		return errors.Wrap(err, "file does not exist")
	} else {
		fd, err := os.Open(sigFilePath)
		if err != nil {
			return errors.Wrap(err, "cannot open target file")
		}
		fingerprint = xferspdy.NewFingerprintFromReader(fd, 1024)
	}

	b.Signature = rsync.NewSignature(fingerprint)

	return nil
}

func checksum(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return []byte{}, errors.Wrap(err, "cannot open file for checksum computation")
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return []byte{}, errors.Wrap(err, "cannot compute checksum")
	}

	return h.Sum(nil), nil
}

func (b *BackupItem) GetLocalChecksum() ([]byte, error) {
	b.GetLog().Debug("Getting local item checksum")

	h, err := checksum(b.Path)
	if err != nil || len(h) == 0 {
		return []byte{}, errors.Wrap(err, "cannot compute checksum")
	}

	return h, nil
}

func (b *BackupItem) ComputeChecksum() error {
	b.GetLog().Debug("Getting item checksum")
	job, err := backup_job.GetByUuid(b.JobUuid)
	if err != nil {
		return errors.Wrap(err, "cannot find backup item associated job")
	}
	if job.Uuid == "" {
		return fmt.Errorf("empty job uuid loaded from db")
	}
	b.Job = job
	b.Exists = true

	filePath, err := getSourceDiffPath(b)
	if err != nil {
		b.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get checksum source path. Using /dev/null instead")
		b.Exists = false
		return nil
	}

	if _, err := os.Lstat(filePath); os.IsNotExist(err) || filePath == "/dev/null" {
		b.Exists = false
		return nil
	}

	h, err := checksum(filePath)
	if err != nil || len(h) == 0 {
		return errors.Wrap(err, "cannot compute checksum")
	}

	b.Checksum = h

	return nil
}

func (b *BackupItem) GetDiff() {
	b.GetLog().Debug("Getting item diff from signature")

	diff := xferspdy.NewDiff(b.Path, *b.Signature.Sig)
	b.Diff = diff
}

func (b *BackupItem) ApplyDiff() error {
	// TODO: Use bufwriter to avoid loading file in memory
	b.GetLog().Debug("Applying diff")

	job, err := backup_job.GetByUuid(b.JobUuid)
	if err != nil {
		return errors.Wrap(err, "cannot find backup item associated job")
	}
	if job.Uuid == "" {
		return fmt.Errorf("empty job uuid loaded from db")
	}
	b.Job = job

	// TODO: Handle access times
	if b.IsDir {
		b.GetLog().Debug("Backup item is a directory. Creating directory structure")
		if err := createDir(b); err != nil {
			return errors.Wrap(err, "cannot create destination directory")
		}
		return nil
	}

	if b.IsSymlink {
		b.GetLog().Debug("Backup item is a symlink. Creating symlink")
		if err := createSymlink(b); err != nil {
			return errors.Wrap(err, "cannot create symlink")
		}
		return nil
	}

	if job.BackupType.Type == backup_type.Diff {
		if b.CreateHardlink {
			if permissionsMatch(b) {
				b.GetLog().Debug("Backup item does not have any diff. Creating hardlink")
				if err := createHardLink(b); err != nil {
					return errors.Wrap(err, "cannot create hard link with item from previous backup")
				}
				return nil
			} else {
				b.GetLog().Debug("Backup item does not have any diff but permissions differ with previous backup data. Copying file to apply correct permissions")
				if err := copyFromPrevious(b); err != nil {
					return errors.Wrap(err, "cannot copy item from previous backup")
				}
				return nil
			}
		}
	}

	if err := patchFile(b); err != nil {
		return errors.Wrap(err, "cannot patch file")
	}

	return nil
}

func (b *BackupItem) SaveFile(file *multipart.FileHeader) error {
	b.GetLog().Debug("Save uploaded raw file")

	job, err := backup_job.GetByUuid(b.JobUuid)
	if err != nil {
		return errors.Wrap(err, "cannot find backup item associated job")
	}
	if job.Uuid == "" {
		return fmt.Errorf("empty job uuid loaded from db")
	}
	b.Job = job

	// TODO: Handle access times
	if b.IsDir {
		if err := createDir(b); err != nil {
			return errors.Wrap(err, "cannot create destination directory")
		}
		return nil
	}

	if b.IsSymlink {
		if err := createSymlink(b); err != nil {
			return errors.Wrap(err, "cannot create symlink")
		}
		return nil
	}

	tmpPath := getTmpFilePath(b)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return errors.Wrap(err, "cannot create temporary file")
	}
	defer tmpFile.Close()

	srcFile, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "cannot open file for raw file save")
	}
	defer srcFile.Close()

	_, err = io.Copy(tmpFile, srcFile)
	if err != nil {
		return errors.Wrap(err, "cannot copy file from raw send")
	}

	dest := getDestinationBackupPath(b)
	if err := moveTmpDiffFile(tmpPath, dest); err != nil {
		return errors.Wrap(err, "cannot move tmp file to final backup destination")
	}

	// TODO: Apply rights and time settings
	if err := setItemRights(dest, b); err != nil {
		return errors.Wrap(err, "cannot set diffed item permissions")
	}

	return nil
}

func createDir(b *BackupItem) error {
	dest := getDestinationBackupPath(b)
	if err := os.MkdirAll(dest, b.Permissions); err != nil {
		return errors.Wrap(err, "cannot create directory structure")
	}

	if err := setItemRights(dest, b); err != nil {
		return errors.Wrap(err, "cannot set directory item permissions")
	}

	return nil
}

func patchFile(b *BackupItem) error {
	diffFilePath, err := getSourceDiffPath(b)
	if err != nil {
		b.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get diff source path. Using /dev/null instead")
		diffFilePath = "/dev/null"
	}

	src := getTmpFilePath(b)
	outputFile, err := os.OpenFile(src, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return errors.Wrap(err, "cannot create temporary diff file")
	}
	// Apply diff on /dev/null and write it to output file
	if err := xferspdy.PatchFile(b.Diff, diffFilePath, outputFile); err != nil {
		return errors.Wrap(err, "cannot apply diff")
	}
	dest := getDestinationBackupPath(b)
	if err := moveTmpDiffFile(src, dest); err != nil {
		return errors.Wrap(err, "cannot move tmp diff file to final backup destination")
	}

	if err := setItemRights(dest, b); err != nil {
		return errors.Wrap(err, "cannot set diffed item permissions")
	}

	return nil
}

func copyFromPrevious(b *BackupItem) error {
	previousFilePath, err := getSourceDiffPath(b)
	if err != nil {
		return errors.Wrap(err, "cannot get previous backup version path")
	}

	dest := getDestinationBackupPath(b)
	if err := copyFile(previousFilePath, dest); err != nil {
		return errors.Wrap(err, "cannot previous backup file version to final backup destination")
	}

	if err := setItemRights(dest, b); err != nil {
		return errors.Wrap(err, "cannot set item permissions")
	}

	return nil
}

func createHardLink(b *BackupItem) error {
	b.GetLog().Debug("Creating hardlink")

	hardLinkOriginalFile, err := getSourceDiffPath(b)
	if err != nil {
		return errors.Wrap(err, "cannot find previous job hardlink target")
	}

	hardlinkNewTarget := getDestinationBackupPath(b)
	baseDir := filepath.Dir(hardlinkNewTarget)
	if _, err := os.Lstat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return errors.Wrap(err, "cannot create missing directory structure for destination file")
		}
	}

	if err := os.Link(hardLinkOriginalFile, hardlinkNewTarget); err != nil {
		return errors.Wrap(err, "cannot create hard link")
	}

	return nil
}

func createSymlink(b *BackupItem) error {
	b.GetLog().Debug("Creating symlink")

	symlinkPath := getDestinationBackupPath(b)
	var targetPath string
	if filepath.IsAbs(b.SymlinkTarget) {
		targetPath = getDestinationBackupPath(&BackupItem{
			JobUuid: b.JobUuid,
			Path:    b.SymlinkTarget,
		})
	} else {
		targetPath = b.SymlinkTarget
	}

	symlinkPath = filepath.Clean(symlinkPath)
	targetPath = filepath.Clean(targetPath)

	if err := os.Symlink(targetPath, symlinkPath); err != nil {
		return errors.Wrap(err, "cannot create symlink")
	}

	return nil
}

// Perform copy instead of os.Rename because os.Rename fails if src and dest are on different partitions
func moveTmpDiffFile(src string, dest string) error {
	baseDir := filepath.Dir(dest)
	if _, err := os.Lstat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return errors.Wrap(err, "cannot create missing directory structure for destination file")
		}
	}

	inputFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "cannot open source file")
	}
	outputFile, err := os.Create(dest)
	if err != nil {
		inputFile.Close()
		return errors.Wrap(err, "cannot open destination file")
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return errors.Wrap(err, "cannot write to output file")
	}

	// The copy was successful, so now delete the original file
	if err = os.Remove(src); err != nil {
		return errors.Wrap(err, "cannot remove original file")
	}

	return nil
}

func copyFile(src string, dest string) error {
	baseDir := filepath.Dir(dest)
	if _, err := os.Lstat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return errors.Wrap(err, "cannot create missing directory structure for destination file")
		}
	}

	inputFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "cannot open source file")
	}
	outputFile, err := os.Create(dest)
	if err != nil {
		inputFile.Close()
		return errors.Wrap(err, "cannot open destination file")
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return errors.Wrap(err, "cannot write to output file")
	}

	return nil
}

func getSourceDiffPath(b *BackupItem) (string, error) {
	if b.Job.Uuid == "" {
		return "", fmt.Errorf("job with empty uuid")
	}

	if b.Job.BackupType.Type == backup_type.Full {
		return "/dev/null", nil
	}
	previousJob, err := backup_job.GetPreviousJob(b.Job)
	if err != nil {
		return "", errors.Wrap(err, "Cannot find previous job")
	}
	if previousJob.Uuid == "" {
		return "", fmt.Errorf("previous job uuid empty")
	}

	return getDestinationBackupPath(&BackupItem{
		Path:    b.Path,
		JobUuid: previousJob.Uuid,
	}), nil
}

func getTmpFilePath(b *BackupItem) string {
	// TODO: Get a different filepath for each file to avoid concurrency issues
	return filepath.Clean(fmt.Sprintf("%s/relique_%s", os.TempDir(), b.JobUuid))
}

func getDestinationBackupPath(b *BackupItem) string {
	return filepath.Clean(fmt.Sprintf("%s/%s/%s", server_daemon_config.Config.BackupStoragePath, b.JobUuid, b.Path))
}

func permissionsMatch(b *BackupItem) bool {
	sourcePath, err := getSourceDiffPath(b)
	if err != nil {
		b.GetLog().WithFields(log.Fields{
			"err": err,
		}).Debug("Cannot get previous item path for permission check")
		return false
	}

	sourceInfo, err := os.Lstat(sourcePath)
	if err != nil {
		b.GetLog().WithFields(log.Fields{
			"err": err,
		}).Debug("Cannot get previous item file info")
		return false
	}

	if sourceInfo.Mode().Perm() != b.Permissions {
		return false
	}

	return true
}

func setItemRights(dest string, rights *BackupItem) error {
	if err := os.Chmod(dest, rights.Permissions); err != nil {
		return errors.Wrap(err, "cannot set item chmod permissions")
	}
	if err := os.Chown(dest, rights.UID, rights.GID); err != nil {
		return errors.Wrap(err, "cannot set item ownership")
	}

	return nil
}
