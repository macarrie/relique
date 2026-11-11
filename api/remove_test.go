package api

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/job_status"
	"github.com/macarrie/relique/internal/job_type"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
	"github.com/macarrie/relique/internal/utils"
)

func setupRemoveTestEnv(t *testing.T) *repo.RepositoryLocal {
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

func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
}

type testArtifacts struct {
	uuid    string
	catalog string
	storage string
}

func (a testArtifacts) dataFile() string   { return filepath.Join(a.storage, "_data", "file") }
func (a testArtifacts) logFile() string    { return filepath.Join(a.storage, "_logs", "rsync.log") }
func (a testArtifacts) statsFile() string  { return filepath.Join(a.catalog, "stats.toml") }
func (a testArtifacts) moduleFile() string { return filepath.Join(a.catalog, "module.toml") }

// createJob stores a finished job in db with all its files on disk, and its image when withImage is set
func createJob(t *testing.T, r repo.Repository, jobType uint8, done bool, withImage bool) testArtifacts {
	t.Helper()
	mod := module.Module{
		ModuleType:  "generic",
		Name:        "mod",
		BackupType:  backup_type.New(backup_type.Full),
		BackupPaths: []string{"/etc"},
	}
	j := job.NewBackup(client.New("test_client", "127.0.0.1"), mod, r)
	j.JobType = job_type.New(jobType)
	j.Done = done
	j.Status = job_status.New(job_status.Success)
	j.StartTime = time.Now()
	if _, err := j.Save(); err != nil {
		t.Fatal(err)
	}

	catalog := j.GetCatalogPath()
	if err := os.MkdirAll(catalog, 0755); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]interface{}{"module.toml": j.Module, "client.toml": j.Client, "repo.toml": j.Repository} {
		if err := utils.SerializeToFile(data, filepath.Join(catalog, name)); err != nil {
			t.Fatal(err)
		}
	}
	storage, err := j.GetStorageFolderPath()
	if err != nil {
		t.Fatal(err)
	}

	a := testArtifacts{uuid: j.Uuid, catalog: catalog, storage: storage}
	touch(t, a.statsFile())
	touch(t, a.logFile())

	if withImage {
		img := image.New(j.Client, j.Module, j.Repository)
		img.Uuid = j.Uuid
		if _, err := img.Save(); err != nil {
			t.Fatal(err)
		}
		touch(t, a.dataFile())
	}

	return a
}

func assertExists(t *testing.T, path string, want bool) {
	t.Helper()
	_, err := os.Stat(path)
	if got := err == nil; got != want {
		t.Errorf("path %s exists = %v, want %v", path, got, want)
	}
}

func assertRows(t *testing.T, uuid string, wantJob bool, wantImage bool) {
	t.Helper()
	if got, _ := job.Exists(uuid); got != wantJob {
		t.Errorf("job exists = %v, want %v", got, wantJob)
	}
	if got, _ := image.Exists(uuid); got != wantImage {
		t.Errorf("image exists = %v, want %v", got, wantImage)
	}
}

func TestJobRemoveKeepsImage(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, true, true)

	if err := JobRemove(a.uuid, false); err != nil {
		t.Fatalf("JobRemove() error = %v", err)
	}

	assertRows(t, a.uuid, false, true)
	assertExists(t, a.logFile(), false)
	assertExists(t, a.statsFile(), false)
	assertExists(t, a.dataFile(), true)
	assertExists(t, a.moduleFile(), true)
	if _, err := ImageGet(a.uuid); err != nil {
		t.Errorf("image must still be loadable after job removal: %v", err)
	}

	// Removing the remaining image cleans every folder
	if err := ImageRemove(a.uuid, false); err != nil {
		t.Fatalf("ImageRemove() error = %v", err)
	}
	assertRows(t, a.uuid, false, false)
	assertExists(t, a.storage, false)
	assertExists(t, a.catalog, false)
}

func TestJobRemoveCleanImage(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, true, true)
	other := createJob(t, r, job_type.Backup, true, true)

	if err := JobRemove(a.uuid, true); err != nil {
		t.Fatalf("JobRemove() error = %v", err)
	}

	assertRows(t, a.uuid, false, false)
	assertExists(t, a.storage, false)
	assertExists(t, a.catalog, false)

	// Other jobs are untouched
	assertRows(t, other.uuid, true, true)
	assertExists(t, other.dataFile(), true)
	assertExists(t, other.logFile(), true)
}

func TestJobRemoveRestoreJob(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Restore, true, false)

	if err := JobRemove(a.uuid, false); err != nil {
		t.Fatalf("JobRemove() error = %v", err)
	}

	assertRows(t, a.uuid, false, false)
	assertExists(t, a.storage, false)
	assertExists(t, a.catalog, false)
}

func TestJobRemoveNotDone(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, false, false)

	if err := JobRemove(a.uuid, true); err == nil {
		t.Errorf("JobRemove() of a running job should return an error")
	}
	assertRows(t, a.uuid, true, false)
	assertExists(t, a.logFile(), true)
}

func TestJobRemoveUnknown(t *testing.T) {
	setupRemoveTestEnv(t)
	if err := JobRemove("does-not-exist", true); err == nil {
		t.Errorf("JobRemove() of an unknown job should return an error")
	}
}

func TestImageRemoveKeepsJob(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, true, true)

	if err := ImageRemove(a.uuid, false); err != nil {
		t.Fatalf("ImageRemove() error = %v", err)
	}

	assertRows(t, a.uuid, true, false)
	assertExists(t, a.dataFile(), false)
	assertExists(t, filepath.Join(a.storage, "_data"), false)
	assertExists(t, a.logFile(), true)
	assertExists(t, a.statsFile(), true)
	if _, err := JobGet(a.uuid); err != nil {
		t.Errorf("job must still be loadable after image removal: %v", err)
	}
}

func TestImageRemoveCleanJob(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, true, true)

	if err := ImageRemove(a.uuid, true); err != nil {
		t.Fatalf("ImageRemove() error = %v", err)
	}

	assertRows(t, a.uuid, false, false)
	assertExists(t, a.storage, false)
	assertExists(t, a.catalog, false)
}

func TestImageRemoveUnknown(t *testing.T) {
	setupRemoveTestEnv(t)
	if err := ImageRemove("does-not-exist", true); err == nil {
		t.Errorf("ImageRemove() of an unknown image should return an error")
	}
}

func TestSafeRemoveAll(t *testing.T) {
	tmp := t.TempDir()
	uuid := "1234"
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "uuid_folder", path: filepath.Join(tmp, uuid), wantErr: false},
		{name: "inside_uuid_folder", path: filepath.Join(tmp, uuid, "_data"), wantErr: false},
		{name: "unrelated_folder", path: tmp, wantErr: true},
		{name: "root", path: "/", wantErr: true},
		{name: "deep_inside_uuid_folder", path: filepath.Join(tmp, uuid, "a", "b"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := safeRemoveAll(uuid, tt.path); (err != nil) != tt.wantErr {
				t.Errorf("safeRemoveAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
