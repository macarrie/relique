package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/pagination"
	"strconv"
)

const DEFAULT_LIMIT = 10
const DEFAULT_OFFSET = 0

func getPagination(c *gin.Context) pagination.Pagination {
	limitFromQuery := c.DefaultQuery("limit", fmt.Sprintf("%d", DEFAULT_LIMIT))
	limit, err := strconv.Atoi(limitFromQuery)
	if err != nil {
		log.WithFields(log.Fields{
			"value": DEFAULT_LIMIT,
			"err":   err,
		}).Warning("Cannot get pagination limit parameter from request, using default value")
		limit = DEFAULT_LIMIT
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		log.WithFields(log.Fields{
			"value": DEFAULT_LIMIT,
		}).Warning("Cannot get pagination offset parameter from request, using default value")
		offset = DEFAULT_OFFSET
	}

	return pagination.Pagination{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	}
}
