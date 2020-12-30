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
	times := strings.Split(string(b), "-")
	if len(times) != 2 {
		return fmt.Errorf("timerange can only have two components, found %d in '%s'", len(times), string(b))
	}

	start, err := time.Parse("15:04", strings.TrimSpace(times[0]))
	if err != nil {
		return errors.Wrap(err, "cannot parse range start date")
	}

	end, err := time.Parse("15:04", strings.TrimSpace(times[1]))
	if err != nil {
		return errors.Wrap(err, "cannot parse range start date")
	}

	if !start.Before(end) {
		return fmt.Errorf("start date must be before end date for range '%s'", string(b))
	}

	*r = Timerange{
		Start: start,
		End:   end,
	}

	return nil
}

func (r *Timerange) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Timerange) String() string {
	layout := "15:04"
	return fmt.Sprintf("%s-%s", r.Start.Format(layout), r.End.Format(layout))
}

func (r *Timerange) Active(now time.Time) bool {
	hour, min, sec := now.Clock()
	timeNow := time.Time{}.Add(time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second)

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
