package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml"

	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/repo"
)

func GetStoragePath(r repo.Repository, repoName string, uuid string) (string, error) {
	var repository repo.Repository

	if r == nil {
		if repoName == "" {
			return "", fmt.Errorf("cannot get storage path because used repository is unknown")
		}

		repoFromConfig, err := repo.GetByName(config.Current.Repositories, repoName)
		if err != nil {
			return "", fmt.Errorf("cannot get repository from configuration: %w", err)
		}

		repository = repoFromConfig
	} else {
		repository = r
	}

	if repository.GetType() == "local" {
		localRepo := repository.(*repo.RepositoryLocal)
		return filepath.Clean(fmt.Sprintf("%s/%s/", localRepo.Path, uuid)), nil
	}

	return "", fmt.Errorf("repo type not implemented")
}

func FormatDatetime(t time.Time) string {
	if t.IsZero() {
		return "---"
	}

	return t.Format("2006/01/02 15:04:05")
}

func FormatDuration(d time.Duration) string {
	return d.String()
}

func SerializeToFile[S interface{}](data S, destinationFile string) error {
	tomlContents, serializeErr := toml.Marshal(data)
	if serializeErr != nil {
		return fmt.Errorf("cannot serialize data to toml: %w", serializeErr)
	}
	if err := os.WriteFile(destinationFile, tomlContents, 0644); err != nil {
		return fmt.Errorf("cannot export data to file: %w", err)
	}

	return nil
}
