package schedule

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config"
)

type Schedule struct {
	Name      string
	Monday    string
	Tuesday   string
	Wednesday string
	Thursday  string
	Friday    string
	Saturday  string
	Sunday    string
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

	byteValue, _ := ioutil.ReadAll(tomlFile)

	var schedule Schedule
	if err := toml.Unmarshal(byteValue, &schedule); err != nil {
		return Schedule{}, err
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
			}).Warn("Cannot load client configuration from file")
			return err
		}

		if filepath.Ext(path) == ".toml" {
			files = append(files, path)
		}
		return nil
	})

	var schedules []Schedule
	for _, file := range files {
		client, err := loadFromFile(file)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": file,
			}).Warn("Cannot load client configuration from file")
			continue
		}

		schedules = append(schedules, client)
	}

	return schedules, nil
}
