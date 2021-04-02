package common

import (
	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/relique/internal/logging"
	client "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/schedule"
	"github.com/spf13/viper"
)

var customConfigFilePath string
var customConfigFile bool

type Configuration struct {
	Version                   string `mapstructure:",omitempty"`
	Clients                   []client.Client
	Schedules                 []schedule.Schedule
	BindAddr                  string `mapstructure:"bind_addr"`
	PublicAddress             string `mapstructure:"public_address"`
	Port                      uint32
	SSLCert                   string `mapstructure:"ssl_cert"`
	SSLKey                    string `mapstructure:"ssl_key"`
	StrictSSLCertificateCheck bool   `mapstructure:"strict_ssl_certificate_check"`
	ClientCfgPath             string `mapstructure:"client_cfg_path"`
	SchedulesCfgPath          string `mapstructure:"schedules_cfg_path"`
	BackupStoragePath         string `mapstructure:"backup_storage_path"`
}

func Load(fileName string) (Configuration, error) {
	viper.SetConfigType("toml")

	if customConfigFile {
		viper.SetConfigFile(customConfigFilePath)
	} else {
		viper.SetConfigName(fileName)
		viper.AddConfigPath("$HOME/.config/relique/")
		viper.AddConfigPath("/etc/relique/")
	}

	viper.SetEnvPrefix("RELIQUE")

	if err := viper.ReadInConfig(); err != nil {
		return Configuration{}, err
	}

	setDefaultValues()

	var conf Configuration
	if err := viper.Unmarshal(&conf); err != nil {
		return Configuration{}, err
	}

	log.WithFields(log.Fields{
		"file": viper.ConfigFileUsed(),
	}).Info("Configuration file loaded")

	return conf, nil
}

func setDefaultValues() {
	viper.SetDefault("version", "")
	viper.SetDefault("bind_addr", "0.0.0.0")
	viper.SetDefault("public_address", "localhost")
	viper.SetDefault("port", 8433)
	viper.SetDefault("ssl_cert", "/etc/relique/certs/cert.pem")
	viper.SetDefault("ssl_key", "/etc/relique/certs/key.pem")
	viper.SetDefault("strict_ssl_certificate_check", true)
	viper.SetDefault("client_cfg_path", "clients")
	viper.SetDefault("schedules_cfg_path", "schedules")
	viper.SetDefault("backup_storage_path", "/opt/relique")
}

func UseFile(filePath string) {
	log.WithFields(log.Fields{
		"file": filePath,
	}).Info("Using specified configuration file")
	customConfigFile = true
	customConfigFilePath = filePath
}

// TODO: Configuration validity checks
func Check() error {
	var errorList *multierror.Error

	return errorList.ErrorOrNil()
}
