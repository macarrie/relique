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
		name    string
		args    args
		want    JobStatus
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "pending",
			args: args{
				val: "pending",
			},
			want:    JobStatus{Status: Pending},
			wantErr: false,
		},
		{
			name: "active",
			args: args{
				val: "active",
			},
			want:    JobStatus{Status: Active},
			wantErr: false,
		},
		{
			name: "success",
			args: args{
				val: "success",
			},
			want:    JobStatus{Status: Success},
			wantErr: false,
		},
		{
			name: "incomplete",
			args: args{
				val: "incomplete",
			},
			want:    JobStatus{Status: Incomplete},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				val: "error",
			},
			want:    JobStatus{Status: Error},
			wantErr: false,
		},
		{
			name: "unknown",
			args: args{
				val: "pouet",
			},
			want:    JobStatus{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
