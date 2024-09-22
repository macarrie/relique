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

func TestJobStatus_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    JobStatus
		wantErr bool
	}{
		{
			name:    "pending",
			args:    args{data: []byte("pending")},
			want:    JobStatus{Status: Pending},
			wantErr: false,
		},
		{
			name:    "active",
			args:    args{data: []byte("active")},
			want:    JobStatus{Status: Active},
			wantErr: false,
		},
		{
			name:    "success",
			args:    args{data: []byte("success")},
			want:    JobStatus{Status: Success},
			wantErr: false,
		},
		{
			name:    "incomplete",
			args:    args{data: []byte("incomplete")},
			want:    JobStatus{Status: Incomplete},
			wantErr: false,
		},
		{
			name:    "error",
			args:    args{data: []byte("error")},
			want:    JobStatus{Status: Error},
			wantErr: false,
		},
		{
			name:    "unknown",
			args:    args{data: []byte("unknown")},
			want:    JobStatus{Status: Unknown},
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{data: []byte("invalid")},
			want:    JobStatus{Status: Unknown},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fromText JobStatus
			err := fromText.UnmarshalText(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(fromText, tt.want) {
				t.Errorf("UnmarshalText() got = %v, want %v", fromText, tt.want)
			}
		})
	}
}

func TestJobType_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		bt      JobStatus
		want    []byte
		wantErr bool
	}{
		{
			name:    "pending",
			bt:      JobStatus{Status: Pending},
			want:    []byte("pending"),
			wantErr: false,
		},
		{
			name:    "active",
			bt:      JobStatus{Status: Active},
			want:    []byte("active"),
			wantErr: false,
		},
		{
			name:    "success",
			bt:      JobStatus{Status: Success},
			want:    []byte("success"),
			wantErr: false,
		},
		{
			name:    "incomplete",
			bt:      JobStatus{Status: Incomplete},
			want:    []byte("incomplete"),
			wantErr: false,
		},
		{
			name:    "error",
			bt:      JobStatus{Status: Error},
			want:    []byte("error"),
			wantErr: false,
		},
		{
			name:    "unknown",
			bt:      JobStatus{Status: Unknown},
			want:    []byte("unknown"),
			wantErr: false,
		},
		{
			name:    "invalid variant",
			bt:      JobStatus{Status: 123},
			want:    []byte("unknown"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.bt.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalText() got = %v, want %v", got, tt.want)
			}
		})
	}
}
