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
	Variant           string   `json:"variant"`
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
		Variant:           m.GetVariant(),
	}

	return d
}

func (d ModuleDisplay) Summary() string {
	return fmt.Sprintf("%s (type '%s/%s')", d.Name, d.ModuleType, d.Variant)
}

func (d ModuleDisplay) Details() string {
	return fmt.Sprintf(
		`	Name: 			%s
	Type: 			%s
	Variant: 		%s
	Schedules: 		%s
	Backup type: 		%s
	Backup paths: 		%s
	Pre backup script: 	%s
	Post backup script: 	%s
	Pre restore script:	%s
	Post restore script: 	%s`,
		d.Name,
		d.ModuleType,
		d.Variant,
		strings.Join(d.Schedules, ", "),
		d.BackupType,
		strings.Join(d.BackupPaths, ", "),
		orNone(d.PreBackupScript),
		orNone(d.PostBackupScript),
		orNone(d.PreRestoreScript),
		orNone(d.PostRestoreScript),
	)
}

func (d ModuleDisplay) TableHeaders() []string {
	return []string{
		"Name",
		"Type",
		"Variant",
		"Backup type",
		"Schedules",
		"Backup paths",
	}
}

func (d ModuleDisplay) TableRow() []string {
	return []string{
		d.Name,
		d.ModuleType,
		d.Variant,
		d.BackupType,
		strings.Join(d.Schedules, ", "),
		strings.Join(d.BackupPaths, ", "),
	}
}

func orNone(s string) string {
	if s == "" {
		return "---"
	}

	return s
}
