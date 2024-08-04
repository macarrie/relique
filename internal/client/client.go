package client

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

func New(name string, address string) Client {
	return Client{
		Name:    name,
		Address: address,
	}
}

func LoadFromFile(file string) (cl Client, err error) {
	slog.Debug("Loading client configuration from file", slog.String("path", file))

	f, err := os.Open(file)
	if err != nil {
		return Client{}, fmt.Errorf("cannot open file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("cannot close file correctly: %w", err)
		}
	}()

	content, _ := io.ReadAll(f)

	var client Client
	if err := toml.Unmarshal(content, &client); err != nil {
		return Client{}, fmt.Errorf("cannot parse toml file: %w", err)
	}

	return client, nil
}

func LoadFromPath(p string) ([]Client, error) {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("specified path does not exist or cannot be opened: %w", err)
	}

	var files []string

	_ = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.With(
				slog.Any("error", err),
				slog.String("path", path),
			).Error("Cannot load client configuration from file")
			return err
		}

		if filepath.Ext(path) == ".toml" {
			files = append(files, path)
		}
		return nil
	})

	var clients []Client
	for _, file := range files {
		cl, err := LoadFromFile(file)
		if err != nil {
			slog.With(
				slog.Any("err", err),
				slog.String("path", file),
			).Error("Cannot load client configuration from file")
			continue
		}

		if cl.Valid() {
			clients = append(clients, cl)
		} else {
			cl.GetLog().Error("Client has invalid configuration. This client will not be loaded into configuration")
		}
	}

	return clients, nil
}
