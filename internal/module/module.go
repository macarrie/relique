package module

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/pelletier/go-toml"
)

func LoadFromFile(file string) (m Module, err error) {
	slog.With(
		slog.String("path", file),
	).Debug("Loading module configuration parameters from file")

	f, err := os.Open(file)
	if err != nil {
		return Module{}, fmt.Errorf("cannot open file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("cannot close file correctly: %w", cerr)
		}
	}()

	content, _ := io.ReadAll(f)

	var module Module
	if err := toml.Unmarshal(content, &module); err != nil {
		return Module{}, fmt.Errorf("cannot parse toml file: %w", err)
	}

	if err := module.Valid(); err != nil {
		return Module{}, fmt.Errorf("invalid module loaded from file: %w", err)
	}

	return module, nil
}

func GetByName(list []Module, name string) (Module, error) {
	for _, mod := range list {
		if mod.Name == name {
			return mod, nil
		}
	}
	return Module{}, fmt.Errorf("cannot find module named '%s'", name)
}
