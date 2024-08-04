package repo

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/pelletier/go-toml"
)

type RepositoryLocal struct {
	Name    string `json:"name" toml:"name"`
	Type    string `json:"type" toml:"type"`
	Path    string `json:"path" toml:"path"`
	Default bool   `json:"default" toml:"default"`
}

func RepoLocalNew(name string, path string, isDefault bool) RepositoryLocal {
	return RepositoryLocal{
		Name:    name,
		Type:    "local",
		Path:    path,
		Default: isDefault,
	}
}

func (r *RepositoryLocal) GetName() string {
	return r.Name
}

func (r *RepositoryLocal) GetType() string {
	return r.Type
}

func (r *RepositoryLocal) GetLog() *slog.Logger {
	return slog.With(
		slog.String("name", r.GetName()),
		slog.String("type", r.GetType()),
		slog.String("path", r.Path),
		slog.Bool("default", r.IsDefault()),
	)
}

func (r *RepositoryLocal) Write(rootPath string) error {
	var path string = filepath.Clean(fmt.Sprintf("%s/%s.toml",
		rootPath,
		strings.ToLower(sanitize.Accents(sanitize.BaseName(r.GetName()))),
	))

	repoToml, repoErr := toml.Marshal(r)
	if repoErr != nil {
		return fmt.Errorf("cannot serialize repository info to toml data: %w", repoErr)
	}
	if err := os.WriteFile(path, repoToml, 0644); err != nil {
		return fmt.Errorf("cannot export repository info to file: %w", err)
	}

	r.GetLog().With(
		slog.String("path", path),
	).Debug("Saved repository to file")

	return nil
}

func (r *RepositoryLocal) IsDefault() bool {
	return r.Default
}
