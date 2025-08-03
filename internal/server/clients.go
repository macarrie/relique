package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/client"
)

func webAPIListClients(c *gin.Context) {
	page := getPagination(c)
	// search := getJobSearchParams(c)

	clList := api.ClientList(page, api_helpers.ClientSearch{})
	c.JSON(http.StatusOK, clList)
}

func webAPIGetClient(c *gin.Context) {
	name := c.Param("name")
	cl, err := api.ClientGet(name)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("name", name),
		).Error("Cannot find client in config")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, cl)
}

func webAPIGetClientPing(c *gin.Context) {
	name := c.Param("name")
	cl, err := api.ClientGet(name)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("name", name),
		).Error("Cannot find client in config")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	pingResult := api.ClientSSHPing(cl)
	var message string
	if pingResult == nil {
		message = ""
	} else {
		message = pingResult.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"ping_error": message,
	})
}

// Create client
func webAPIPostClient(c *gin.Context) {
	var cl client.Client
	slog.Warn("CREATING CLIENT")
	if bindErr := c.ShouldBindJSON(&cl); bindErr != nil {
		c.AbortWithError(http.StatusBadRequest, bindErr)
		return
	}

	if err := api.ClientCreate(cl.Name, cl.Address, cl.SSHUser, cl.SSHPort); err != nil {
		cl.GetLog().With(
			slog.Any("error", err),
		).Error("Cannot create client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(cl.Modules) > 0 {
		slog.Warn("SAVING CLIENT MODULES")
		if err := api.ClientSave(cl); err != nil {
			cl.GetLog().With(
				slog.Any("error", err),
			).Error("Cannot save client with modules info added")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.AbortWithStatus(http.StatusCreated)
}
