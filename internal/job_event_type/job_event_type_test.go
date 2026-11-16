package job_event_type

import (
	"reflect"
	"testing"
)

func TestFromString(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want JobEventType
	}{
		{name: "job_start", args: args{val: "job_start"}, want: JobEventType{Type: JobStart}},
		{name: "job_retry", args: args{val: "job_retry"}, want: JobEventType{Type: JobRetry}},
		{name: "diff_reference", args: args{val: "diff_reference"}, want: JobEventType{Type: DiffReference}},
		{name: "fallback_to_full", args: args{val: "fallback_to_full"}, want: JobEventType{Type: FallbackToFull}},
		{name: "job_setup", args: args{val: "job_setup"}, want: JobEventType{Type: JobSetup}},
		{name: "client_ping", args: args{val: "client_ping"}, want: JobEventType{Type: ClientPing}},
		{name: "client_ping_error", args: args{val: "client_ping_error"}, want: JobEventType{Type: ClientPingError}},
		{name: "sync_start", args: args{val: "sync_start"}, want: JobEventType{Type: SyncStart}},
		{name: "sync_end", args: args{val: "sync_end"}, want: JobEventType{Type: SyncEnd}},
		{name: "sync_error", args: args{val: "sync_error"}, want: JobEventType{Type: SyncError}},
		{name: "image_created", args: args{val: "image_created"}, want: JobEventType{Type: ImageCreated}},
		{name: "image_error", args: args{val: "image_error"}, want: JobEventType{Type: ImageError}},
		{name: "job_end", args: args{val: "job_end"}, want: JobEventType{Type: JobEnd}},
		{name: "job_error", args: args{val: "job_error"}, want: JobEventType{Type: JobError}},
		{name: "unknown", args: args{val: "unknown"}, want: JobEventType{Type: Unknown}},
		{name: "random_value", args: args{val: "pouet"}, want: JobEventType{Type: Unknown}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromString(tt.args.val)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobEventType_String(t *testing.T) {
	tests := []struct {
		name string
		t    uint8
		want string
	}{
		{name: "job_start", t: JobStart, want: "job_start"},
		{name: "sync_error", t: SyncError, want: "sync_error"},
		{name: "job_error", t: JobError, want: "job_error"},
		{name: "unknown", t: Unknown, want: "unknown"},
		{name: "zero_value", t: 0, want: "unknown"},
		{name: "out_of_range", t: 200, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			et := JobEventType{Type: tt.t}
			if got := et.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringRoundTrip(t *testing.T) {
	for i := uint8(JobStart); i < Unknown; i++ {
		et := New(i)
		if got := FromString(et.String()); got.Type != i {
			t.Errorf("FromString(%q) = %v, want %v", et.String(), got.Type, i)
		}
	}
}
