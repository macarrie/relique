package api

import (
	"fmt"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/module"
	"github.com/samber/lo"
)

func ModuleList(p api_helpers.PaginationParams) (api_helpers.PaginatedResponse[module.Module], error) {
	var returnModuleList []module.Module
	limit := p.Limit

	moduleList, err := module.GetLocallyInstalled(config.Current.ModuleInstallPath)
	if limit == 0 {
		returnModuleList = moduleList
	} else {
		returnModuleList = lo.Slice(moduleList, 0, int(p.Limit))
	}
	return api_helpers.PaginatedResponse[module.Module]{
		Count:      uint64(len(moduleList)),
		Pagination: p,
		Data:       returnModuleList,
	}, err
}

func ModuleGet(mod_name string) (module.Module, error) {
	modList, err := ModuleList(api_helpers.PaginationParams{})
	if err != nil {
		return module.Module{}, fmt.Errorf("cannot get installed module list: %w", err)
	}

	mod, ok := lo.Find(modList.Data, func(m module.Module) bool {
		return m.Name == mod_name
	})

	if ok {
		return mod, nil
	}
	return module.Module{}, fmt.Errorf("cannot find module in installed modules")
}

func ModuleInstall(modulesInstallPath string, path string, local bool, archive bool, force bool) error {
	return module.Install(modulesInstallPath, path, local, archive, force, true)
}
