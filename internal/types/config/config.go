package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetConfigurationSubpath(relative_path string) string {
	base := filepath.Dir(viper.ConfigFileUsed())

	return fmt.Sprintf("%s/%s", base, relative_path)
}
