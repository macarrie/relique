package client

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	consts "github.com/macarrie/relique/internal/types"

	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/pkg/errors"

	"github.com/macarrie/relique/internal/types/module"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config"
	"github.com/pelletier/go-toml"
)

type Client struct {
	Name            string          `json:"name" toml:"name"`
	Address         string          `json:"address" toml:"address"`
	Port            uint32          `json:"port" toml:"port"`
	Modules         []module.Module `json:"modules"`
	Version         string          `json:"version"`
	ServerAddress   string          `json:"server_address" toml:"server_address"`
	ServerPort      uint32          `json:"server_port" toml:"server_port"`
	APIAlive        uint8           `json:"api_alive"`
	SSHAlive        uint8           `json:"ssh_alive"`
	SSHAliveMessage string          `json:"ssh_alive_message"`
}

func (c *Client) String() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.Address)
}

func LoadFromFile(file string) (cl Client, err error) {
	log.WithFields(log.Fields{
		"path": file,
	}).Debug("Loading client configuration from file")

	f, err := os.Open(file)
	if err != nil {
		return Client{}, errors.Wrap(err, "cannot open file")
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = errors.Wrap(cerr, "cannot close file correctly")
		}
	}()

	content, _ := io.ReadAll(f)

	var client Client
	if err := toml.Unmarshal(content, &client); err != nil {
		return Client{}, errors.Wrap(err, "cannot parse toml file")
	}

	modules := client.Modules
	var filteredModulesList []module.Module
	for i := range modules {
		if err := modules[i].LoadDefaultConfiguration(); err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"module": client.Modules[i].ModuleType,
			}).Error("Cannot find default configuration parameters for module. Make sure that this module is correctly installed")
		}
		if err := modules[i].Valid(); err == nil {
			filteredModulesList = append(filteredModulesList, modules[i])
		} else {
			modules[i].GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Module has invalid configuration. This module will not be loaded into configuration")
		}

	}
	client.Modules = filteredModulesList

	return client, nil
}

func LoadFromPath(p string) ([]Client, error) {
	absPath := config.GetConfigurationSubpath(p)

	var files []string

	_ = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": path,
			}).Error("Cannot load client configuration from file")
			return err
		}

		if filepath.Ext(path) == ".toml" {
			files = append(files, path)
		}
		return nil
	})

	var clients []Client
	for _, file := range files {
		client, err := LoadFromFile(file)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": file,
			}).Error("Cannot load client configuration from file")
			continue
		}

		if client.Valid() {
			clients = append(clients, client)
		} else {
			client.GetLog().Error("Client has invalid configuration. This client will not be loaded into configuration")
		}
	}

	return clients, nil
}

func FillSchedulesStruct(clients []Client, schedules []schedule.Schedule) ([]Client, error) {
	var retList []Client
	for _, client := range clients {
		var mods []module.Module
		for _, mod := range client.Modules {
			var scheds []schedule.Schedule
			for _, scheduleName := range mod.ScheduleNames {
				foundScheduleDef := false
				for _, s := range schedules {
					if s.Name == scheduleName {
						foundScheduleDef = true
						scheds = append(scheds, s)
					}
				}
				if !foundScheduleDef {
					return []Client{}, fmt.Errorf("cannot find schedule '%s' definition for module '%s' of client '%s'", scheduleName, mod.Name, client.Name)
				}
			}
			mod.Schedules = scheds
			mods = append(mods, mod)
		}
		client.Modules = mods
		retList = append(retList, client)
	}

	return retList, nil
}

func FillServerPublicAddress(clients []Client, addr string, port uint32) []Client {
	var retList []Client
	for _, client := range clients {
		client.ServerAddress = addr
		client.ServerPort = port
		retList = append(retList, client)
	}

	return retList
}

func InitAliveStatus(clients []Client) []Client {
	var retList []Client
	for _, client := range clients {
		client.SSHAlive = consts.UNKNOWN
		client.APIAlive = consts.UNKNOWN
		retList = append(retList, client)
	}

	return retList
}

func FillConfigVersion(clients []Client, version string) []Client {
	var retList []Client
	for _, client := range clients {
		client.Version = version
		retList = append(retList, client)
	}

	return retList
}

func (c *Client) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"name":    c.Name,
		"address": c.Address,
	})
}

func (c *Client) Valid() bool {
	if c.Name == "" || c.Address == "" {
		return false
	}

	return true
}
