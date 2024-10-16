package client

import (
	"cmp"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/pelletier/go-toml"

	"github.com/macarrie/relique/internal/module"
)

var DEFAULT_SSH_USER string = "relique"
var DEFAULT_SSH_PORT int = 22

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

	if err := toml.Unmarshal(content, &cl); err != nil {
		return Client{}, fmt.Errorf("cannot parse toml file: %w", err)
	}

	modules := cl.Modules
	var filteredModulesList []module.Module
	for i := range modules {
		if err := modules[i].LoadDefaultConfiguration(); err != nil {
			modules[i].GetLog().With(
				slog.Any("error", err),
				slog.String("module_type", modules[i].ModuleType),
			).Error("Cannot find default configuration parameters for module. Make sure that this module is correctly installed")
			continue
		}
		if err := modules[i].Valid(); err == nil {
			filteredModulesList = append(filteredModulesList, modules[i])
		} else {
			modules[i].GetLog().With(
				slog.Any("error", err),
				slog.Any("client", cl.Name),
			).Error("Module has invalid configuration. This module will not be loaded into configuration")
		}

	}
	cl.Modules = filteredModulesList

	return cl, nil
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
			).Error("Cannot browse client configuration")
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

	// Sort by client_name to have consistent output
	slices.SortFunc(clients, func(a, b Client) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return clients, nil
}
