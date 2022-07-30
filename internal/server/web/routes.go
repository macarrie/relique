package web

import (
	"fmt"
	"net/http"

	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func webPath(path string) string {
	web_root := serverConfig.Config.UiPath
	return fmt.Sprintf("%s/%s", web_root, path)
}

func getRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(static.Serve("/static", static.LocalFile(webPath("static"), true)))

	router.NoRoute(func(c *gin.Context) {
		c.File(webPath("index.html"))
	})

	root := router.Group("/")
	{
		root.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		})

		root.GET("/favicon.ico", func(c *gin.Context) {
			c.File(webPath("favicon.ico"))
		})
	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping)

		v1.POST("/jobs/", getBackupJob)
		v1.POST("/jobs/start", postJobStart)

		v1.POST("/retention/clean", postRetentionClean)

		v1.GET("/clients", getClients)
	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
