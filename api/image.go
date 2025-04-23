package api

import (
	"fmt"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/image"
)

func ImageList(p api_helpers.PaginationParams, s api_helpers.ImageSearch) (api_helpers.PaginatedResponse[image.Image], error) {
	imgCount, err := image.Count()
	if err != nil {
		return api_helpers.PaginatedResponse[image.Image]{}, fmt.Errorf("cannot count total images: %w", err)
	}

	imgs, err := image.Search(p, s, config.Current.ModuleInstallPath)
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
	img, err := image.GetByUuid(uuid)
	if err != nil {
		return image.Image{}, fmt.Errorf("cannot get image from db: %w", err)
	}

	return img, nil
}
