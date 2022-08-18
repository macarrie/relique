package web

import (
	"net/http"
	"strconv"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/client"
	serverApi "github.com/macarrie/relique/pkg/api/server"

	"github.com/gin-gonic/gin"
	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"
)

func getClients(c *gin.Context) {
	clients := serverConfig.Config.Clients
	for index, cl := range clients {
		id, _ := client.GetID(cl.Name, nil)
		if id != 0 {
			clients[index].ID = id
		}

		_ = serverApi.PingSSHClient(&clients[index])
	}

	c.JSON(http.StatusOK, clients)
}

func getClient(c *gin.Context) {
	idStr := c.Param("id")
	id, strconvErr := strconv.ParseInt(idStr, 10, 64)
	if strconvErr != nil {
		log.Error("Cannot parse client id from request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cl, getClientErr := client.GetByID(id)
	if getClientErr != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).Error("Cannot find client in database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Get client from config to get modules details
	for _, configClient := range serverConfig.Config.Clients {
		if configClient.ID == id {
			cl = configClient
		}
	}

	if err := serverApi.PingSSHClient(&cl); err != nil {
		cl.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot ping client")
	}

	c.JSON(http.StatusOK, cl)
}
