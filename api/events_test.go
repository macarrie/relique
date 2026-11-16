package api

import (
	"testing"

	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/job_event_type"
	"github.com/macarrie/relique/internal/job_type"
)

func recordTestEvents(uuid string) {
	j := job.Job{Uuid: uuid}
	j.RecordEvent(job_event_type.JobStart, "start")
	j.RecordEvent(job_event_type.SyncStart, "Sync started for path /etc")
	j.RecordEvent(job_event_type.JobEnd, "end")
}

func TestJobEvents(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, true, false)
	recordTestEvents(a.uuid)

	events, err := JobEvents(a.uuid)
	if err != nil {
		t.Fatalf("JobEvents() error = %v", err)
	}

	wantTypes := []uint8{job_event_type.JobStart, job_event_type.SyncStart, job_event_type.JobEnd}
	if len(events) != len(wantTypes) {
		t.Fatalf("JobEvents() returned %d events, want %d", len(events), len(wantTypes))
	}
	for i, want := range wantTypes {
		if events[i].Type.Type != want {
			t.Errorf("event %d type = %v, want %v", i, events[i].Type.Type, want)
		}
	}
}

func TestJobEventsUnknownJob(t *testing.T) {
	setupRemoveTestEnv(t)

	if _, err := JobEvents("does-not-exist"); err == nil {
		t.Errorf("JobEvents() on unknown job should return an error")
	}
}

func TestJobRemoveDeletesEvents(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, true, false)
	other := createJob(t, r, job_type.Backup, true, false)
	recordTestEvents(a.uuid)
	recordTestEvents(other.uuid)

	if err := JobRemove(a.uuid, false); err != nil {
		t.Fatalf("JobRemove() error = %v", err)
	}

	if _, err := JobEvents(a.uuid); err == nil {
		t.Errorf("JobEvents() on removed job should return an error")
	}
	events, err := JobEvents(other.uuid)
	if err != nil || len(events) != 3 {
		t.Errorf("other job events = %d (error %v), want 3", len(events), err)
	}
}
