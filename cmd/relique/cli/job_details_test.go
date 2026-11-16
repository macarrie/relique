package cli

import (
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/job_event"
	"github.com/macarrie/relique/internal/job_event_type"
	"github.com/macarrie/relique/internal/repo"
	"github.com/macarrie/relique/internal/utils"
)

func TestMarkdownInline(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "Sync started for path /etc", want: "Sync started for path /etc"},
		{name: "empty", in: "  ", want: "---"},
		{name: "newlines", in: "line one\r\nline two\n", want: "line one line two"},
		{name: "markdown_chars", in: "/var/my_dir/*.log [x] <y> #1 `z`", want: "/var/my\\_dir/\\*.log \\[x\\] \\<y\\> \\#1 \\`z\\`"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := markdownInline(tt.in); got != tt.want {
				t.Errorf("markdownInline() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewJobDetailsViewEvents(t *testing.T) {
	at := time.Date(2026, 11, 16, 10, 35, 44, 0, time.Local)
	events := []job_event.JobEvent{
		{JobUuid: "u", Time: at, Type: job_event_type.New(job_event_type.SyncError), Description: "rsync failed\n"},
	}

	v := newJobDetailsView(job.Job{Uuid: "u"}, events, true)
	if !v.ShowEvents || len(v.Events) != 1 {
		t.Fatalf("view events = %v (show %v), want 1 event shown", v.Events, v.ShowEvents)
	}
	want := jobEventRow{Time: "10:35:44", Type: "sync_error", Description: "rsync failed"}
	if v.Events[0] != want {
		t.Errorf("event row = %+v, want %+v", v.Events[0], want)
	}

	if v := newJobDetailsView(job.Job{Uuid: "u"}, nil, false); v.ShowEvents || len(v.Events) != 0 {
		t.Errorf("view without --events should not show events")
	}
}

func TestJobDetailsTemplateEvents(t *testing.T) {
	render := func(showEvents bool, events []job_event.JobEvent) string {
		t.Helper()
		localRepo := repo.RepoLocalNew("local", "/tmp", true)
		v := newJobDetailsView(job.Job{Uuid: "u", Repository: &localRepo}, events, showEvents)
		tpl := template.Must(template.New("job_details").Funcs(utils.TemplateFuncs()).Parse(jobDetailsTemplate))
		var b strings.Builder
		if err := tpl.Execute(&b, v); err != nil {
			t.Fatalf("template error: %v", err)
		}
		return b.String()
	}

	ev := []job_event.JobEvent{{Time: time.Now(), Type: job_event_type.New(job_event_type.JobStart), Description: "started"}}
	if out := render(true, ev); !strings.Contains(out, "## Events") || !strings.Contains(out, "`job_start` · started") {
		t.Errorf("template with events missing event list:\n%s", out)
	}
	if out := render(true, nil); !strings.Contains(out, "No events recorded") {
		t.Errorf("template with no events missing empty message:\n%s", out)
	}
	if out := render(false, nil); strings.Contains(out, "## Events") || !strings.Contains(out, "--events") {
		t.Errorf("template without --events should only hint at the flag:\n%s", out)
	}
}
