package api

import (
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/module"
)

func ModuleList() ([]module.Module, error) {
	return module.GetLocallyInstalled(config.Current.ModuleInstallPath)
}

func ModuleInstall(modulesInstallPath string, path string, local bool, archive bool, force bool) error {
	return module.Install(modulesInstallPath, path, local, archive, force, true)
}
