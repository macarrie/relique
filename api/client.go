package api

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/samber/lo"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/module"
)

func ClientCreate(name string, address string) error {
	cl := client.New(name, address)
	if err := cl.Write(config.GetClientsCfgPath()); err != nil {
		return err
	}

	return nil
}

func ClientList(p api_helpers.PaginationParams, s api_helpers.ClientSearch) api_helpers.PaginatedResponse[client.Client] {
	limit := p.Limit
	clientList := config.Current.Clients
	// Filters
	if s.ModuleName != "" {
		clientList = lo.Filter(clientList, func(item client.Client, index int) bool {
			mods := lo.Filter(item.Modules, func(m module.Module, index int) bool {
				return m.Name == s.ModuleName
			})
			return len(mods) != 0
		})
	}
	if s.ModuleType != "" {
		clientList = lo.Filter(clientList, func(item client.Client, index int) bool {
			mods := lo.Filter(item.Modules, func(m module.Module, index int) bool {
				return m.ModuleType == s.ModuleType
			})
			return len(mods) != 0
		})
	}

	// Count after filters
	count := len(clientList)
	if limit != 0 {
		clientList = lo.Slice(clientList, 0, int(p.Limit))
	}
	return api_helpers.PaginatedResponse[client.Client]{
		Count:      uint64(count),
		Pagination: p,
		Data:       clientList,
	}
}

func ClientGet(name string) (client.Client, error) {
	clientList := ClientList(api_helpers.PaginationParams{}, api_helpers.ClientSearch{})
	for _, cl := range clientList.Data {
		if cl.Name == name {
			return cl, nil
		}
	}

	return client.Client{}, fmt.Errorf("cannot find client '%s' in configuration", name)
}

func ClientSSHPing(c client.Client) error {
	c.GetLog().Debug("Checking SSH connexion with client")

	var sshUser string
	if c.SSHUser != "" {
		sshUser = c.SSHUser
	} else {
		sshUser = client.DEFAULT_SSH_USER
	}

	var sshPort int
	if c.SSHPort != 0 {
		sshPort = c.SSHPort
	} else {
		sshPort = client.DEFAULT_SSH_PORT
	}

	sshPingCmd := exec.Command("ssh", "-f", "-o BatchMode=yes", "-p", fmt.Sprint(sshPort), fmt.Sprintf("%s@%s", sshUser, c.Address), "echo 'ping'")
	slog.With(
		slog.String("cmd", sshPingCmd.String()),
		slog.String("client", c.Name),
	).Debug("Trying to ping client with following command")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	sshPingCmd.Stdout = &stdout
	sshPingCmd.Stderr = &stderr

	if err := sshPingCmd.Run(); err != nil {
		return fmt.Errorf("cannot ping client via ssh: %s", stderr.String())
	}

	if stderr.String() != "" {
		return fmt.Errorf("cannot ping client via ssh:'%s'", stderr.String())
	}

	c.GetLog().Info("Client SSH ping successful")
	return nil
}
