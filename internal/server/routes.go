package server

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var f embed.FS

func staticHandler(engine *gin.Engine) {
	dist, _ := fs.Sub(f, "dist")
	fileServer := http.FileServer(http.FS(dist))

	engine.Use(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			// Check if the requested file exists
			_, err := fs.Stat(dist, strings.TrimPrefix(c.Request.URL.Path, "/"))
			if os.IsNotExist(err) {
				// If the file does not exist, serve index.html
				c.Request.URL.Path = "/"
			}

			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
}

func getRoutes() *gin.Engine {
	router := gin.Default()

	staticHandler(router)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping)
		v1.GET("/config", webAPIGetConfig)
		v1.GET("/config/version", webAPIGetVersion)

		v1.GET("/jobs", webAPIListJobs)
		v1.GET("/jobs/:uuid", webAPIGetJob)

		v1.GET("/clients", webAPIListClients)
		v1.GET("/clients/:name", webAPIGetClient)
		v1.GET("/clients/:name/ping", webAPIGetClientPing)

		v1.GET("/modules", webAPIListModules)
		v1.GET("/modules/:name", webAPIGetModule)

		v1.GET("/images", webAPIListImages)
		v1.GET("/images/:uuid", webAPIGetImage)
		v1.GET("/images/stats", webAPIGetImageStats)

		v1.GET("/repositories", webAPIListRepos)
		v1.GET("/repositories/:name", webAPIGetRepo)

	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
