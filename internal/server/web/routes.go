package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"

	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"

	"github.com/gin-gonic/gin"
)

func webPath(path string) string {
	web_root := serverConfig.Config.UiPath
	return fmt.Sprintf("%s/%s", web_root, path)
}

func getRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(static.Serve("/ui/static", static.LocalFile(webPath("static"), true)))

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		fmt.Println(path)
		fmt.Println(method)
		if strings.HasPrefix(path, "/ui") {
			c.File(webPath("index.html"))
		}
	})

	root := router.Group("/")
	{
		root.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/ui")
		})
	}

	ui := router.Group("/ui")
	{
		ui.GET("/favicon.ico", func(c *gin.Context) {
			c.File(webPath("favicon.ico"))
		})

		ui.GET("/", func(c *gin.Context) {
			c.File(webPath("index.html"))
		})

	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping)

		v1.POST("/jobs/", searchJob)
		v1.POST("/jobs/start", postJobStart)
		v1.GET("/jobs/:uuid", getJob)
		v1.GET("/jobs/:uuid/logs", getJobLogs)

		v1.POST("/retention/clean", postRetentionClean)

		v1.GET("/clients", getClients)
		v1.GET("/clients/:name", getClient)

		v1.GET("/modules", getModules)
		v1.GET("/modules/:name", getModule)
	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
