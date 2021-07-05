package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/macarrie/relique/internal/types/job_type"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/relique_job"
)

func postJobLaunchScript(c *gin.Context) {
	scriptType, err := strconv.Atoi(c.Param("script_type"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid script type parameter received in request from relique server")
		return
	}

	if scriptType != relique_job.PreScript && scriptType != relique_job.PostScript {
		c.String(http.StatusBadRequest, "invalid script type parameter received in request from relique server")
		return
	}

	var job relique_job.ReliqueJob
	if err := c.ShouldBind(&job); err != nil {
		c.String(http.StatusBadRequest, "cannot parse received job parameters")
		return
	}

	if scriptType == relique_job.PreScript {
		if err := job.StartPreScript(); err != nil {
			job.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Error encountered during module pre script execution")
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := job.StartPostScript(); err != nil {
			job.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Error encountered during module post script execution")
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func postJobSetup(c *gin.Context) {
	var job relique_job.ReliqueJob
	if err := c.ShouldBind(&job); err != nil {
		c.String(http.StatusBadRequest, "cannot parse received job parameters")
		return
	}

	if err := job.PreFlightCheck(); err != nil {
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Error detected during job pre flight check")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if job.JobType.Type == job_type.Restore {
		// Create folders
		for _, path := range job.Module.BackupPaths {
			var pathToCreate string
			if job.RestoreDestination == "" {
				pathToCreate = filepath.Clean(path)
			} else {
				pathToCreate = filepath.Clean(fmt.Sprintf("%s/%s", job.RestoreDestination, path))
			}

			job.GetLog().WithFields(log.Fields{
				"path": pathToCreate,
			}).Debug("Creating path for data restoration")
			if err := os.MkdirAll(pathToCreate, 0755); err != nil {
				job.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": pathToCreate,
				}).Error("Cannot create path for restoration")
				c.String(http.StatusInternalServerError, err.Error())
			}
		}
	}

	c.Status(http.StatusOK)
}
