package image

import (
	"fmt"
	"os"
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
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
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

func testModule(name string) module.Module {
	return module.Module{
		ModuleType:  "generic",
		Name:        name,
		BackupType:  backup_type.BackupType{Type: backup_type.Full},
		BackupPaths: []string{"/tmp/example"},
	}
}

// writeCatalog writes the files GetByUuid needs to rebuild an image from the catalog
func writeCatalog(t *testing.T, img Image) {
	t.Helper()
	catalogPath := img.GetCatalogPath()
	if err := os.MkdirAll(catalogPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := utils.SerializeToFile(img.Module, fmt.Sprintf("%s/module.toml", catalogPath)); err != nil {
		t.Fatal(err)
	}
	if err := utils.SerializeToFile(img.Client, fmt.Sprintf("%s/client.toml", catalogPath)); err != nil {
		t.Fatal(err)
	}
	if err := utils.SerializeToFile(img.Repository, fmt.Sprintf("%s/repo.toml", catalogPath)); err != nil {
		t.Fatal(err)
	}
}

func saveTestImage(t *testing.T, r repo.Repository, moduleName string) Image {
	t.Helper()
	img := New(testClient(), testModule(moduleName), r)
	if _, err := img.Save(); err != nil {
		t.Fatalf("cannot save test image: %v", err)
	}
	writeCatalog(t, img)
	return img
}

func TestNew(t *testing.T) {
	r := repo.RepoLocalNew("local", "/tmp/storage", true)
	img := New(testClient(), testModule("mod"), &r)

	if img.Uuid == "" {
		t.Errorf("New() did not generate a UUID")
	}
	if img.CreatedAt.IsZero() {
		t.Errorf("New() did not set CreatedAt")
	}
	if img.Client.Name != "test_client" || img.Module.Name != "mod" || img.Repository.GetName() != "local" {
		t.Errorf("New() did not set linked objects: %+v", img)
	}
	if other := New(testClient(), testModule("mod"), &r); other.Uuid == img.Uuid {
		t.Errorf("New() generated the same UUID twice")
	}
	if img.GetLog() == nil {
		t.Errorf("GetLog() returned nil")
	}
}

func TestGetSizeOnDisk(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "a"), make([]byte, 100), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "sub", "b"), make([]byte, 250), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want uint64
	}{
		{name: "folder_with_files", path: tmp, want: 350},
		{name: "missing_folder", path: filepath.Join(tmp, "does_not_exist"), want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSizeOnDisk(tt.path)
			if err != nil {
				t.Fatalf("GetSizeOnDisk() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetSizeOnDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillStats(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "a"), make([]byte, 42), 0644); err != nil {
		t.Fatal(err)
	}

	img := Image{}
	stats := rsync_lib.Stats{NumberOfFiles: 10, NumberOfRegularFiles: 7, NumberOfDirectories: 3}
	if err := img.FillStats(stats, tmp); err != nil {
		t.Fatalf("FillStats() error = %v", err)
	}
	if img.SizeOnDisk != 42 || img.NumberOfElements != 10 || img.NumberOfFiles != 7 || img.NumberOfFolders != 3 {
		t.Errorf("FillStats() = %+v", img)
	}
}

func TestPaths(t *testing.T) {
	r := setupTestEnv(t)
	img := New(testClient(), testModule("mod"), r)

	storagePath, err := img.GetStorageFolderPath()
	if err != nil {
		t.Fatalf("GetStorageFolderPath() error = %v", err)
	}
	if want := filepath.Join(r.Path, img.Uuid); storagePath != want {
		t.Errorf("GetStorageFolderPath() = %v, want %v", storagePath, want)
	}
	if !strings.HasSuffix(img.GetCatalogPath(), filepath.Join("catalog", img.Uuid)) {
		t.Errorf("GetCatalogPath() = %v", img.GetCatalogPath())
	}
}

func TestSaveAndGetByUuid(t *testing.T) {
	r := setupTestEnv(t)

	img := New(testClient(), testModule("mod"), r)
	img.NumberOfFiles = 5
	img.SizeOnDisk = 1024

	id, err := img.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if id == 0 || img.ID != id {
		t.Errorf("Save() id = %v, img.ID = %v", id, img.ID)
	}

	// Catalog files missing: image must not be loadable
	if _, err := GetByUuid(img.Uuid); err == nil {
		t.Errorf("GetByUuid() without catalog files should return an error")
	}

	writeCatalog(t, img)
	got, err := GetByUuid(img.Uuid)
	if err != nil {
		t.Fatalf("GetByUuid() error = %v", err)
	}
	if got.ID != id || got.NumberOfFiles != 5 || got.SizeOnDisk != 1024 {
		t.Errorf("GetByUuid() = %+v", got)
	}
	if got.ClientName != "test_client" || got.ModuleName != "mod" || got.RepoName != "local" {
		t.Errorf("GetByUuid() names = %v/%v/%v", got.ClientName, got.ModuleName, got.RepoName)
	}
	if got.Module.Name != "mod" || got.Client.Name != "test_client" || got.Repository.GetName() != "local" {
		t.Errorf("GetByUuid() did not load linked objects from catalog")
	}

	// Saving an image with an ID performs an update
	img.NumberOfFiles = 12
	updatedID, err := img.Save()
	if err != nil {
		t.Fatalf("Save() as update error = %v", err)
	}
	if updatedID != id {
		t.Errorf("Save() as update id = %v, want %v", updatedID, id)
	}
	got, _ = GetByUuid(img.Uuid)
	if got.NumberOfFiles != 12 {
		t.Errorf("Save() as update did not persist changes, got %v", got.NumberOfFiles)
	}
}

func TestGetByUuidNotFound(t *testing.T) {
	setupTestEnv(t)
	if _, err := GetByUuid("does-not-exist"); err == nil {
		t.Errorf("GetByUuid() on unknown uuid should return an error")
	}
}

func TestUpdateUnknownImage(t *testing.T) {
	r := setupTestEnv(t)
	img := New(testClient(), testModule("mod"), r)
	img.ID = 1234

	if _, err := img.Update(nil); err == nil {
		t.Errorf("Update() on image absent from db should return an error")
	}
	if _, err := img.Save(); err == nil {
		t.Errorf("Save() on image with unknown ID should return an error")
	}
}

func TestSearchAndCount(t *testing.T) {
	r := setupTestEnv(t)

	saveTestImage(t, r, "mod_a")
	saveTestImage(t, r, "mod_a")
	last := saveTestImage(t, r, "mod_b")

	// Image present in db but without catalog files must be skipped by search
	orphan := New(testClient(), testModule("mod_a"), r)
	if _, err := orphan.Save(); err != nil {
		t.Fatal(err)
	}

	count, err := Count()
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 4 {
		t.Errorf("Count() = %v, want 4", count)
	}

	tests := []struct {
		name      string
		p         api_helpers.PaginationParams
		s         api_helpers.ImageSearch
		wantCount int
	}{
		{name: "all", wantCount: 3},
		{name: "by_module", s: api_helpers.ImageSearch{ModuleName: "mod_a"}, wantCount: 2},
		{name: "by_client", s: api_helpers.ImageSearch{ClientName: "unknown"}, wantCount: 0},
		{name: "limit", p: api_helpers.PaginationParams{Limit: 1}, wantCount: 0}, // Most recent is the orphan image
		{name: "limit_offset", p: api_helpers.PaginationParams{Limit: 1, Offset: 1}, wantCount: 1},
		{name: "after_future_date", s: api_helpers.ImageSearch{After: time.Now().Add(24 * time.Hour).Format(time.DateTime)}, wantCount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Search(tt.p, tt.s, "")
			if err != nil {
				t.Fatalf("Search() error = %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("Search() returned %d images, want %d", len(got), tt.wantCount)
			}
		})
	}

	got, _ := Search(api_helpers.PaginationParams{}, api_helpers.ImageSearch{}, "")
	if len(got) == 0 || got[0].Uuid != last.Uuid {
		t.Errorf("Search() results are not ordered by most recent first")
	}
}

func TestApplySearchParams(t *testing.T) {
	tests := []struct {
		name     string
		s        api_helpers.ImageSearch
		wantSQL  []string
		wantArgs int
	}{
		{name: "empty", s: api_helpers.ImageSearch{}, wantArgs: 0},
		{
			name:     "all_fields",
			s:        api_helpers.ImageSearch{ModuleName: "m", ClientName: "c", Before: "2024-01-01", After: "2023-01-01"},
			wantSQL:  []string{"module_name = ?", "client_name = ?", "datetime(created_at) < datetime(?)", "datetime(created_at) > datetime(?)"},
			wantArgs: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args, err := ApplySearchParams(sq.Select("uuid").From("images"), tt.s).ToSql()
			if err != nil {
				t.Fatalf("ToSql() error = %v", err)
			}
			for _, want := range tt.wantSQL {
				if !strings.Contains(query, want) {
					t.Errorf("query %q missing %q", query, want)
				}
			}
			if len(args) != tt.wantArgs {
				t.Errorf("got %d args, want %d", len(args), tt.wantArgs)
			}
		})
	}
}
