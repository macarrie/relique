package relique_job

import log "github.com/macarrie/relique/internal/logging"

type JobSearchParams struct {
	Module             string `json:"module,omitempty" mapstructure:"module"`
	ModuleVariant      string `json:"variant,omitempty" mapstructure:"variant"`
	Status             string `json:"status,omitempty" mapstructure:"status"`
	Client             string `json:"client,omitempty" mapstructure:"client"`
	Uuid               string `json:"uuid,omitempty" mapstructure:"uuid"`
	BackupType         string `json:"backup_type,omitempty" mapstructure:"backup_type"`
	JobType            string `json:"job_type,omitempty" mapstructure:"job_type"`
	RestoreJobUuid     string `json:"restore_job_uuid,omitempty" mapstructure:"restore_job_uuid"`
	RestoreDestination string `json:"restore_destination,omitempty" mapstructure:"restore_destination"`
	Count              string `json:"count,omitempty" mapstructure:"count"`
}

func (p *JobSearchParams) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"uuid":           p.Uuid,
		"client":         p.Client,
		"module":         p.Module,
		"module_variant": p.ModuleVariant,
		"backup_type":    p.BackupType,
		"status":         p.Status,
	})
}
