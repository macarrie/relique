package schedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/pkg/errors"
)

type Timeranges struct {
	Ranges []Timerange `json:"ranges"`
}

func (r *Timeranges) UnmarshalText(b []byte) error {
	rangesStr := strings.Split(string(b), ",")
	var ranges []Timerange
	for _, str := range rangesStr {
		var rng Timerange
		err := rng.UnmarshalText([]byte(str))
		if err != nil {
			return errors.Wrap(err, "cannot parse range")
		}

		ranges = append(ranges, rng)
	}

	*r = Timeranges{Ranges: ranges}

	return nil
}

func (r *Timeranges) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Timeranges) String() string {
	var arr []string
	for _, timerange := range r.Ranges {
		arr = append(arr, timerange.String())
	}

	return strings.Join(arr, ",")
}

func (r *Timeranges) Active(now time.Time) bool {
	for _, timerange := range r.Ranges {
		if timerange.Active(now) {
			return true
		}
	}

	return false
}

func (r *Timeranges) Valid() error {
	var objErrors *multierror.Error

	for _, timerange := range r.Ranges {
		if err := timerange.Valid(); err != nil {
			objErrors = multierror.Append(objErrors, errors.Wrap(err, fmt.Sprintf("invalid range '%s'", timerange.String())))
		}
	}

	return objErrors.ErrorOrNil()
}
