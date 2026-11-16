package job_event

import (
	"log/slog"
	"time"

	"github.com/macarrie/relique/internal/job_event_type"
)

// JobEvent is a timestamped record of a state change in a job
type JobEvent struct {
	// Database ID
	ID int64 `json:"-" toml:"-"`

	JobUuid     string                      `json:"job_uuid" toml:"job_uuid"`
	Time        time.Time                   `json:"time" toml:"time"`
	Type        job_event_type.JobEventType `json:"type" toml:"type"`
	Description string                      `json:"description" toml:"description"`
}

func (e *JobEvent) GetLog() *slog.Logger {
	return slog.With(
		slog.String("job_uuid", e.JobUuid),
		slog.String("event_type", e.Type.String()),
	)
}
