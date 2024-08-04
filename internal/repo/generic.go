package repo

import (
	"fmt"
)

type GenericRepository struct {
	Name    string `json:"name" toml:"name"`
	Type    string `json:"type" toml:"type"`
	Default bool   `json:"default" toml:"default"`
}

func (g *GenericRepository) GetName() string {
	return g.Name
}

func (g *GenericRepository) GetType() string {
	return g.Type
}

func (g *GenericRepository) Write(path string) error {
	return fmt.Errorf("TODO")
}

func (g *GenericRepository) IsDefault() bool {
	return g.Default
}
