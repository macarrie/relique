package module

import (
	"fmt"
	"strings"

	"github.com/macarrie/relique/internal/types/displayable"
)

type ModuleDisplay struct {
	ModuleType        string   `json:"module_type"`
	Name              string   `json:"name"`
	BackupType        string   `json:"backup_type"`
	Schedules         []string `json:"schedules"`
	BackupPaths       []string `json:"backup_paths"`
	PreBackupScript   string   `json:"pre_backup_script"`
	PostBackupScript  string   `json:"post_backup_script"`
	PreRestoreScript  string   `json:"pre_restore_script"`
	PostRestoreScript string   `json:"post_restore_script"`
}

func (m Module) Display() displayable.Struct {
	var d displayable.Struct = ModuleDisplay{
		ModuleType:        m.ModuleType,
		Name:              m.Name,
		BackupType:        m.BackupType.String(),
		Schedules:         m.ScheduleNames,
		BackupPaths:       m.BackupPaths,
		PreBackupScript:   m.PreBackupScript,
		PostBackupScript:  m.PostBackupScript,
		PreRestoreScript:  m.PreRestoreScript,
		PostRestoreScript: m.PostRestoreScript,
	}

	return d
}

func (d ModuleDisplay) Summary() string {
	// TODO: Pretty display
	return fmt.Sprintf("Module summary: %s (type '%s')", d.Name, d.ModuleType)
}

func (d ModuleDisplay) Details() string {
	return fmt.Sprintf("Module DETAILS \n"+
		"----------- \n"+
		"\tName: %s\n"+
		"\tType: %s\n"+
		"\tSchedules: %s\n"+
		"\tBackup type: %s\n"+
		"\tBackup paths: %s\n"+
		"\tPre backup script: %s\n"+
		"\tPost backup script: %s\n"+
		"\tPre restore script: %s\n"+
		"\tPost restore script: %s\n",
		d.Name,
		d.ModuleType,
		strings.Join(d.Schedules, ", "),
		d.BackupType,
		strings.Join(d.BackupPaths, ", "),
		d.PreBackupScript,
		d.PostBackupScript,
		d.PreRestoreScript,
		d.PostRestoreScript)
}

func (d ModuleDisplay) TableHeaders() []string {
	return []string{
		"Name",
		"Type",
		"Backup type",
		"Schedules",
		"Backup paths",
	}
}

func (d ModuleDisplay) TableRow() []string {
	return []string{
		d.Name,
		d.ModuleType,
		d.BackupType,
		strings.Join(d.Schedules, ", "),
		strings.Join(d.BackupPaths, ", "),
	}
}
