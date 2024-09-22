package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
)

var customConfigFilePath string
var customConfigFile bool
var Current Configuration
var Loaded bool

var CLIENTS_DEFAULT_FOLDER string = "clients"
var REPOS_DEFAULT_FOLDER string = "repositories"
var MODULES_DEFAULT_FOLDER string = "/var/lib/relique/modules"
var DB_DEFAULT_FOLDER string = "db"

type Configuration struct {
	Clients      []client.Client   `json:"clients" toml:"clients"`
	Repositories []repo.Repository `json:"repositories" toml:"repositories"`
	Modules      []module.Module   `json:"modules" toml:"modules"`

	ClientCfgPath     string `mapstructure:"client_cfg_path" json:"client_cfg_path" toml:"client_cfg_path"`
	RepoCfgPath       string `mapstructure:"repo_cfg_path" json:"repo_cfg_path" toml:"repo_cfg_path"`
	ModuleInstallPath string `mapstructure:"module_install_path" json:"module_install_path" toml:"module_install_path"`
	DBPath            string `mapstructure:"db_path" json:"db_path" toml:"db_path"`
}

func New() {
	Current = Configuration{}
	setDefaultValues()
}

func Load(fileName string) error {
	viper.SetConfigType("toml")

	if customConfigFile {
		viper.SetConfigFile(customConfigFilePath)
	} else {
		viper.SetConfigName(fileName)
		viper.AddConfigPath("$HOME/.config/relique/")
		viper.AddConfigPath("/usr/local/etc/relique/")
		viper.AddConfigPath("/etc/relique/")
	}

	viper.SetEnvPrefix("RELIQUE")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var cfg Configuration
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	setDefaultValues()
	module.MODULES_INSTALL_PATH = cfg.ModuleInstallPath

	clients, err := client.LoadFromPath(getAbsCfgDir(cfg.ClientCfgPath, CLIENTS_DEFAULT_FOLDER))
	if err != nil {
		return fmt.Errorf("cannot load clients configuration: %w", err)
	}
	cfg.Clients = clients

	repos, err := repo.LoadFromPath(getAbsCfgDir(cfg.RepoCfgPath, REPOS_DEFAULT_FOLDER))
	if err != nil {
		return fmt.Errorf("cannot load repositories configuration: %w", err)
	}
	cfg.Repositories = repos

	Current = cfg
	Loaded = true
	slog.Debug("Loaded config file", slog.String("file", viper.ConfigFileUsed()))

	return nil
}

func Write(path string) error {
	viper.SetConfigFile(path)

	cfgFile, cfgErr := toml.Marshal(Current)
	if cfgErr != nil {
		return fmt.Errorf("cannot serialize config info to toml data: %w", cfgErr)
	}
	if err := os.WriteFile(path, cfgFile, 0644); err != nil {
		return fmt.Errorf("cannot export config info to file: %w", err)
	}

	return nil
}

func setDefaultValues() {
	if Current.ClientCfgPath == "" {
		Current.ClientCfgPath = CLIENTS_DEFAULT_FOLDER
	}
	if Current.RepoCfgPath == "" {
		Current.RepoCfgPath = REPOS_DEFAULT_FOLDER
	}
	if Current.ModuleInstallPath == "" {
		Current.ModuleInstallPath = MODULES_DEFAULT_FOLDER
		module.MODULES_INSTALL_PATH = MODULES_DEFAULT_FOLDER
	}
	if Current.DBPath == "" {
		Current.DBPath = DB_DEFAULT_FOLDER
	}
}

func getAbsoluteCfgPath(relative_path string) string {
	base := filepath.Dir(viper.ConfigFileUsed())

	return fmt.Sprintf("%s/%s", base, relative_path)
}

func getAbsCfgDir(cfgKey string, defaultValue string) string {
	var dir string
	if cfgKey == "" {
		dir = defaultValue
	} else {
		dir = cfgKey
	}

	return getAbsoluteCfgPath(dir)
}

func GetClientsCfgPath() string {
	return getAbsCfgDir(Current.ClientCfgPath, CLIENTS_DEFAULT_FOLDER)
}

func GetReposCfgPath() string {
	return getAbsCfgDir(Current.RepoCfgPath, REPOS_DEFAULT_FOLDER)
}

func GetDBPath() string {
	return getAbsCfgDir(Current.DBPath, DB_DEFAULT_FOLDER)
}

func UseFile(filePath string) {
	slog.Info("Using specified configuration file", slog.String("file_path", filePath))
	customConfigFile = true
	customConfigFilePath = filePath
}

// TODO: Configuration validity checks
func Check() error {
	var errorList *multierror.Error

	return errorList.ErrorOrNil()
}
