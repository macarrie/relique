package api

import (
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

func ModuleInstall(modulesInstallPath string, path string, local bool, archive bool, force bool) error {
	return module.Install(modulesInstallPath, path, local, archive, force, true)
}
