package cli

import (
	"strings"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/job_event"
	"github.com/macarrie/relique/internal/job_status"
	"github.com/macarrie/relique/internal/job_type"
	"github.com/macarrie/relique/internal/utils"
)

// jobWithEvents is the raw output of "relique job show --events": job fields plus its events
type jobWithEvents struct {
	job.Job
	Events []job_event.JobEvent `json:"events" toml:"events"`
}

// jobEventRow is a job event formatted for a markdown list item
type jobEventRow struct {
	Time        string
	Type        string
	Description string
}

// jobDetailsView holds the job data displayed by "relique job show", with values precomputed for the template
type jobDetailsView struct {
	job.Job
	Type        string
	StatusLabel string
	Duration    string
	SSHUser     string
	SSHPort     int
	IsRestore   bool
	HasStats    bool
	ShowEvents  bool
	Events      []jobEventRow
}

func newJobDetailsView(j job.Job, events []job_event.JobEvent, showEvents bool) jobDetailsView {
	v := jobDetailsView{
		Job:         j,
		Type:        jobTypeString(j),
		StatusLabel: statusLabel(j),
		Duration:    utils.FormatDuration(j.Duration()),
		SSHUser:     j.Client.SSHUser,
		SSHPort:     j.Client.SSHPort,
		IsRestore:   j.JobType.Type == job_type.Restore,
		HasStats:    j.Stats.NumberOfFiles > 0 || j.Stats.TotalBytesReceived > 0 || j.Stats.TotalBytesSent > 0,
	}
	if v.SSHUser == "" {
		v.SSHUser = client.DEFAULT_SSH_USER
	}
	if v.SSHPort == 0 {
		v.SSHPort = client.DEFAULT_SSH_PORT
	}
	if j.StartTime.IsZero() {
		v.Duration = "---"
	}

	v.ShowEvents = showEvents
	for _, e := range events {
		v.Events = append(v.Events, jobEventRow{
			// Time only: events of a job happen within its start and end dates displayed in the overview
			Time:        e.Time.Format("15:04:05"),
			Type:        e.Type.String(),
			Description: markdownInline(e.Description),
		})
	}

	return v
}

// markdownInline makes free text safe to display on a single markdown line
func markdownInline(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if s == "" {
		return "---"
	}
	return markdownEscaper.Replace(s)
}

var markdownEscaper = strings.NewReplacer(
	`\`, `\\`,
	"`", "\\`",
	"*", `\*`,
	"_", `\_`,
	"[", `\[`,
	"]", `\]`,
	"<", `\<`,
	">", `\>`,
	"#", `\#`,
)

// statusLabel prefixes the job status with a symbol readable without colors
func statusLabel(j job.Job) string {
	symbol := "?"
	switch j.Status.Status {
	case job_status.Pending:
		symbol = "…"
	case job_status.Active:
		symbol = "▶"
	case job_status.Success:
		symbol = "✔"
	case job_status.Incomplete:
		symbol = "⚠"
	case job_status.Error:
		symbol = "✘"
	}

	return symbol + " " + j.Status.String()
}

const jobDetailsTemplate = `# Job {{ .Uuid }}

**{{ .StatusLabel }}** · {{ .Type }} job · {{ .Client.Name }} / {{ .Module.Name }} · {{ .Duration }}

## Overview

| Property | Value |
| --- | --- |
| Status | {{ .StatusLabel }} |
| Type | {{ .Type }} |
| Started | {{ datetime .StartTime }} |
| Ended | {{ datetime .EndTime }} |
| Duration | {{ .Duration }} |
{{- if .PreviousJobUuid }}
| Diff reference | ` + "`{{ .PreviousJobUuid }}`" + ` |
{{- end }}
{{- if .IsRestore }}
| Restored image | ` + "`{{ .RestoreImageUuid }}`" + ` |
{{- end }}

## Client

| Property | Value |
| --- | --- |
| Name | {{ .Client.Name }} |
| Address | {{ .Client.Address }} |
| SSH | ` + "`{{ .SSHUser }}@{{ .Client.Address }}:{{ .SSHPort }}`" + ` |

## Module

| Property | Value |
| --- | --- |
| Name | {{ .Module.Name }} |
| Module type | {{ if .Module.ModuleType }}{{ .Module.ModuleType }}{{ else }}---{{ end }} |
| Backup paths | {{ if .Module.BackupPaths }}{{ join .Module.BackupPaths ", " }}{{ else }}---{{ end }} |
| Inclusions | {{ if .Module.Include }}{{ join .Module.Include ", " }}{{ else }}---{{ end }} |
| Exclusions | {{ if .Module.Exclude }}{{ join .Module.Exclude ", " }}{{ else }}---{{ end }} |
| Exclude CVS | {{ .Module.ExcludeCVS }} |

## Repository

| Property | Value |
| --- | --- |
| Name | {{ .Repository.GetName }} |
| Type | {{ .Repository.GetType }} |

## Statistics
{{ if .HasStats }}
| Statistic | Value |
| --- | --- |
| Elements | {{ .Stats.NumberOfFiles }} ({{ .Stats.NumberOfRegularFiles }} files, {{ .Stats.NumberOfDirectories }} directories) |
| Created | {{ .Stats.NumberOfCreatedFiles }} |
| Deleted | {{ .Stats.NumberOfDeletedFiles }} |
| Total size | {{ bytes .Stats.TotalFileSize }} |
| Transferred files size | {{ bytes .Stats.TotalTransferredFileSize }} |
| Literal data | {{ bytes .Stats.LiteralData }} |
| Matched data | {{ bytes .Stats.MatchedData }} |
| Sent / received | {{ bytes .Stats.TotalBytesSent }} / {{ bytes .Stats.TotalBytesReceived }} |
{{ else }}
_No statistics recorded for this job._
{{ end }}
{{- if .ShowEvents }}
## Events
{{ if .Events }}
{{- range .Events }}
- {{ .Time }} · ` + "`{{ .Type }}`" + ` · {{ .Description }}
{{- end }}
{{ else }}
_No events recorded for this job._
{{ end }}
{{- else }}
-----

Show events recorded during this job with ` + "`--events`" + `
{{ end }}`
