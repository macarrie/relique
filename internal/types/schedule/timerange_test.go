package schedule

import (
	"reflect"
	"testing"
	"time"
)

func TestTimerange_Active(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name      string
		timerange Timerange
		args      args
		want      bool
	}{
		{
			name: "before",
			timerange: Timerange{
				Start: time.Time{}.Add(12 * time.Hour),
				End:   time.Time{}.Add(14 * time.Hour),
			},
			args: args{
				now: time.Time{}.Add(11 * time.Hour),
			},
			want: false,
		},
		{
			name: "after",
			timerange: Timerange{
				Start: time.Time{}.Add(12 * time.Hour),
				End:   time.Time{}.Add(14 * time.Hour),
			},
			args: args{
				now: time.Time{}.Add(15 * time.Hour),
			},
			want: false,
		},
		{
			name: "active",
			timerange: Timerange{
				Start: time.Time{}.Add(12 * time.Hour),
				End:   time.Time{}.Add(14 * time.Hour),
			},
			args: args{
				now: time.Time{}.Add(13 * time.Hour),
			},
			want: true,
		},
		{
			name: "start_boundary",
			timerange: Timerange{
				Start: time.Time{}.Add(12 * time.Hour),
				End:   time.Time{}.Add(14 * time.Hour),
			},
			args: args{
				now: time.Time{}.Add(12 * time.Hour),
			},
			want: true,
		},
		{
			name: "end_boundary",
			timerange: Timerange{
				Start: time.Time{}.Add(12 * time.Hour),
				End:   time.Time{}.Add(14 * time.Hour),
			},
			args: args{
				now: time.Time{}.Add(14 * time.Hour),
			},
			want: false,
		},
		{
			name:      "empty",
			timerange: Timerange{},
			args: args{
				now: time.Now(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.timerange.Active(tt.args.now); got != tt.want {
				t.Errorf("Active() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimerange_MarshalText(t *testing.T) {
	tests := []struct {
		name      string
		timerange Timerange
		want      []byte
		wantErr   bool
	}{
		{
			name: "start_boundary",
			timerange: Timerange{
				Start: time.Time{}.Add(12*time.Hour + 34*time.Minute),
				End:   time.Time{}.Add(14*time.Hour + 52*time.Minute),
			},
			want: []byte("12:34-14:52"),
		},
		{
			name:      "empty",
			timerange: Timerange{},
			want:      []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.timerange.MarshalText()
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

func TestTimerange_UnmarshalText(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name      string
		timerange Timerange
		args      args
		wantErr   bool
	}{
		{
			name: "standard",
			timerange: Timerange{
				Start: time.Time{}.Add(1*time.Hour + 23*time.Minute),
				End:   time.Time{}.Add(12*time.Hour + 34*time.Minute),
			},
			args:    args{b: []byte("01:23-12:34")},
			wantErr: false,
		},
		{
			name: "with_whitespace",
			timerange: Timerange{
				Start: time.Time{}.Add(1*time.Hour + 23*time.Minute),
				End:   time.Time{}.Add(12*time.Hour + 34*time.Minute),
			},
			args:    args{b: []byte("  01:23  -   12:34  ")},
			wantErr: false,
		},
		{
			name:      "element_count",
			timerange: Timerange{},
			args:      args{b: []byte("01:00-02:00-03:00")},
			wantErr:   true,
		},
		{
			name:      "end_before_start",
			timerange: Timerange{},
			args:      args{b: []byte("02:00-01:00")},
			wantErr:   true,
		},
		{
			name:      "invalid_start",
			timerange: Timerange{},
			args:      args{b: []byte("66:66-02:00")},
			wantErr:   true,
		},
		{
			name:      "invalid_end",
			timerange: Timerange{},
			args:      args{b: []byte("01:00-66:66")},
			wantErr:   true,
		},
		{
			name:      "empty",
			timerange: Timerange{},
			args:      args{b: []byte("")},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.timerange.UnmarshalText(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
