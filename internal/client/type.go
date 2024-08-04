package client

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennygrant/sanitize"

	"github.com/macarrie/relique/internal/module"

	"github.com/pelletier/go-toml"
)

type Client struct {
	Name    string          `json:"name" toml:"name"`
	Address string          `json:"address" toml:"address"`
	SSHUser string          `json:"ssh_user" toml:"ssh_user"`
	SSHPort int             `json:"ssh_port" toml:"ssh_port"`
	Modules []module.Module `json:"modules" toml:"modules"`
}

func (c *Client) Write(rootPath string) error {
	if !c.Valid() {
		return fmt.Errorf("cannot write invalid client to file")
	}

	var path string = filepath.Clean(fmt.Sprintf("%s/%s.toml",
		rootPath,
		strings.ToLower(sanitize.Accents(sanitize.BaseName(c.Name))),
	))

	clientFile, clientErr := toml.Marshal(c)
	if clientErr != nil {
		return fmt.Errorf("cannot serialize client info to toml data: %w", clientErr)
	}
	if err := os.WriteFile(path, clientFile, 0644); err != nil {
		return fmt.Errorf("cannot export client info to file: %w", err)
	}

	c.GetLog().With(
		slog.String("path", path),
	).Debug("Saved client to file")
	return nil
}

func (c *Client) String() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.Address)
}

func (c *Client) GetLog() *slog.Logger {
	return slog.With(
		slog.String("name", c.Name),
		slog.String("address", c.Address),
	)
}

func (c *Client) Valid() bool {
	if c.Name == "" || c.Address == "" {
		return false
	}

	return true
}
