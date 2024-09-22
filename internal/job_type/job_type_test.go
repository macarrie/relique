package job_type

import (
	"reflect"
	"testing"
)

func TestNew(t1 *testing.T) {
	type fields struct {
		Type uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name:   "backup",
			fields: fields{Backup},
			want:   Backup,
		},
		{
			name:   "restore",
			fields: fields{Restore},
			want:   Restore,
		},
		{
			name:   "unknown",
			fields: fields{123},
			want:   123,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := New(tt.fields.Type)
			if got := t.Type; got != tt.want {
				t1.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobType_String(t1 *testing.T) {
	type fields struct {
		Type uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "backup",
			fields: fields{Backup},
			want:   "backup",
		},
		{
			name:   "restore",
			fields: fields{Restore},
			want:   "restore",
		},
		{
			name:   "unknown",
			fields: fields{123},
			want:   "unknown",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &JobType{
				Type: tt.fields.Type,
			}
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want JobType
	}{
		{
			name: "backup",
			args: args{"backup"},
			want: JobType{Type: Backup},
		},
		{
			name: "restore",
			args: args{"restore"},
			want: JobType{Type: Restore},
		},
		{
			name: "unknown",
			args: args{"pouet"},
			want: JobType{Type: Unknown},
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

func TestJobType_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    JobType
		wantErr bool
	}{
		{
			name:    "backup",
			args:    args{data: []byte("backup")},
			want:    JobType{Type: Backup},
			wantErr: false,
		},
		{
			name:    "restore",
			args:    args{data: []byte("restore")},
			want:    JobType{Type: Restore},
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{data: []byte("invalid")},
			want:    JobType{Type: Unknown},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fromText JobType
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
		bt      JobType
		want    []byte
		wantErr bool
	}{
		{
			name:    "backup",
			bt:      JobType{Type: Backup},
			want:    []byte("backup"),
			wantErr: false,
		},
		{
			name:    "restore",
			bt:      JobType{Type: Restore},
			want:    []byte("restore"),
			wantErr: false,
		},
		{
			name:    "invalid variant",
			bt:      JobType{Type: 123},
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
