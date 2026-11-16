package job_event

import (
	"testing"
	"time"

	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/job_event_type"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	if err := db.Init(t.TempDir()); err != nil {
		t.Fatalf("cannot init test db: %v", err)
	}
}

func TestSaveAndList(t *testing.T) {
	setupTestDB(t)

	start := time.Now()
	events := []JobEvent{
		{JobUuid: "job-a", Time: start, Type: job_event_type.New(job_event_type.JobStart), Description: "backup job (full)"},
		{JobUuid: "job-b", Time: start, Type: job_event_type.New(job_event_type.JobStart), Description: "other job"},
		{JobUuid: "job-a", Time: start.Add(2 * time.Second), Type: job_event_type.New(job_event_type.JobEnd), Description: "success"},
		// Same timestamp as the first event: insertion order must be kept
		{JobUuid: "job-a", Time: start, Type: job_event_type.New(job_event_type.ClientPing), Description: ""},
	}
	for i := range events {
		id, err := events[i].Save()
		if err != nil || id == 0 {
			t.Fatalf("Save() id = %v, error = %v", id, err)
		}
	}

	got, err := ListByJobUuid("job-a")
	if err != nil {
		t.Fatalf("ListByJobUuid() error = %v", err)
	}

	wantTypes := []uint8{job_event_type.JobStart, job_event_type.ClientPing, job_event_type.JobEnd}
	if len(got) != len(wantTypes) {
		t.Fatalf("ListByJobUuid() returned %d events, want %d", len(got), len(wantTypes))
	}
	for i, want := range wantTypes {
		if got[i].Type.Type != want {
			t.Errorf("event %d type = %v, want %v", i, got[i].Type.Type, want)
		}
		if got[i].JobUuid != "job-a" {
			t.Errorf("event %d job uuid = %v, want job-a", i, got[i].JobUuid)
		}
	}
	if got[0].Description != "backup job (full)" {
		t.Errorf("event 0 description = %q", got[0].Description)
	}
	if !got[2].Time.Equal(events[2].Time) {
		t.Errorf("event 2 time = %v, want %v", got[2].Time, events[2].Time)
	}
}

func TestListUnknownJob(t *testing.T) {
	setupTestDB(t)

	got, err := ListByJobUuid("does-not-exist")
	if err != nil {
		t.Fatalf("ListByJobUuid() error = %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("ListByJobUuid() = %v, want empty non-nil slice", got)
	}
}

func TestDeleteByJobUuid(t *testing.T) {
	setupTestDB(t)

	for _, uuid := range []string{"job-a", "job-a", "job-b"} {
		e := New(uuid, job_event_type.JobStart, "")
		if _, err := e.Save(); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	tx, err := db.Handler().Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := DeleteByJobUuid(tx, "job-a"); err != nil {
		tx.Rollback()
		t.Fatalf("DeleteByJobUuid() error = %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if got, _ := ListByJobUuid("job-a"); len(got) != 0 {
		t.Errorf("job-a still has %d events after delete", len(got))
	}
	if got, _ := ListByJobUuid("job-b"); len(got) != 1 {
		t.Errorf("job-b has %d events after deleting job-a events, want 1", len(got))
	}
}
