package job

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/viper"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job_status"
	"github.com/macarrie/relique/internal/job_type"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
	"github.com/macarrie/relique/internal/rsync_task"
	rsync_lib "github.com/macarrie/relique/internal/rsync_task/lib"
	"github.com/macarrie/relique/internal/utils"
)

// setupTestEnv initializes a temporary config root, catalog, local repository and sqlite database
func setupTestEnv(t *testing.T) *repo.RepositoryLocal {
	t.Helper()
	tmp := t.TempDir()

	viper.SetConfigFile(filepath.Join(tmp, "relique.toml"))
	config.Current.CatalogPath = "catalog"
	localRepo := repo.RepoLocalNew("local", filepath.Join(tmp, "storage"), true)
	config.Current.Repositories = []repo.Repository{&localRepo}

	if err := db.Init(tmp); err != nil {
		t.Fatalf("cannot init test db: %v", err)
	}

	t.Cleanup(func() {
		viper.Reset()
		config.Current = config.Configuration{}
	})

	return &localRepo
}

func testClient() client.Client {
	return client.New("test_client", "127.0.0.1")
}

func testModule(name string, bType uint8, paths ...string) module.Module {
	if len(paths) == 0 {
		paths = []string{"/tmp/example"}
	}
	return module.Module{
		ModuleType:  "generic",
		Name:        name,
		BackupType:  backup_type.New(bType),
		BackupPaths: paths,
	}
}

// writeCatalog writes the files GetByUuid needs to rebuild a job from the catalog
func writeCatalog(t *testing.T, j Job) {
	t.Helper()
	catalogPath := j.GetCatalogPath()
	if err := os.MkdirAll(catalogPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := utils.SerializeToFile(j.Module, fmt.Sprintf("%s/module.toml", catalogPath)); err != nil {
		t.Fatal(err)
	}
	if err := utils.SerializeToFile(j.Client, fmt.Sprintf("%s/client.toml", catalogPath)); err != nil {
		t.Fatal(err)
	}
	if err := utils.SerializeToFile(j.Repository, fmt.Sprintf("%s/repo.toml", catalogPath)); err != nil {
		t.Fatal(err)
	}
}

// saveDoneJob stores a finished backup job in db and catalog
func saveDoneJob(t *testing.T, r repo.Repository, mod module.Module) Job {
	t.Helper()
	j := NewBackup(testClient(), mod, r)
	j.Done = true
	j.Status = job_status.New(job_status.Success)
	j.StartTime = time.Now()
	if _, err := j.Save(); err != nil {
		t.Fatalf("cannot save test job: %v", err)
	}
	writeCatalog(t, j)
	saveImageFor(t, j)
	return j
}

// saveImageFor stores the image generated from a job, sharing its uuid
func saveImageFor(t *testing.T, j Job) image.Image {
	t.Helper()
	img := image.New(j.Client, j.Module, j.Repository)
	img.Uuid = j.Uuid
	if _, err := img.Save(); err != nil {
		t.Fatalf("cannot save test image: %v", err)
	}
	return img
}

func TestNewBackup(t *testing.T) {
	r := repo.RepoLocalNew("local", "/tmp/storage", true)
	mod := testModule("mod", backup_type.Diff)
	j := NewBackup(testClient(), mod, &r)

	if j.Uuid == "" {
		t.Errorf("NewBackup() did not generate a UUID")
	}
	if j.JobType.Type != job_type.Backup {
		t.Errorf("NewBackup() job type = %v, want backup", j.JobType.String())
	}
	if j.Status.Status != job_status.Pending {
		t.Errorf("NewBackup() status = %v, want pending", j.Status.String())
	}
	if j.BackupType.Type != backup_type.Diff {
		t.Errorf("NewBackup() backup type = %v, want module backup type", j.BackupType.String())
	}
	if j.GetLog() == nil {
		t.Errorf("GetLog() returned nil")
	}
}

func TestNewRestore(t *testing.T) {
	r := repo.RepoLocalNew("local", "/tmp/storage", true)
	img := image.New(testClient(), testModule("mod", backup_type.Full), &r)
	target := client.New("target", "10.0.0.1")

	tests := []struct {
		name           string
		restorePaths   map[string]string
		wantModuleName string
	}{
		{name: "module_restore", restorePaths: nil, wantModuleName: "mod"},
		{name: "custom_paths_restore", restorePaths: map[string]string{"/a": "/b"}, wantModuleName: "on-demand"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewRestore(img, target, tt.restorePaths)
			if j.JobType.Type != job_type.Restore {
				t.Errorf("NewRestore() job type = %v, want restore", j.JobType.String())
			}
			if j.RestoreImageUuid != img.Uuid {
				t.Errorf("NewRestore() RestoreImageUuid = %v, want %v", j.RestoreImageUuid, img.Uuid)
			}
			if j.Client.Name != "target" {
				t.Errorf("NewRestore() client = %v, want target", j.Client.Name)
			}
			if j.Module.Name != tt.wantModuleName {
				t.Errorf("NewRestore() module name = %v, want %v", j.Module.Name, tt.wantModuleName)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	start := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		job  Job
		want time.Duration
	}{
		{name: "not_started", job: Job{}, want: 0},
		{name: "finished", job: Job{StartTime: start, EndTime: start.Add(90*time.Second + 500*time.Millisecond)}, want: 90 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.job.Duration(); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}

	running := Job{StartTime: time.Now().Add(-5 * time.Second)}
	if got := running.Duration(); got < 4*time.Second || got > 6*time.Second {
		t.Errorf("Duration() on running job = %v, want ~5s", got)
	}
}

func TestSaveAndGetByUuid(t *testing.T) {
	r := setupTestEnv(t)

	j := NewBackup(testClient(), testModule("mod", backup_type.Full), r)
	id, err := j.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if id == 0 || j.ID != id {
		t.Errorf("Save() id = %v, j.ID = %v", id, j.ID)
	}

	if _, err := GetByUuid(j.Uuid); err == nil {
		t.Errorf("GetByUuid() without catalog files should return an error")
	}

	writeCatalog(t, j)
	got, err := GetByUuid(j.Uuid)
	if err != nil {
		t.Fatalf("GetByUuid() error = %v", err)
	}
	if got.ID != id || got.Status.Status != job_status.Pending || got.JobType.Type != job_type.Backup || got.Done {
		t.Errorf("GetByUuid() = %+v", got)
	}
	if got.Module.Name != "mod" || got.Client.Name != "test_client" || got.Repository.GetName() != "local" {
		t.Errorf("GetByUuid() did not load linked objects from catalog")
	}

	// Update through Save
	j.Status = job_status.New(job_status.Success)
	j.Done = true
	if _, err := j.Save(); err != nil {
		t.Fatalf("Save() as update error = %v", err)
	}
	got, _ = GetByUuid(j.Uuid)
	if got.Status.Status != job_status.Success || !got.Done {
		t.Errorf("Save() as update did not persist changes: %+v", got)
	}
}

func TestGetByUuidNotFound(t *testing.T) {
	setupTestEnv(t)
	if _, err := GetByUuid("does-not-exist"); err == nil {
		t.Errorf("GetByUuid() on unknown uuid should return an error")
	}
}

func TestUpdateUnknownJob(t *testing.T) {
	r := setupTestEnv(t)
	j := NewBackup(testClient(), testModule("mod", backup_type.Full), r)
	j.ID = 1234

	if _, err := j.Update(nil); err == nil {
		t.Errorf("Update() on job absent from db should return an error")
	}
	if _, err := j.Save(); err == nil {
		t.Errorf("Save() on job with unknown ID should return an error")
	}
}

func TestSearchAndCount(t *testing.T) {
	r := setupTestEnv(t)

	saveDoneJob(t, r, testModule("mod_a", backup_type.Full))
	saveDoneJob(t, r, testModule("mod_a", backup_type.Diff))
	saveDoneJob(t, r, testModule("mod_b", backup_type.Full))

	tests := []struct {
		name      string
		p         api_helpers.PaginationParams
		s         api_helpers.JobSearch
		wantCount int
	}{
		{name: "all", wantCount: 3},
		{name: "by_module", s: api_helpers.JobSearch{ModuleName: "mod_a"}, wantCount: 2},
		{name: "by_backup_type", s: api_helpers.JobSearch{BackupType: backup_type.Diff}, wantCount: 1},
		{name: "by_status", s: api_helpers.JobSearch{Status: job_status.Error}, wantCount: 0},
		{name: "by_job_type", s: api_helpers.JobSearch{JobType: job_type.Backup}, wantCount: 3},
		{name: "by_client", s: api_helpers.JobSearch{ClientName: "unknown"}, wantCount: 0},
		{name: "limit_offset", p: api_helpers.PaginationParams{Limit: 2, Offset: 2}, wantCount: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Search(tt.p, tt.s)
			if err != nil {
				t.Fatalf("Search() error = %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("Search() returned %d jobs, want %d", len(got), tt.wantCount)
			}

			count, err := Count(tt.s)
			if err != nil {
				t.Fatalf("Count() error = %v", err)
			}
			if tt.p.Limit == 0 && count != uint64(tt.wantCount) {
				t.Errorf("Count() = %d, want %d", count, tt.wantCount)
			}
		})
	}
}

func TestApplySearchParams(t *testing.T) {
	s := api_helpers.JobSearch{
		ModuleName: "m",
		ClientName: "c",
		BackupType: backup_type.Full,
		JobType:    job_type.Backup,
		Status:     job_status.Success,
		Before:     "2024-01-01",
		After:      "2023-01-01",
	}
	query, args, err := ApplySearchParams(sq.Select("uuid").From("jobs"), s).ToSql()
	if err != nil {
		t.Fatalf("ToSql() error = %v", err)
	}
	for _, want := range []string{"module_name = ?", "client_name = ?", "backup_type = ?", "job_type = ?", "status = ?", "datetime(start_time) < datetime(?)", "datetime(start_time) > datetime(?)"} {
		if !strings.Contains(query, want) {
			t.Errorf("query %q missing %q", query, want)
		}
	}
	if len(args) != 7 {
		t.Errorf("got %d args, want 7", len(args))
	}
}

func TestGetPrevious(t *testing.T) {
	r := setupTestEnv(t)
	mod := testModule("mod", backup_type.Full)

	current := NewBackup(testClient(), mod, r)
	if _, err := GetPrevious(current, backup_type.New(backup_type.Full)); err == nil {
		t.Errorf("GetPrevious() with no previous job should return an error")
	}

	saveDoneJob(t, r, mod)
	latest := saveDoneJob(t, r, mod)
	// Not done jobs must be ignored
	pending := NewBackup(testClient(), mod, r)
	if _, err := pending.Save(); err != nil {
		t.Fatal(err)
	}

	got, err := GetPrevious(current, backup_type.New(backup_type.Full))
	if err != nil {
		t.Fatalf("GetPrevious() error = %v", err)
	}
	if got.Uuid != latest.Uuid {
		t.Errorf("GetPrevious() = %v, want latest done job %v", got.Uuid, latest.Uuid)
	}

	if _, err := GetPrevious(current, backup_type.New(backup_type.Diff)); err == nil {
		t.Errorf("GetPrevious() for diff with no previous diff job should return an error")
	}

	otherModule := NewBackup(testClient(), testModule("other", backup_type.Full), r)
	if _, err := GetPrevious(otherModule, backup_type.New(backup_type.Full)); err == nil {
		t.Errorf("GetPrevious() must filter on module name")
	}
}

func TestSetupBackupFull(t *testing.T) {
	r := setupTestEnv(t)
	j := NewBackup(testClient(), testModule("mod", backup_type.Full, "/etc", "/var/lib"), r)

	if err := j.SetupBackup(); err != nil {
		t.Fatalf("SetupBackup() error = %v", err)
	}
	if j.Status.Status != job_status.Active || j.StartTime.IsZero() {
		t.Errorf("SetupBackup() did not mark job as active")
	}
	if j.Client.SSHUser != client.DEFAULT_SSH_USER {
		t.Errorf("SetupBackup() SSHUser = %v, want default", j.Client.SSHUser)
	}
	if len(j.Tasks) != 2 {
		t.Fatalf("SetupBackup() created %d tasks, want 2", len(j.Tasks))
	}

	storagePath, _ := j.GetStorageFolderPath()
	for _, sub := range []string{"_data", "_logs"} {
		if _, err := os.Stat(filepath.Join(storagePath, sub)); err != nil {
			t.Errorf("SetupBackup() did not create %s folder: %v", sub, err)
		}
	}
	for _, f := range []string{"module.toml", "client.toml", "repo.toml"} {
		if _, err := os.Stat(filepath.Join(j.GetCatalogPath(), f)); err != nil {
			t.Errorf("SetupBackup() did not write catalog file %s: %v", f, err)
		}
	}

	args := strings.Join(j.Tasks[0].Task.Rsync.Cmd.Args, " ")
	if !strings.Contains(args, fmt.Sprintf("%s@127.0.0.1:/etc", client.DEFAULT_SSH_USER)) {
		t.Errorf("SetupBackup() task source is wrong: %v", args)
	}

	if _, err := GetByUuid(j.Uuid); err != nil {
		t.Errorf("job not retrievable after SetupBackup(): %v", err)
	}
}

func TestSetupBackupSSHPort(t *testing.T) {
	r := setupTestEnv(t)

	tests := []struct {
		name     string
		port     int
		wantPort int
	}{
		{name: "default_port", port: 0, wantPort: client.DEFAULT_SSH_PORT},
		{name: "custom_port_kept", port: 2222, wantPort: 2222},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := testClient()
			cl.SSHPort = tt.port
			j := NewBackup(cl, testModule("mod", backup_type.Full), r)
			if err := j.SetupBackup(); err != nil {
				t.Fatalf("SetupBackup() error = %v", err)
			}
			if j.Client.SSHPort != tt.wantPort {
				t.Errorf("SetupBackup() SSHPort = %v, want %v", j.Client.SSHPort, tt.wantPort)
			}
		})
	}
}

func TestSetupBackupDiff(t *testing.T) {
	r := setupTestEnv(t)
	mod := testModule("mod", backup_type.Diff)

	// No previous job: diff backup falls back to full
	first := NewBackup(testClient(), mod, r)
	if err := first.SetupBackup(); err != nil {
		t.Fatalf("SetupBackup() error = %v", err)
	}
	if first.BackupType.Type != backup_type.Full {
		t.Errorf("SetupBackup() without previous job backup type = %v, want full", first.BackupType.String())
	}
	if first.PreviousJobUuid != "" {
		t.Errorf("SetupBackup() without previous job set PreviousJobUuid = %v", first.PreviousJobUuid)
	}

	first.Done = true
	if _, err := first.Save(); err != nil {
		t.Fatal(err)
	}
	saveImageFor(t, first)

	// Previous full job exists: diff is computed from it
	second := NewBackup(testClient(), mod, r)
	if err := second.SetupBackup(); err != nil {
		t.Fatalf("SetupBackup() error = %v", err)
	}
	if second.BackupType.Type != backup_type.Diff {
		t.Errorf("SetupBackup() backup type = %v, want diff", second.BackupType.String())
	}
	if second.PreviousJobUuid != first.Uuid {
		t.Errorf("SetupBackup() PreviousJobUuid = %v, want %v", second.PreviousJobUuid, first.Uuid)
	}
	firstStorage, _ := first.GetStorageFolderPath()
	args := strings.Join(second.Tasks[0].Task.Rsync.Cmd.Args, " ")
	if !strings.Contains(args, "--link-dest "+filepath.Join(firstStorage, "_data")) {
		t.Errorf("diff task does not link to previous job data: %v", args)
	}

	second.Done = true
	if _, err := second.Save(); err != nil {
		t.Fatal(err)
	}
	saveImageFor(t, second)

	// Previous diff job exists: diff is computed from latest diff
	third := NewBackup(testClient(), mod, r)
	if err := third.SetupBackup(); err != nil {
		t.Fatalf("SetupBackup() error = %v", err)
	}
	if third.PreviousJobUuid != second.Uuid {
		t.Errorf("SetupBackup() PreviousJobUuid = %v, want previous diff %v", third.PreviousJobUuid, second.Uuid)
	}
}

func TestSetupRestore(t *testing.T) {
	r := setupTestEnv(t)
	source := saveDoneJob(t, r, testModule("mod", backup_type.Full, "/etc", "/home"))
	sourceImg := image.New(source.Client, source.Module, source.Repository)
	sourceImg.Uuid = source.Uuid

	t.Run("missing_image_uuid", func(t *testing.T) {
		j := NewRestore(sourceImg, testClient(), nil)
		j.RestoreImageUuid = ""
		if err := j.SetupRestore(); err == nil {
			t.Errorf("SetupRestore() without image uuid should return an error")
		}
	})

	t.Run("unknown_image_uuid", func(t *testing.T) {
		j := NewRestore(sourceImg, testClient(), nil)
		j.RestoreImageUuid = "does-not-exist"
		if err := j.SetupRestore(); err == nil {
			t.Errorf("SetupRestore() with unknown image uuid should return an error")
		}
	})

	t.Run("module_paths", func(t *testing.T) {
		j := NewRestore(sourceImg, testClient(), nil)
		if err := j.SetupRestore(); err != nil {
			t.Fatalf("SetupRestore() error = %v", err)
		}
		if len(j.Tasks) != 2 {
			t.Errorf("SetupRestore() created %d tasks, want 2", len(j.Tasks))
		}
		if _, err := GetByUuid(j.Uuid); err != nil {
			t.Errorf("job not retrievable after SetupRestore(): %v", err)
		}
	})

	t.Run("custom_paths", func(t *testing.T) {
		j := NewRestore(sourceImg, testClient(), map[string]string{"/etc": "/restore/etc"})
		if err := j.SetupRestore(); err != nil {
			t.Fatalf("SetupRestore() error = %v", err)
		}
		if len(j.Tasks) != 1 {
			t.Fatalf("SetupRestore() created %d tasks, want 1", len(j.Tasks))
		}
		args := strings.Join(j.Tasks[0].Task.Rsync.Cmd.Args, " ")
		if !strings.Contains(args, "127.0.0.1:/restore/etc") {
			t.Errorf("SetupRestore() task destination is wrong: %v", args)
		}
	})
}

func TestStart(t *testing.T) {
	if _, err := exec.LookPath("rsync"); err != nil {
		t.Skip("rsync binary not available")
	}

	tests := []struct {
		name       string
		srcExists  bool
		binary     string
		wantStatus uint8
		wantImage  bool
	}{
		{name: "successful_sync", srcExists: true, wantStatus: job_status.Success, wantImage: true},
		// rsync exits with code 23 (partial transfer) when the source does not exist
		{name: "partial_sync", srcExists: false, wantStatus: job_status.Incomplete, wantImage: true},
		// `false` exits with code 1, which is not a partial transfer code
		{name: "failed_sync", srcExists: true, binary: "false", wantStatus: job_status.Error, wantImage: false},
		{name: "rsync_cannot_start", srcExists: true, binary: "/does/not/exist/rsync", wantStatus: job_status.Error, wantImage: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestEnv(t)

			src := filepath.Join(t.TempDir(), "src")
			if tt.srcExists {
				if err := os.MkdirAll(src, 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(src, "file.txt"), []byte("hello"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			j := NewBackup(testClient(), testModule("mod", backup_type.Full, src), r)
			if err := j.SetupBackup(); err != nil {
				t.Fatalf("SetupBackup() error = %v", err)
			}

			// Replace remote tasks by a local sync to avoid needing SSH access
			storagePath, _ := j.GetStorageFolderPath()
			j.Tasks = []rsync_task.RsyncTask{
				rsync_task.NewFullBackup(src, filepath.Join(storagePath, "_data"), filepath.Join(storagePath, "_logs"), src, nil, false, nil),
			}
			if tt.binary != "" {
				j.Tasks[0].Task.Rsync = rsync_lib.NewRsync(src, filepath.Join(storagePath, "_data"), rsync_lib.RsyncOptions{RsyncBinaryPath: tt.binary})
			}

			if err := j.Start(); err != nil {
				t.Fatalf("Start() error = %v", err)
			}
			if !j.Done || j.EndTime.IsZero() {
				t.Errorf("Start() did not mark job as done")
			}
			if j.Status.Status != tt.wantStatus {
				want := job_status.New(tt.wantStatus)
				t.Errorf("Start() status = %v, want %v", j.Status.String(), want.String())
			}
			if _, err := os.Stat(filepath.Join(j.GetCatalogPath(), "stats.toml")); err != nil {
				t.Errorf("Start() did not write stats file: %v", err)
			}
			if _, err := os.Stat(j.Tasks[0].LogFile); err != nil {
				t.Errorf("Start() did not write task log file: %v", err)
			}

			img, err := image.GetByUuid(j.Uuid)
			if !tt.wantImage {
				if err == nil {
					t.Errorf("Start() generated an image for a failed job")
				}
				return
			}
			if err != nil {
				t.Fatalf("Start() did not generate image: %v", err)
			}
			if tt.srcExists {
				if img.NumberOfFiles != 1 {
					t.Errorf("image NumberOfFiles = %v, want 1", img.NumberOfFiles)
				}
				if _, err := os.Stat(filepath.Join(storagePath, "_data", src, "file.txt")); err != nil {
					t.Errorf("file not synced to storage: %v", err)
				}
			}
		})
	}
}

func TestDeleteAndExists(t *testing.T) {
	r := setupTestEnv(t)
	j := saveDoneJob(t, r, testModule("mod", backup_type.Full))

	exists, err := Exists(j.Uuid)
	if err != nil || !exists {
		t.Fatalf("Exists() before delete = %v, %v, want true", exists, err)
	}

	if err := j.Delete(); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	exists, err = Exists(j.Uuid)
	if err != nil || exists {
		t.Errorf("Exists() after delete = %v, %v, want false", exists, err)
	}
	// Image sharing the job uuid is not removed with the job
	if imgExists, _ := image.Exists(j.Uuid); !imgExists {
		t.Errorf("Delete() of a job must not delete its image")
	}

	unknown := Job{Uuid: "does-not-exist"}
	if err := unknown.Delete(); err == nil {
		t.Errorf("Delete() of an unknown job should return an error")
	}
}

func TestGetPreviousSkipsJobsWithoutImage(t *testing.T) {
	r := setupTestEnv(t)
	mod := testModule("mod", backup_type.Full)

	older := saveDoneJob(t, r, mod)
	latest := saveDoneJob(t, r, mod)
	img := image.Image{Uuid: latest.Uuid}
	if err := img.Delete(); err != nil {
		t.Fatal(err)
	}

	current := NewBackup(testClient(), mod, r)
	got, err := GetPrevious(current, backup_type.New(backup_type.Full))
	if err != nil {
		t.Fatalf("GetPrevious() error = %v", err)
	}
	if got.Uuid != older.Uuid {
		t.Errorf("GetPrevious() = %v, want %v (latest job has no image anymore)", got.Uuid, older.Uuid)
	}
}

func TestSetupRestoreWithoutSourceJob(t *testing.T) {
	r := setupTestEnv(t)
	source := saveDoneJob(t, r, testModule("mod", backup_type.Full, "/etc"))
	if err := source.Delete(); err != nil {
		t.Fatal(err)
	}

	sourceImg, err := image.GetByUuid(source.Uuid)
	if err != nil {
		t.Fatal(err)
	}
	j := NewRestore(sourceImg, testClient(), nil)
	if err := j.SetupRestore(); err != nil {
		t.Fatalf("SetupRestore() from an image whose job was removed error = %v", err)
	}
	if len(j.Tasks) != 1 {
		t.Errorf("SetupRestore() created %d tasks, want 1", len(j.Tasks))
	}
}
