package client

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"

	"github.com/macarrie/relique/internal/db"
	"github.com/pkg/errors"

	"github.com/macarrie/relique/internal/types/module"

	"github.com/spf13/viper"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config"
)

type Client struct {
	ID            int64
	Name          string `mapstructure:"name"`
	Address       string `mapstructure:"address"`
	Port          uint32 `mapstructure:"port"`
	Modules       []module.Module
	Version       string
	ServerAddress string `mapstructure:"server_address"`
	ServerPort    uint32 `mapstructure:"server_port"`
}

func (c *Client) String() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.Address)
}

func loadFromFile(file string) (Client, error) {
	log.WithFields(log.Fields{
		"path": file,
	}).Debug("Loading client configuration from file")

	clientViper := viper.New()
	clientViper.SetConfigType("toml")
	clientViper.SetConfigFile(file)

	if err := clientViper.ReadInConfig(); err != nil {
		return Client{}, err
	}

	var client Client
	if err := clientViper.Unmarshal(&client); err != nil {
		return Client{}, err
	}

	for i := range client.Modules {
		client.Modules[i].ComputeBackupTypeFromString()
		if err := client.Modules[i].LoadDefaultConfiguration(); err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"module": client.Modules[i].ModuleType,
			}).Fatal("Cannot find default configuration parameters for module. Make sure that this module is correctly installed")
		}
	}

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
			}).Warn("Cannot load client configuration from file")
			return err
		}

		if filepath.Ext(path) == ".toml" {
			files = append(files, path)
		}
		return nil
	})

	var clients []Client
	for _, file := range files {
		client, err := loadFromFile(file)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": file,
			}).Warn("Cannot load client configuration from file")
			continue
		}

		clients = append(clients, client)
	}

	return clients, nil
}

func GetID(name string) (int64, error) {
	request := sq.Select(
		"id",
	).From(
		"clients",
	).Where(
		"name = ?",
		name,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	row := db.Read().QueryRow(query, args...)
	defer db.RUnlock()

	var id int64
	if err := row.Scan(&id); err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, errors.Wrap(err, "cannot search retrieve client ID in db")
	}

	return id, nil
}

func GetByID(id int64) (Client, error) {
	log.WithFields(log.Fields{
		"id": id,
	}).Trace("Looking for client in database")

	request := sq.Select(
		"id",
		"config_version",
		"name",
		"address",
		"port",
		"server_address",
		"server_port",
	).From(
		"clients",
	).Where(
		"id = ?",
		id,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return Client{}, errors.Wrap(err, "cannot build sql query")
	}

	row := db.Read().QueryRow(query, args...)
	defer db.RUnlock()

	var cl Client
	if err := row.Scan(&cl.ID,
		&cl.Version,
		&cl.Name,
		&cl.Address,
		&cl.Port,
		&cl.ServerAddress,
		&cl.ServerPort,
	); err == sql.ErrNoRows {
		return Client{}, nil
	} else if err != nil {
		return Client{}, errors.Wrap(err, "cannot retrieve client from db")
	}

	return cl, nil
}

func (c *Client) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"name":    c.Name,
		"address": c.Address,
		"id":      c.ID,
	})
}

func (c *Client) Save() (int64, error) {
	id, err := GetID(c.Name)
	if err != nil {
		return 0, errors.Wrap(err, "cannot search for possibly existing client ID")
	}

	if id != 0 {
		c.ID = id
		return c.Update()
	}

	c.GetLog().Debug("Saving client into database")

	request := sq.Insert("clients").Columns(
		"config_version",
		"name",
		"address",
		"port",
		"server_address",
		"server_port",
	).Values(
		c.Version,
		c.Name,
		c.Address,
		c.Port,
		c.ServerAddress,
		c.ServerPort,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	result, err := db.Write().Exec(
		query,
		args...,
	)
	defer db.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "cannot save client into db")
	}

	c.ID, _ = result.LastInsertId()

	return c.ID, nil
}

func (c *Client) Update() (int64, error) {
	c.GetLog().Debug("Updating client details into database")

	if c.ID == 0 {
		return 0, fmt.Errorf("cannot update client with ID 0")
	}

	request := sq.Update("clients").SetMap(sq.Eq{
		"config_version": c.Version,
		"name":           c.Name,
		"address":        c.Address,
		"port":           c.Port,
		"server_address": c.ServerAddress,
		"server_port":    c.ServerPort,
	}).Where(
		" id = ?",
		c.ID,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	_, err = db.Write().Exec(query, args...)
	defer db.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "cannot update client into db")
	}

	return c.ID, nil
}
