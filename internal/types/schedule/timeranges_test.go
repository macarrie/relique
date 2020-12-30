package schedule

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeranges_Active(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		ranges []Timerange
		args   args
		want   bool
	}{
		{
			name: "active",
			ranges: []Timerange{
				Timerange{
					Start: time.Time{}.Add(1*time.Hour + 11*time.Minute),
					End:   time.Time{}.Add(2*time.Hour + 22*time.Minute),
				},
				Timerange{
					Start: time.Time{}.Add(3*time.Hour + 33*time.Minute),
					End:   time.Time{}.Add(4*time.Hour + 44*time.Minute),
				},
			},
			args: args{now: time.Time{}.Add(4 * time.Hour)},
			want: true,
		},
		{
			name: "overlap",
			ranges: []Timerange{
				Timerange{
					Start: time.Time{}.Add(1*time.Hour + 11*time.Minute),
					End:   time.Time{}.Add(4*time.Hour + 44*time.Minute),
				},
				Timerange{
					Start: time.Time{}.Add(2*time.Hour + 22*time.Minute),
					End:   time.Time{}.Add(3*time.Hour + 33*time.Minute),
				},
			},
			args: args{now: time.Time{}.Add(4 * time.Hour)},
			want: true,
		},
		{
			name: "not_active",
			ranges: []Timerange{
				Timerange{
					Start: time.Time{}.Add(1*time.Hour + 11*time.Minute),
					End:   time.Time{}.Add(2*time.Hour + 22*time.Minute),
				},
				Timerange{
					Start: time.Time{}.Add(3*time.Hour + 33*time.Minute),
					End:   time.Time{}.Add(3*time.Hour + 50*time.Minute),
				},
			},
			args: args{now: time.Time{}.Add(4 * time.Hour)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Timeranges{
				Ranges: tt.ranges,
			}
			if got := r.Active(tt.args.now); got != tt.want {
				t.Errorf("Active() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeranges_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		ranges  []Timerange
		want    []byte
		wantErr bool
	}{
		{
			name: "normal",
			ranges: []Timerange{
				Timerange{
					Start: time.Time{}.Add(1*time.Hour + 11*time.Minute),
					End:   time.Time{}.Add(2*time.Hour + 22*time.Minute),
				},
				Timerange{
					Start: time.Time{}.Add(3*time.Hour + 33*time.Minute),
					End:   time.Time{}.Add(4*time.Hour + 44*time.Minute),
				},
			},
			want:    []byte("01:11-02:22,03:33-04:44"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Timeranges{
				Ranges: tt.ranges,
			}
			got, err := r.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalText() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestTimeranges_UnmarshalText(t *testing.T) {
	type fields struct {
		Ranges []Timerange
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{Ranges: []Timerange{
				Timerange{
					Start: time.Time{}.Add(1*time.Hour + 11*time.Minute),
					End:   time.Time{}.Add(2*time.Hour + 22*time.Minute),
				},
				Timerange{
					Start: time.Time{}.Add(3*time.Hour + 33*time.Minute),
					End:   time.Time{}.Add(4*time.Hour + 44*time.Minute),
				},
			}},
			args:    args{b: []byte("01:11-02:22,03:33-04:44")},
			wantErr: false,
		},
		{
			name: "with_whitespaces",
			fields: fields{Ranges: []Timerange{
				Timerange{
					Start: time.Time{}.Add(1*time.Hour + 11*time.Minute),
					End:   time.Time{}.Add(2*time.Hour + 22*time.Minute),
				},
				Timerange{
					Start: time.Time{}.Add(3*time.Hour + 33*time.Minute),
					End:   time.Time{}.Add(4*time.Hour + 44*time.Minute),
				},
			}},
			args:    args{b: []byte("  1:11  - 02:22 ,   3:33   - 4:44   ")},
			wantErr: false,
		},
		{
			name:    "parse_error",
			fields:  fields{Ranges: []Timerange{}},
			args:    args{b: []byte("pouet")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Timeranges{
				Ranges: tt.fields.Ranges,
			}
			if err := r.UnmarshalText(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
