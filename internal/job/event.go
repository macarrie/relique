package job

import (
	"log/slog"
	"strings"

	"github.com/macarrie/relique/internal/job_event"
	"github.com/macarrie/relique/internal/job_event_type"
)

// RecordEvent stores a job event. Failures are logged and never interrupt the job
func (j *Job) RecordEvent(t uint8, description string) {
	// External commands output (ssh, rsync) often ends with trailing newlines
	description = strings.TrimSpace(description)
	e := job_event.New(j.Uuid, t, description)
	if _, err := e.Save(); err != nil {
		j.GetLog().With(
			slog.Any("error", err),
			slog.String("event_type", e.Type.String()),
			slog.String("description", description),
		).Warn("Cannot save job event")
	}
}

// RecordError stores a job_error event with the error message as description
func (j *Job) RecordError(err error) {
	if err == nil {
		return
	}
	j.RecordEvent(job_event_type.JobError, err.Error())
}
