package api

import (
	"fmt"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/image"
)

func ImageList(p api_helpers.PaginationParams) (api_helpers.PaginatedResponse[image.Image], error) {
	// TODO: Handle pagination
	imgCount, err := image.Count()
	if err != nil {
		return api_helpers.PaginatedResponse[image.Image]{}, fmt.Errorf("cannot count total images: %w", err)
	}

	imgs, err := image.Search(p, config.Current.ModuleInstallPath)
	if err != nil {
		return api_helpers.PaginatedResponse[image.Image]{}, fmt.Errorf("cannot get images from database: %w", err)
	}

	return api_helpers.PaginatedResponse[image.Image]{
		Count:      imgCount,
		Pagination: p,
		Data:       imgs,
	}, nil
}

func ImageGet(uuid string) (image.Image, error) {
	return image.Image{}, fmt.Errorf("TODO ImageGet")
}
