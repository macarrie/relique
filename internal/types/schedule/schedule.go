package schedule

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config"
	"github.com/pelletier/go-toml"
)

type Schedule struct {
	Name      string     `json:"name"`
	Monday    Timeranges `json:"monday"`
	Tuesday   Timeranges `json:"tuesday"`
	Wednesday Timeranges `json:"wednesday"`
	Thursday  Timeranges `json:"thursday"`
	Friday    Timeranges `json:"friday"`
	Saturday  Timeranges `json:"saturday"`
	Sunday    Timeranges `json:"sunday"`
}

func loadFromFile(file string) (Schedule, error) {
	log.WithFields(log.Fields{
		"path": file,
	}).Debug("Loading schedule configuration from file")

	tomlFile, err := os.Open(file)
	if err != nil {
		return Schedule{}, err
	}
	defer tomlFile.Close()

	byteValue, _ := io.ReadAll(tomlFile)

	var schedule Schedule
	if err := toml.Unmarshal(byteValue, &schedule); err != nil {
		return Schedule{}, err
	}

	if err := schedule.Valid(); err != nil {
		return Schedule{}, errors.Wrap(err, "invalid schedule loaded from file")
	}

	return schedule, nil
}

func LoadFromPath(p string) ([]Schedule, error) {
	absPath := config.GetConfigurationSubpath(p)

	var files []string

	_ = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": path,
			}).Error("Cannot walk path to load schedules")
			return err
		}

		if filepath.Ext(path) == ".toml" {
			files = append(files, path)
		}
		return nil
	})

	var schedules []Schedule
	for _, file := range files {
		sched, err := loadFromFile(file)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": file,
			}).Error("Cannot load schedule configuration from file")
			continue
		}

		schedules = append(schedules, sched)
	}

	return schedules, nil
}

func (s *Schedule) Active(now time.Time) bool {
	var rangeToCheck Timeranges
	switch now.Weekday() {
	case time.Monday:
		rangeToCheck = s.Monday
	case time.Tuesday:
		rangeToCheck = s.Tuesday
	case time.Wednesday:
		rangeToCheck = s.Wednesday
	case time.Thursday:
		rangeToCheck = s.Thursday
	case time.Friday:
		rangeToCheck = s.Friday
	case time.Saturday:
		rangeToCheck = s.Saturday
	case time.Sunday:
		rangeToCheck = s.Sunday
	}

	return rangeToCheck.Active(now)
}

func (s *Schedule) Valid() error {
	var objErrors *multierror.Error

	if s.Name == "" {
		objErrors = multierror.Append(objErrors, fmt.Errorf("missing name"))
	}
	if err := s.Monday.Valid(); err != nil {
		objErrors = multierror.Append(objErrors, errors.Wrap(err, "invalid monday ranges"))
	}
	if err := s.Tuesday.Valid(); err != nil {
		objErrors = multierror.Append(objErrors, errors.Wrap(err, "invalid tuesday ranges"))
	}
	if err := s.Wednesday.Valid(); err != nil {
		objErrors = multierror.Append(objErrors, errors.Wrap(err, "invalid wednesday ranges"))
	}
	if err := s.Thursday.Valid(); err != nil {
		objErrors = multierror.Append(objErrors, errors.Wrap(err, "invalid thursday ranges"))
	}
	if err := s.Friday.Valid(); err != nil {
		objErrors = multierror.Append(objErrors, errors.Wrap(err, "invalid friday ranges"))
	}
	if err := s.Saturday.Valid(); err != nil {
		objErrors = multierror.Append(objErrors, errors.Wrap(err, "invalid saturday ranges"))
	}
	if err := s.Sunday.Valid(); err != nil {
		objErrors = multierror.Append(objErrors, errors.Wrap(err, "invalid sunday ranges"))
	}

	return objErrors.ErrorOrNil()
}
