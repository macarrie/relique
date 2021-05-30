package schedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/pkg/errors"
)

type Timerange struct {
	Start time.Time
	End   time.Time
}

func (r *Timerange) UnmarshalText(b []byte) error {
	timerangeStr := strings.TrimSpace(string(b))
	if timerangeStr == "" {
		*r = Timerange{}
		return nil
	}

	times := strings.Split(timerangeStr, "-")
	if len(times) != 2 && len(times) != 0 {
		return fmt.Errorf("timerange can only have zero or two components, found %d in '%s'", len(times), timerangeStr)
	}

	start, err := time.Parse("15:04", strings.TrimSpace(times[0]))
	startHour, startMin, _ := start.Clock()
	if err != nil {
		return errors.Wrap(err, "cannot parse range start date")
	}

	end, err := time.Parse("15:04", strings.TrimSpace(times[1]))
	endHour, endMin, _ := end.Clock()
	if err != nil {
		return errors.Wrap(err, "cannot parse range start date")
	}

	if !start.Before(end) {
		return fmt.Errorf("start date must be before end date for range '%s'", string(b))
	}

	*r = Timerange{
		Start: time.Time{}.Add(time.Duration(startHour)*time.Hour + time.Duration(startMin)*time.Minute),
		End:   time.Time{}.Add(time.Duration(endHour)*time.Hour + time.Duration(endMin)*time.Minute),
	}

	return nil
}

func (r *Timerange) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Timerange) String() string {
	layout := "15:04"
	if r.Start.IsZero() && r.End.IsZero() {
		return ""
	}

	return fmt.Sprintf("%s-%s", r.Start.Format(layout), r.End.Format(layout))
}

func (r *Timerange) Active(now time.Time) bool {
	hour, min, _ := now.Clock()
	timeNow := time.Time{}.Add(time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute)

	if timeNow.Equal(r.Start) {
		return true
	}
	if timeNow.After(r.Start) && timeNow.Before(r.End) {
		return true
	}

	return false
}

func (r *Timerange) Valid() error {
	var objErrors *multierror.Error

	if r.End.Before(r.Start) {
		objErrors = multierror.Append(objErrors, fmt.Errorf("range start must be before range end"))
	}

	return objErrors.ErrorOrNil()
}
