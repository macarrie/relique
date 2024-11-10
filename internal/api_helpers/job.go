package api_helpers

type JobSearch struct {
	ModuleName string `json:"module"`
	ClientName string `json:"client"`
	BackupType uint8  `json:"backup_type"`
	JobType    uint8  `json:"job_type"`
	Status     uint8  `json:"status"`
}
