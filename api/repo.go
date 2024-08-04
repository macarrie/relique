package api

import (
	"fmt"
	"os"

	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/repo"
)

func RepoList() []repo.Repository {
	return config.Current.Repositories
}

func RepoGet(name string) (repo.Repository, error) {
	return repo.GetByName(RepoList(), name)
}

func RepoCreateLocal(name string, path string, isDefault bool) error {
	// Check if repository name is already taken
	if repo, _ := repo.GetByName(config.Current.Repositories, name); repo.GetName() != "" {
		return fmt.Errorf("a repository of same name already exists ('%s')", repo.GetName())
	}

	// Check if a default repository already exists
	if isDefault {
		if repo, _ := repo.GetDefault(config.Current.Repositories); repo.GetName() != "" {
			return fmt.Errorf("a default repository already exists ('%s')", repo.GetName())
		}
	}

	// Check if path already exists
	if _, err := os.Stat(path); err == nil && !os.IsNotExist(err) {
		return fmt.Errorf("specified folder '%s' already exists, aborting local repository creation to avoid polluting folder", path)
	}

	r := repo.RepoLocalNew(name, path, isDefault)

	fmt.Printf("REPO: %+v\n", r)

	// Save repo to config file
	if err := r.Write(config.GetReposCfgPath()); err != nil {
		return fmt.Errorf("cannot write repository configuration to file: %w", err)
	}

	// Create repo folder
	if err := os.Mkdir(path, 0755); err != nil {
		return fmt.Errorf("cannot create local repository folder '%s': %w", path, err)
	}

	return nil
}
