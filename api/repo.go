package api

import (
	"fmt"
	"os"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/repo"
	"github.com/samber/lo"
)

func RepoList(p api_helpers.PaginationParams, s api_helpers.RepoSearch) api_helpers.PaginatedResponse[repo.Repository] {
	limit := p.Limit
	repoList := config.Current.Repositories
	// Filters
	if s.RepoType != "" {
		repoList = lo.Filter(repoList, func(item repo.Repository, index int) bool {
			return item.GetType() == s.RepoType
		})
	}

	// Count after filters
	count := len(repoList)
	if limit != 0 {
		repoList = lo.Slice(repoList, 0, int(p.Limit))
	}
	return api_helpers.PaginatedResponse[repo.Repository]{
		Count:      uint64(count),
		Pagination: p,
		Data:       repoList,
	}
}

func RepoGet(name string) (repo.Repository, error) {
	return repo.GetByName(config.Current.Repositories, name)
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
