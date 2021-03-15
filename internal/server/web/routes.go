package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getRoutes() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping)
		v1.POST("/backup/register_job", postBackupRegisterJob)
		v1.PUT("/backup/jobs/:uuid/status", putBackupJobStatus)
		v1.PUT("/backup/jobs/:uuid/done", putBackupJobDone)
		v1.POST("/backup/jobs/:id/checksum", postBackupJobChecksum)
		v1.POST("/backup/jobs/:id/signature", postBackupJobSignature)
		v1.POST("/backup/jobs/:id/apply_diff", postBackupJobApplyDiff)
		v1.POST("/backup/jobs/:id/file", postBackupJobFile)
		v1.GET("/backup/jobs/:id/file", getBackupJobFile)
		v1.GET("/backup/jobs/:id/file_list", getBackupJobFileList)
		v1.POST("/backup/jobs/", getBackupJob)
	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
