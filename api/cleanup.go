package api

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"

	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/utils"
)

// cleanupArtifacts removes files left on disk for a job/image uuid once its database rows have been deleted.
// A backup job and the image generated from it share the same uuid, catalog folder and storage folder:
//   - <storage>/_data is owned by the image
//   - <storage>/_logs and <catalog>/stats.toml are owned by the job
//   - <catalog>/{module,client,repo}.toml are needed by both and are removed only when neither exists anymore
func cleanupArtifacts(uuid string, storagePath string) error {
	if uuid == "" {
		return fmt.Errorf("cannot cleanup files without uuid")
	}

	jobExists, err := job.Exists(uuid)
	if err != nil {
		return fmt.Errorf("cannot check if job still exists: %w", err)
	}
	imgExists, err := image.Exists(uuid)
	if err != nil {
		return fmt.Errorf("cannot check if image still exists: %w", err)
	}

	catalogPath := utils.GetCatalogPath(uuid)
	var toRemove []string
	switch {
	case !jobExists && !imgExists:
		toRemove = append(toRemove, catalogPath)
		if storagePath != "" {
			toRemove = append(toRemove, storagePath)
		}
	case !jobExists:
		toRemove = append(toRemove, filepath.Join(catalogPath, "stats.toml"))
		if storagePath != "" {
			toRemove = append(toRemove, filepath.Join(storagePath, "_logs"))
		}
	case !imgExists:
		if storagePath != "" {
			toRemove = append(toRemove, filepath.Join(storagePath, "_data"))
		}
	default:
		return nil
	}

	var result *multierror.Error
	if storagePath == "" {
		result = multierror.Append(result, fmt.Errorf("unknown storage path, data stored in repository has not been removed"))
	}

	for _, path := range toRemove {
		if err := safeRemoveAll(uuid, path); err != nil {
			result = multierror.Append(result, err)
			continue
		}
		slog.With(
			slog.String("uuid", uuid),
			slog.String("path", path),
		).Debug("Removed path from disk")
	}

	return result.ErrorOrNil()
}

// safeRemoveAll only removes paths located inside a folder named after the uuid, to avoid wiping unrelated data
// if a path was computed from incomplete configuration
func safeRemoveAll(uuid string, path string) error {
	path = filepath.Clean(path)
	if filepath.Base(path) != uuid && filepath.Base(filepath.Dir(path)) != uuid {
		return fmt.Errorf("cannot remove '%s': path is not inside a folder named after uuid '%s'", path, uuid)
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("cannot remove '%s': %w", path, err)
	}

	return nil
}
