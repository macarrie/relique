package api

import (
	"fmt"

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
