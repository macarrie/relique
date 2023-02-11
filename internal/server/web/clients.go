package web

import (
	"fmt"
	"net/http"

	log "github.com/macarrie/relique/internal/logging"
	clientObject "github.com/macarrie/relique/internal/types/client"
	serverApi "github.com/macarrie/relique/pkg/api/server"

	"github.com/gin-gonic/gin"
	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"
)

func getClients(c *gin.Context) {
	clients := serverConfig.Config.Clients
	for index := range clients {
		_ = serverApi.PingSSHClient(&clients[index])
	}

	c.JSON(http.StatusOK, clients)
}

func getClient(c *gin.Context) {
	clName := c.Param("name")
	var cl clientObject.Client

	// Get client from config to get modules details
	found := false
	fmt.Printf("CL name: %+v\n", clName)
	for _, configClient := range serverConfig.Config.Clients {
		fmt.Printf("Current client: %+v\n", configClient)
		if configClient.Name == clName {
			found = true
			cl = configClient
		}
	}
	if !found {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	fmt.Printf("CL: %+v\n", cl)

	if err := serverApi.PingSSHClient(&cl); err != nil {
		cl.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot ping client")
	}

	c.JSON(http.StatusOK, cl)
}
