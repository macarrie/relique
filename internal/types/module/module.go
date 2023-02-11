package module

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/hashicorp/go-multierror"

	"github.com/pkg/errors"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/pelletier/go-toml"
)

var MODULES_INSTALL_PATH = "/var/lib/relique/modules"
var ModulesInstallPathReadInConfig bool
var IsTest = false

type Module struct {
	ModuleType        string                 `json:"module_type" toml:"module_type"`
	Name              string                 `json:"name" toml:"name"`
	BackupType        backup_type.BackupType `json:"backup_type" toml:"backup_type"`
	Schedules         []schedule.Schedule    `json:"schedules" toml:"-"`
	ScheduleNames     []string               `json:"-" toml:"schedules"`
	AvailableVariants []string               `json:"available_variants" toml:"available_variants"`
	BackupPaths       []string               `json:"backup_paths" toml:"backup_paths"`
	PreBackupScript   string                 `json:"pre_backup_script" toml:"pre_backup_script"`
	PostBackupScript  string                 `json:"post_backup_script" toml:"post_backup_script"`
	PreRestoreScript  string                 `json:"pre_restore_script" toml:"pre_restore_script"`
	PostRestoreScript string                 `json:"post_restore_script" toml:"post_restore_script"`
	Variant           string                 `json:"variant" toml:"variant"`
	Params            map[string]interface{} `json:"params" toml:"params"`
}

// Set default value for dbPath according to OS if not already set in configuration file
func SetModulePathDefaultValue() {
	if IsTest {
		// Let unit test suite define module install path
		return
	}

	if ModulesInstallPathReadInConfig {
		return
	}

	switch runtime.GOOS {
	case "freebsd":
		MODULES_INSTALL_PATH = "/usr/local/relique/modules/"
	default:
		MODULES_INSTALL_PATH = "/var/lib/relique/modules/"
	}

	log.WithFields(log.Fields{
		"path": MODULES_INSTALL_PATH,
	}).Debug("Set default install path for modules")
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

func (m *Module) GetAbsScriptPath(module_name string, path string) string {
	return filepath.Clean(fmt.Sprintf("%s/%s/scripts/%s", MODULES_INSTALL_PATH, module_name, path))
}

func (m *Module) GetAvailableVariants() error {
	var availableVariants []string
	itemPath := fmt.Sprintf("%s/%s", MODULES_INSTALL_PATH, m.Name)
	files, err := os.ReadDir(itemPath)
	if err != nil {
		return errors.Wrap(err, "cannot list variants for module")
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

	if m.PreBackupScript == "" {
		m.PreBackupScript = defaults.PreBackupScript
	}

	if m.PostBackupScript == "" {
		m.PostBackupScript = defaults.PostBackupScript
	}

	if m.PreRestoreScript == "" {
		m.PreRestoreScript = defaults.PreRestoreScript
	}

	if m.PostRestoreScript == "" {
		m.PostRestoreScript = defaults.PostRestoreScript
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

func LoadFromFile(file string) (m Module, err error) {
	log.WithFields(log.Fields{
		"path": file,
	}).Debug("Loading module configuration parameters from file")

	f, err := os.Open(file)
	if err != nil {
		return Module{}, errors.Wrap(err, "cannot open file")
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = errors.Wrap(cerr, "cannot close file correctly")
		}
	}()

	content, _ := io.ReadAll(f)

	var module Module
	if err := toml.Unmarshal(content, &module); err != nil {
		return Module{}, errors.Wrap(err, "cannot parse toml file")
	}

	if err := module.Valid(); err != nil {
		return Module{}, errors.Wrap(err, "invalid module loaded from file")
	}

	return module, nil
}

func (m *Module) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"name":        m.Name,
		"type":        m.ModuleType,
		"backup_type": m.BackupType.String(),
		"variant":     m.GetVariant(),
	})
}

func (m *Module) Valid() error {
	var objErrors *multierror.Error
	if m.ModuleType == "" {
		objErrors = multierror.Append(objErrors, fmt.Errorf("empty module type"))
	}
	if m.Name == "" {
		objErrors = multierror.Append(objErrors, fmt.Errorf("empty module name"))
	}
	if m.BackupType.Type == backup_type.Unknown {
		objErrors = multierror.Append(objErrors, fmt.Errorf("unknown backup type"))
	}

	return objErrors.ErrorOrNil()
}

func (m *Module) ExtraParamsEnvVars(prefix string) []string {
	envTable := []string{}

	for param, val := range m.Params {
		envTable = append(envTable, fmt.Sprintf("%s%s=%s", prefix, param, fmt.Sprintf("%+v", val)))
	}

	return envTable
}
