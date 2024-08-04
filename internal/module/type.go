package module

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/macarrie/relique/internal/backup_type"
)

type Module struct {
	ModuleType        string                 `json:"module_type" toml:"module_type"`
	Name              string                 `json:"name" toml:"name"`
	BackupType        backup_type.BackupType `json:"backup_type" toml:"backup_type"`
	Variant           string                 `json:"variant" toml:"variant"`
	AvailableVariants []string               `json:"available_variants" toml:"available_variants"`
	BackupPaths       []string               `json:"backup_paths" toml:"backup_paths"`
	Params            map[string]interface{} `json:"params" toml:"params"`
}

func (m *Module) String() string {
	return fmt.Sprintf("%s/%s", m.Name, m.GetVariant())
}

func (m *Module) GetVariant() string {
	if m.Variant == "" {
		return "default"
	}

	return m.Variant
}

func (m *Module) GetLog() *slog.Logger {
	return slog.With(
		slog.String("name", m.Name),
		slog.String("type", m.ModuleType),
		slog.String("backup_type", m.BackupType.String()),
		slog.String("variant", m.GetVariant()),
	)
}

func (m *Module) Valid() error {
	var objErrors *multierror.Error
	if m.ModuleType == "" {
		objErrors = multierror.Append(objErrors, fmt.Errorf("empty module type"))
	}
	if m.Name == "" {
		objErrors = multierror.Append(objErrors, fmt.Errorf("empty module name"))
	}
	// TODO: Use struct
	if m.BackupType.Type == backup_type.Unknown {
		objErrors = multierror.Append(objErrors, fmt.Errorf("unknown backup type"))
	}

	return objErrors.ErrorOrNil()
}

func (m *Module) GetAvailableVariants() error {
	var availableVariants []string
	itemPath := fmt.Sprintf("%s/%s", MODULES_INSTALL_PATH, m.Name)
	files, err := os.ReadDir(itemPath)
	if err != nil {
		return fmt.Errorf("cannot list variants for module: %w", err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".toml") {
			availableVariants = append(availableVariants, strings.TrimSuffix(file.Name(), ".toml"))
		}
	}

	m.AvailableVariants = availableVariants

	return nil
}

func (m *Module) LoadDefaultConfiguration() error {
	// Load module configuration from file with specified variant
	defaults, err := LoadFromFile(filepath.Clean(fmt.Sprintf("%s/%s/%s.toml", MODULES_INSTALL_PATH, m.ModuleType, m.GetVariant())))
	if err != nil {
		return err
	}

	if len(m.BackupPaths) == 0 {
		m.BackupPaths = defaults.BackupPaths
	}

	if m.Params == nil {
		m.Params = make(map[string]interface{})
	}

	for key := range defaults.Params {
		_, ok := m.Params[key]
		if !ok {
			m.Params[key] = defaults.Params[key]
		}
	}

	return nil
}
