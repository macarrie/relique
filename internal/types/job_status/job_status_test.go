package job_status

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
		want JobStatus
	}{
		{
			name: "pending",
			args: args{
				val: "pending",
			},
			want: JobStatus{Status: Pending},
		},
		{
			name: "active",
			args: args{
				val: "active",
			},
			want: JobStatus{Status: Active},
		},
		{
			name: "success",
			args: args{
				val: "success",
			},
			want: JobStatus{Status: Success},
		},
		{
			name: "incomplete",
			args: args{
				val: "incomplete",
			},
			want: JobStatus{Status: Incomplete},
		},
		{
			name: "error",
			args: args{
				val: "error",
			},
			want: JobStatus{Status: Error},
		},
		{
			name: "unknown",
			args: args{
				val: "pouet",
			},
			want: JobStatus{Status: Unknown},
		},
		{
			name: "random_value",
			args: args{
				val: "pouet",
			},
			want: JobStatus{Status: Unknown},
		},
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

func TestJobStatus_String(t *testing.T) {
	type fields struct {
		Status uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "pending",
			fields: fields{Pending},
			want:   "pending",
		},
		{
			name:   "active",
			fields: fields{Active},
			want:   "active",
		},
		{
			name:   "success",
			fields: fields{Success},
			want:   "success",
		},
		{
			name:   "incomplete",
			fields: fields{Incomplete},
			want:   "incomplete",
		},
		{
			name:   "error",
			fields: fields{Error},
			want:   "error",
		},
		{
			name:   "pouet",
			fields: fields{123},
			want:   "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JobStatus{
				Status: tt.fields.Status,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		status uint8
	}
	tests := []struct {
		name string
		args args
		want JobStatus
	}{
		// TODO: Add test cases.
		{
			name: "new",
			args: args{Incomplete},
			want: JobStatus{Status: Incomplete},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
