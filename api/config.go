package api

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/module"
)

func ConfigGet() (config.Configuration, error) {
	if !config.Loaded {
		if err := config.Load("relique"); err != nil {
			return config.Configuration{}, fmt.Errorf("cannot load relique configuration: %w", err)
		}
	}
	return config.Current, nil
}

func ConfigInit(cfgPath string, modPath string, repoPath string, catalogPath string) error {
	configPath := cfgPath
	moduleInstallPath := modPath
	repoStoragePath := repoPath
	catalogStoragePath := catalogPath
	if cfgPath == "" {
		configPath = "/etc/relique"
		if modPath == "" {
			moduleInstallPath = "/var/lib/relique/modules"
		}
		if repoPath == "" {
			repoStoragePath = "/var/lib/relique/storage"
		}
		if catalogPath == "" {
			catalogStoragePath = "/var/lib/relique/catalog"
		}
	}
	if moduleInstallPath == "" {
		moduleInstallPath = filepath.Clean(fmt.Sprintf("%s/modules", configPath))
	}
	if repoStoragePath == "" {
		repoStoragePath = filepath.Clean(fmt.Sprintf("%s/storage", configPath))
	}
	if catalogStoragePath == "" {
		catalogStoragePath = filepath.Clean(fmt.Sprintf("%s/catalog", configPath))
	}

	configPath = filepath.Clean(configPath)
	moduleInstallPath = filepath.Clean(moduleInstallPath)

	// Check if config folder already exists
	if _, err := os.Stat(configPath); err == nil && !os.IsNotExist(err) {
		return fmt.Errorf("specified folder '%s' already exists, aborting config init to avoid overwriting existing configuration", configPath)
	}

	// Check if module install folder already exists
	if _, err := os.Stat(moduleInstallPath); err == nil && !os.IsNotExist(err) {
		return fmt.Errorf("specified folder '%s' already exists, aborting config init to avoid overwriting existing module install folder", moduleInstallPath)
	}

	// Create config folder
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("cannot create folder '%s': %w", configPath, err)
	}
	slog.With(
		slog.String("path", configPath),
	).Info("Created default relique configuration folder")

	// Create db folder
	dbPath := fmt.Sprintf("%s/%s", configPath, config.DB_DEFAULT_FOLDER)
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return fmt.Errorf("cannot create folder '%s': %w", dbPath, err)
	}
	slog.With(
		slog.String("path", dbPath),
	).Info("Created default relique database folder")

	// Create module config folder
	if err := os.MkdirAll(moduleInstallPath, 0755); err != nil {
		return fmt.Errorf("cannot create folder '%s': %w", moduleInstallPath, err)
	}
	slog.With(
		slog.String("path", moduleInstallPath),
	).Info("Created default relique module install folder")

	configFilePath := filepath.Clean(fmt.Sprintf("%s/relique.toml", configPath))
	config.New()
	config.Current.ModuleInstallPath = moduleInstallPath
	module.MODULES_INSTALL_PATH = moduleInstallPath
	if err := config.Write(configFilePath); err != nil {
		return fmt.Errorf("cannot create default configuration file: %w", err)
	}
	slog.With(
		slog.String("path", configFilePath),
	).Info("Created relique configuration file")

	// Create modules folder
	if err := os.MkdirAll(moduleInstallPath, 0755); err != nil {
		return fmt.Errorf("cannot create module install folder '%s': %w", moduleInstallPath, err)
	}
	slog.With(
		slog.String("path", moduleInstallPath),
	).Debug("Created modules install folder")

	// Install default modules
	if err := ModuleInstall(moduleInstallPath, "https://github.com/macarrie/relique-module-generic", false, false, false); err != nil {
		return fmt.Errorf("cannot install default generic module: %w", err)
	}
	slog.Debug("Installed default generic module")

	// Create clients folder
	clientsFolder := config.GetClientsCfgPath()
	if err := os.Mkdir(clientsFolder, 0755); err != nil {
		return fmt.Errorf("cannot create clients folder '%s': %w", clientsFolder, err)
	}
	slog.With(
		slog.String("path", clientsFolder),
	).Info("Created clients configuration folder")

	if err := ClientCreate("local", "localhost"); err != nil {
		return fmt.Errorf("cannot create default client: %w", err)
	}
	// TODO: Add example module to local client

	// Create catalog config folder
	catalogFolder := config.GetCatalogCfgPath()
	if err := os.Mkdir(catalogFolder, 0755); err != nil {
		return fmt.Errorf("cannot create catalog folder '%s': %w", catalogFolder, err)
	}
	slog.With(
		slog.String("path", catalogFolder),
	).Info("Created catalog folder")


	// Create repositories config folder
	reposFolder := config.GetReposCfgPath()
	if err := os.Mkdir(reposFolder, 0755); err != nil {
		return fmt.Errorf("cannot create repositories folder '%s': %w", reposFolder, err)
	}
	slog.With(
		slog.String("path", reposFolder),
	).Info("Created repositories configuration folder")

	// Create default repo
	if err := RepoCreateLocal("local", repoStoragePath, true); err != nil {
		return fmt.Errorf("cannot create default repository: %w", err)
	}

	return nil
}
