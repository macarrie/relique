package backup_job

import log "github.com/macarrie/relique/internal/logging"

type JobSearchParams struct {
	Module     string `json:"module,omitempty" mapstructure:"module"`
	Status     string `json:"status,omitempty" mapstructure:"status"`
	Client     string `json:"client,omitempty" mapstructure:"client"`
	Uuid       string `json:"uuid,omitempty" mapstructure:"uuid"`
	BackupType string `json:"backup_type,omitempty" mapstructure:"backup_type"`
	Limit      int    `json:"limit" mapstructure:"limit"`
}

func (p *JobSearchParams) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"uuid":        p.Uuid,
		"client":      p.Client,
		"module":      p.Module,
		"backup_type": p.BackupType,
		"status":      p.Status,
	})
}
