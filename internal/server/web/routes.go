package web

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os/exec"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/pkg/errors"

	"github.com/macarrie/relique/internal/types/client"

	"github.com/gin-gonic/gin"
)

func getRoutes() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping)
		v1.POST("/ping_client", pingClient)
		v1.POST("/backup/register_job", postBackupRegisterJob)
		v1.POST("/backup/jobs/:uuid/sync", postBackupJobSync)
		v1.GET("/backup/jobs/:uuid/sync_progress", getBackupJobSyncProgress)
		v1.PUT("/backup/jobs/:uuid/status", putBackupJobStatus)
		v1.PUT("/backup/jobs/:uuid/done", putBackupJobDone)
		v1.POST("/backup/jobs/", getBackupJob)
	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}

func pingClient(c *gin.Context) {
	var params client.ServerPingParams
	if err := c.ShouldBind(&params); err != nil {
		log.Error("Cannot bind server ping parameters received")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ip, _ := c.RemoteIP()
	if params.UseIPv4 && ip.IsLoopback() {
		ip = net.ParseIP("127.0.0.1")
	}
	if params.UseIPv6 && ip.IsLoopback() {
		ip = net.ParseIP("::1")
	}

	if params.UseIPv4 && ip.To4() == nil {
		params.Message = fmt.Sprintf("Client ping is forced to use IPv4 according to user settings but address '%s' cannot be converted to a valid IPv4 address", ip.String())
		c.JSON(http.StatusBadRequest, params)
		return
	}

	if params.UseIPv6 && ip.To16() == nil {
		params.Message = fmt.Sprintf("Client ping is forced to use IPv6 according to user settings but address '%s' cannot be converted to a valid IPv6 address", ip.String())
		c.JSON(http.StatusBadRequest, params)
		return
	}

	if params.UseIPv4 {
		params.ClientAddr = ip.To4().String()
	} else if params.UseIPv6 {
		params.ClientAddr = ip.To16().String()
	} else {
		params.ClientAddr = ip.String()
	}

	log.WithFields(log.Fields{
		"client_ip": params.ClientAddr,
	}).Info("Checking SSH connexion with client")

	sshPingCmd := exec.Command("ssh", "-f", "-o BatchMode=yes", fmt.Sprintf("relique@%s", params.ClientAddr), "echo 'ping'")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	sshPingCmd.Stdout = &stdout
	sshPingCmd.Stderr = &stderr

	err := sshPingCmd.Run()
	if err != nil {
		params.Message = errors.Wrap(err, fmt.Sprintf("cannot ping client via ssh:%s", stderr.String())).Error()
		c.JSON(http.StatusInternalServerError, params)
		return
	}

	c.JSON(http.StatusOK, params)
}
