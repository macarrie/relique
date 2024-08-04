package api

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/config"
)

func ClientCreate(name string, address string) error {
	cl := client.New(name, address)
	if err := cl.Write(config.GetClientsCfgPath()); err != nil {
		return err
	}

	return nil
}

func ClientList() []client.Client {
	return config.Current.Clients
}

func ClientGet(name string) (client.Client, error) {
	clientList := ClientList()

	for _, cl := range clientList {
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
