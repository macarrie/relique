package backup_job

import (
	"fmt"
	"strconv"
	"time"

	"github.com/macarrie/relique/internal/types/displayable"
)

type BackupJobDisplay struct {
	Uuid              string `json:"uuid"`
	Module            string `json:"module"`
	Client            string `json:"client"`
	BackupType        string `json:"backup_type"`
	Status            string `json:"status"`
	Done              string `json:"done"`
	Start             string `json:"start"`
	StartTimestamp    int64  `json:"start_timestamp"`
	End               string `json:"end"`
	EndTimestamp      int64  `json:"end_timestamp"`
	Duration          string `json:"duration"`
	DurationTimestamp int64  `json:"duration_timestamp"`
}

func (j BackupJob) Display() displayable.Struct {
	var d displayable.Struct = BackupJobDisplay{
		Uuid:              j.Uuid,
		Module:            j.Module.Name,
		Client:            j.Client.Name,
		BackupType:        j.BackupType.String(),
		Status:            j.Status.String(),
		Done:              strconv.FormatBool(j.Done),
		Start:             formatDatetime(j.StartTime),
		End:               formatDatetime(j.EndTime),
		StartTimestamp:    j.StartTime.Unix(),
		EndTimestamp:      j.EndTime.Unix(),
		Duration:          formatDuration(j.Duration()),
		DurationTimestamp: int64(j.Duration().Seconds()),
	}

	return d
}

func formatDuration(d time.Duration) string {
	return d.String()
}

func formatDatetime(t time.Time) string {
	if t.IsZero() {
		return "---"
	}

	return t.Format("2006/01/02 15:04:05")
}

func (d BackupJobDisplay) Summary() string {
	return fmt.Sprintf("Job summary: %v", d.Uuid)
}

func (d BackupJobDisplay) Details() string {
	return fmt.Sprintf("Job details: %v", d.Uuid)
}

func (d BackupJobDisplay) TableHeaders() []string {
	return []string{"UUID", "Done", "Status", "Client", "Module", "Backup type", "Duration", "Start", "End"}
}

func (d BackupJobDisplay) TableRow() []string {
	return []string{d.Uuid, d.Done, d.Status, d.Client, d.Module, d.BackupType, d.Duration, d.Start, d.End}
}
