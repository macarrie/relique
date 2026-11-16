package job_event

import (
	"time"

	"github.com/macarrie/relique/internal/job_event_type"
)

func New(jobUuid string, t uint8, description string) JobEvent {
	return JobEvent{
		JobUuid:     jobUuid,
		Time:        time.Now(),
		Type:        job_event_type.New(t),
		Description: description,
	}
}
