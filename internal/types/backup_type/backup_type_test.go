package backup_type

import (
	"reflect"
	"testing"

	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/logging"
)

func SetupTest(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)
	if err := db.InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}
}

func TestBackupType_String(t1 *testing.T) {
	type fields struct {
		Type uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "diff",
			fields: fields{Diff},
			want:   "diff",
		},
		{
			name:   "full",
			fields: fields{Full},
			want:   "full",
		},
		{
			name:   "unknown",
			fields: fields{123},
			want:   "unknown",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &BackupType{
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
		want BackupType
	}{
		{
			name: "diff",
			args: args{"diff"},
			want: BackupType{Type: Diff},
		},
		{
			name: "full",
			args: args{"full"},
			want: BackupType{Type: Full},
		},
		{
			name: "unknown",
			args: args{"pouet"},
			want: BackupType{Type: Unknown},
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

func TestBackupType_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    BackupType
		wantErr bool
	}{
		{
			name:    "diff",
			args:    args{data: []byte("diff")},
			want:    BackupType{Type: Diff},
			wantErr: false,
		},
		{
			name:    "full",
			args:    args{data: []byte("full")},
			want:    BackupType{Type: Full},
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{data: []byte("invalid")},
			want:    BackupType{Type: Unknown},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fromText BackupType
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

func TestBackupType_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		bt      BackupType
		want    []byte
		wantErr bool
	}{
		{
			name:    "diff",
			bt:      BackupType{Type: Diff},
			want:    []byte("diff"),
			wantErr: false,
		},
		{
			name:    "full",
			bt:      BackupType{Type: Full},
			want:    []byte("full"),
			wantErr: false,
		},
		{
			name:    "invalid variant",
			bt:      BackupType{Type: 123},
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
