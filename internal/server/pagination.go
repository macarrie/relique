package server

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/internal/api_helpers"
)

const DEFAULT_LIMIT = 10
const DEFAULT_OFFSET = 0

func getJobSearchParams(c *gin.Context) api_helpers.JobSearch {
	search := api_helpers.JobSearch{}
	query := c.Request.URL.Query()

	if client, ok := query["client"]; ok {
		search.ClientName = client[0]
	}
	if mod, ok := query["module"]; ok {
		search.ModuleName = mod[0]
	}
	if before, ok := query["before"]; ok {
		search.Before = before[0]
	}
	if after, ok := query["after"]; ok {
		search.After = after[0]
	}

	return search
}

func getImageSearchParams(c *gin.Context) api_helpers.ImageSearch {
	search := api_helpers.ImageSearch{}
	query := c.Request.URL.Query()

	if client, ok := query["client"]; ok {
		search.ClientName = client[0]
	}
	if mod, ok := query["module"]; ok {
		search.ModuleName = mod[0]
	}
	if before, ok := query["before"]; ok {
		search.Before = before[0]
	}
	if after, ok := query["after"]; ok {
		search.After = after[0]
	}

	return search
}

func getPagination(c *gin.Context) api_helpers.PaginationParams {
	limitFromQuery := c.DefaultQuery("limit", fmt.Sprintf("%d", DEFAULT_LIMIT))
	limit, err := strconv.Atoi(limitFromQuery)
	if err != nil {
		slog.With(
			slog.Int("value", DEFAULT_LIMIT),
			slog.Any("error", err),
		).Warn("Cannot get pagination limit parameter from request, using default value")
		limit = DEFAULT_LIMIT
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		slog.With(
			slog.Int("value", DEFAULT_LIMIT),
			slog.Any("error", err),
		).Warn("Cannot get pagination offset parameter from request, using default value")
		offset = DEFAULT_OFFSET
	}

	return api_helpers.PaginationParams{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	}
}
