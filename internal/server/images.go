package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/image"
	"github.com/samber/lo"
)

func webAPIListImages(c *gin.Context) {
	page := getPagination(c)
	search := getImageSearchParams(c)

	imgs, err := api.ImageList(page, search)
	if err != nil {
		slog.With(
			slog.Any("error", err),
		).Error("Cannot get image list")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, imgs)
}

func webAPIGetImage(c *gin.Context) {
	uuid := c.Param("uuid")
	img, err := api.ImageGet(uuid)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("uuid", uuid),
		).Error("Cannot find image in database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, img)
}

func webAPIGetImageStats(c *gin.Context) {
	page := getPagination(c)
	search := getImageSearchParams(c)

	imgs, err := api.ImageList(page, search)
	if err != nil {
		slog.With(
			slog.Any("error", err),
		).Error("Cannot get image list")
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	totalSize := lo.SumBy(imgs.Data, func(img image.Image) uint64 {
		return img.SizeOnDisk
	})
	c.JSON(http.StatusOK, gin.H{
		"count":      imgs.Count,
		"total_size": totalSize,
	})
}
