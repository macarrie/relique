package api

import (
	"github.com/macarrie/relique/internal/module"
)

func ModuleList() ([]module.Module, error) {
	return module.GetLocallyInstalled()
}

func ModuleInstall(path string, local bool, archive bool, force bool) error {
	return module.Install(path, local, archive, force, true)
}
