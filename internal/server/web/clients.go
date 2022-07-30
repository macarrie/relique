package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"
)

func getClients(c *gin.Context) {
	clients := serverConfig.Config.Clients
	c.JSON(http.StatusOK, clients)
}
