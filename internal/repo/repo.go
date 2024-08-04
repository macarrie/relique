package repo

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Repository interface {
	GetName() string
	GetType() string
	Write(path string) error
	IsDefault() bool
}

func LoadFromFile(file string) (r Repository, err error) {
	slog.Debug("Loading repository configuration from file", slog.String("path", file))

	f, err := os.Open(file)
	if err != nil {
		return &GenericRepository{}, fmt.Errorf("cannot open file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("cannot close file correctly: %w", err)
		}
	}()

	content, _ := io.ReadAll(f)

	var genericRepo GenericRepository
	if err := toml.Unmarshal(content, &genericRepo); err != nil {
		return &GenericRepository{}, fmt.Errorf("cannot parse toml file: %w", err)
	}

	switch repoType := genericRepo.Type; repoType {
	case "local":
		var localRepo RepositoryLocal
		if err := toml.Unmarshal(content, &localRepo); err != nil {
			return &RepositoryLocal{}, fmt.Errorf("cannot parse toml file: %w", err)
		}
		return &localRepo, nil
	default:
		return &GenericRepository{}, fmt.Errorf("unknown repository type retrieved from file: '%s'", repoType)
	}
}

func LoadFromPath(p string) ([]Repository, error) {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("specified path does not exist or cannot be opened: %w", err)
	}

	var files []string

	_ = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.With(
				slog.Any("error", err),
				slog.String("path", path),
			).Error("Cannot load repository configuration from file")
			return err
		}

		if filepath.Ext(path) == ".toml" {
			files = append(files, path)
		}
		return nil
	})

	var repos []Repository
	for _, file := range files {
		r, err := LoadFromFile(file)
		if err != nil {
			slog.With(
				slog.Any("err", err),
				slog.String("path", file),
			).Error("Cannot load repository configuration from file")
			continue
		}

		repos = append(repos, r)
	}

	return repos, nil
}

func GetByName(list []Repository, name string) (Repository, error) {
	for _, repo := range list {
		if repo.GetName() == name {
			return repo, nil
		}
	}
	return &GenericRepository{}, fmt.Errorf("cannot find repository named '%s'", name)
}

func GetDefault(list []Repository) (Repository, error) {
	for _, repo := range list {
		if repo.IsDefault() {
			return repo, nil
		}
	}
	return &GenericRepository{}, fmt.Errorf("no default repository found in configuration")
}
