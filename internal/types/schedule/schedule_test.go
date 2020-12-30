package schedule

import (
	"testing"
	"time"

	"github.com/macarrie/relique/internal/logging"
)

func TestLoadFromPath(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)

	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
		wantErr bool
	}{
		{
			name: "load",
			args: args{
				p: "../../../test/config/schedules",
			},
			wantLen: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFromPath(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("LoadFromPath() got = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestSchedule_Active(t *testing.T) {
	type fields struct {
		Name      string
		Monday    Timeranges
		Tuesday   Timeranges
		Wednesday Timeranges
		Thursday  Timeranges
		Friday    Timeranges
		Saturday  Timeranges
		Sunday    Timeranges
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "active_monday",
			fields: fields{
				Name: "active_monday",
				Monday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.Add(2 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive_monday",
			fields: fields{
				Name: "inactive_monday",
				Monday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.Add(4 * time.Hour),
			},
			want: false,
		},
		{
			name: "active_tuesday",
			fields: fields{
				Name: "active_tuesday",
				Tuesday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 1).Add(2 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive_tuesday",
			fields: fields{
				Name: "inactive_tuesday",
				Tuesday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 1).Add(4 * time.Hour),
			},
			want: false,
		},
		{
			name: "active_wednesday",
			fields: fields{
				Name: "active_wednesday",
				Wednesday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 2).Add(2 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive_wednesday",
			fields: fields{
				Name: "inactive_wednesday",
				Wednesday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 1).Add(4 * time.Hour),
			},
			want: false,
		},
		{
			name: "active_thursday",
			fields: fields{
				Name: "active_thursday",
				Thursday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 3).Add(2 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive_thursday",
			fields: fields{
				Name: "inactive_thursday",
				Thursday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 1).Add(4 * time.Hour),
			},
			want: false,
		},
		{
			name: "active_friday",
			fields: fields{
				Name: "active_friday",
				Friday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 4).Add(2 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive_friday",
			fields: fields{
				Name: "inactive_friday",
				Friday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 1).Add(4 * time.Hour),
			},
			want: false,
		},
		{
			name: "active_saturday",
			fields: fields{
				Name: "active_saturday",
				Saturday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 5).Add(2 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive_saturday",
			fields: fields{
				Name: "inactive_saturday",
				Saturday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 1).Add(4 * time.Hour),
			},
			want: false,
		},
		{
			name: "active_sunday",
			fields: fields{
				Name: "active_sunday",
				Sunday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 6).Add(2 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive_sunday",
			fields: fields{
				Name: "inactive_sunday",
				Sunday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1 * time.Hour),
							End:   time.Time{}.Add(3 * time.Hour),
						},
					},
				},
			},
			args: args{
				now: time.Time{}.AddDate(0, 0, 1).Add(4 * time.Hour),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schedule{
				Name:      tt.fields.Name,
				Monday:    tt.fields.Monday,
				Tuesday:   tt.fields.Tuesday,
				Wednesday: tt.fields.Wednesday,
				Thursday:  tt.fields.Thursday,
				Friday:    tt.fields.Friday,
				Saturday:  tt.fields.Saturday,
				Sunday:    tt.fields.Sunday,
			}
			if got := s.Active(tt.args.now); got != tt.want {
				t.Errorf("Active() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFromFile(t *testing.T) {
	logging.Setup(true, logging.TEST_LOG_PATH)

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    Schedule
		wantErr bool
	}{
		{
			name: "example",
			args: args{
				file: "../../../test/config/schedules/example.toml",
			},
			want: Schedule{
				Name: "example",
				Monday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1*time.Hour + 10*time.Minute),
							End:   time.Time{}.Add(2*time.Hour + 10*time.Minute),
						},
						{
							Start: time.Time{}.Add(3*time.Hour + 10*time.Minute),
							End:   time.Time{}.Add(4*time.Hour + 10*time.Minute),
						},
					},
				},
				Tuesday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1*time.Hour + 15*time.Minute),
							End:   time.Time{}.Add(2*time.Hour + 15*time.Minute),
						},
						{
							Start: time.Time{}.Add(3*time.Hour + 15*time.Minute),
							End:   time.Time{}.Add(4*time.Hour + 15*time.Minute),
						},
					},
				},
				Wednesday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1*time.Hour + 20*time.Minute),
							End:   time.Time{}.Add(2*time.Hour + 20*time.Minute),
						},
						{
							Start: time.Time{}.Add(3*time.Hour + 20*time.Minute),
							End:   time.Time{}.Add(4*time.Hour + 20*time.Minute),
						},
					},
				},
				Thursday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1*time.Hour + 25*time.Minute),
							End:   time.Time{}.Add(2*time.Hour + 25*time.Minute),
						},
						{
							Start: time.Time{}.Add(3*time.Hour + 25*time.Minute),
							End:   time.Time{}.Add(4*time.Hour + 25*time.Minute),
						},
					},
				},
				Friday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1*time.Hour + 30*time.Minute),
							End:   time.Time{}.Add(2*time.Hour + 30*time.Minute),
						},
						{
							Start: time.Time{}.Add(3*time.Hour + 30*time.Minute),
							End:   time.Time{}.Add(4*time.Hour + 30*time.Minute),
						},
					},
				},
				Saturday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1*time.Hour + 35*time.Minute),
							End:   time.Time{}.Add(2*time.Hour + 35*time.Minute),
						},
						{
							Start: time.Time{}.Add(3*time.Hour + 35*time.Minute),
							End:   time.Time{}.Add(4*time.Hour + 35*time.Minute),
						},
					},
				},
				Sunday: Timeranges{
					Ranges: []Timerange{
						{
							Start: time.Time{}.Add(1*time.Hour + 40*time.Minute),
							End:   time.Time{}.Add(2*time.Hour + 40*time.Minute),
						},
						{
							Start: time.Time{}.Add(3*time.Hour + 40*time.Minute),
							End:   time.Time{}.Add(4*time.Hour + 40*time.Minute),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "default_values",
			args: args{
				file: "../../../test/config/schedules/empty.toml",
			},
			want:    Schedule{},
			wantErr: true,
		},
		{
			name: "parse_error",
			args: args{
				file: "../../../test/config/schedules/parse_error.toml",
			},
			want:    Schedule{},
			wantErr: true,
		},
		{
			name: "unreadable",
			args: args{
				file: "../../../test/config/schedules/unreadable.toml.test",
			},
			want:    Schedule{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadFromFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// DeepEqual does not seem to work on time.Time, use field by field comparison instead
			if got.Name != tt.want.Name ||
				got.Monday.String() != tt.want.Monday.String() ||
				got.Tuesday.String() != tt.want.Tuesday.String() ||
				got.Wednesday.String() != tt.want.Wednesday.String() ||
				got.Thursday.String() != tt.want.Thursday.String() ||
				got.Friday.String() != tt.want.Friday.String() ||
				got.Saturday.String() != tt.want.Saturday.String() ||
				got.Sunday.String() != tt.want.Sunday.String() {
				t.Errorf("loadFromFile() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSchedule_Valid(t *testing.T) {
	type fields struct {
		Name      string
		Monday    Timeranges
		Tuesday   Timeranges
		Wednesday Timeranges
		Thursday  Timeranges
		Friday    Timeranges
		Saturday  Timeranges
		Sunday    Timeranges
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Name: "valid",
			},
			wantErr: false,
		},
		{
			name:    "missing_name",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "invalid_range_monday",
			fields: fields{
				Name: "invalid_range_monday",
				Monday: Timeranges{Ranges: []Timerange{
					{
						Start: time.Time{}.Add(2 * time.Hour),
						End:   time.Time{}.Add(1 * time.Hour),
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid_range_tuesday",
			fields: fields{
				Name: "invalid_range_tuesday",
				Tuesday: Timeranges{Ranges: []Timerange{
					{
						Start: time.Time{}.Add(2 * time.Hour),
						End:   time.Time{}.Add(1 * time.Hour),
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid_range_wednesday",
			fields: fields{
				Name: "invalid_range_wednesday",
				Wednesday: Timeranges{Ranges: []Timerange{
					{
						Start: time.Time{}.Add(2 * time.Hour),
						End:   time.Time{}.Add(1 * time.Hour),
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid_range_thursday",
			fields: fields{
				Name: "invalid_range_thursday",
				Thursday: Timeranges{Ranges: []Timerange{
					{
						Start: time.Time{}.Add(2 * time.Hour),
						End:   time.Time{}.Add(1 * time.Hour),
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid_range_friday",
			fields: fields{
				Name: "invalid_range_friday",
				Friday: Timeranges{Ranges: []Timerange{
					{
						Start: time.Time{}.Add(2 * time.Hour),
						End:   time.Time{}.Add(1 * time.Hour),
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid_range_saturday",
			fields: fields{
				Name: "invalid_range_saturday",
				Saturday: Timeranges{Ranges: []Timerange{
					{
						Start: time.Time{}.Add(2 * time.Hour),
						End:   time.Time{}.Add(1 * time.Hour),
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid_range_sunday",
			fields: fields{
				Name: "invalid_range_sunday",
				Sunday: Timeranges{Ranges: []Timerange{
					{
						Start: time.Time{}.Add(2 * time.Hour),
						End:   time.Time{}.Add(1 * time.Hour),
					},
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schedule{
				Name:      tt.fields.Name,
				Monday:    tt.fields.Monday,
				Tuesday:   tt.fields.Tuesday,
				Wednesday: tt.fields.Wednesday,
				Thursday:  tt.fields.Thursday,
				Friday:    tt.fields.Friday,
				Saturday:  tt.fields.Saturday,
				Sunday:    tt.fields.Sunday,
			}
			if err := s.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
