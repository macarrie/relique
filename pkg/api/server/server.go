// API Methods used by server daemon
package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	config "github.com/macarrie/relique/internal/types/config/server_daemon_config"
	"github.com/macarrie/relique/pkg/api/utils"

	log "github.com/macarrie/relique/internal/logging"
	client "github.com/macarrie/relique/internal/types/client"
	"github.com/pkg/errors"
)

func GetConfigVersion(client client.Client) (string, error) {
	log.WithFields(log.Fields{
		"client": client.Name,
	}).Debug("Checking client configuration version")

	response, err := utils.PerformRequest(config.Config, client.Address, client.Port, "GET", "/api/v1/config/version", nil)
	if err != nil {
		return "", errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "cannot read response body from api requets")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var configVersion struct {
			Version string
		}
		if err := json.Unmarshal(body, &configVersion); err != nil {
			return "", errors.Wrap(err, "cannot parse config version returned from client")
		}

		return configVersion.Version, nil
	}

	return "", fmt.Errorf("cannot get client version, status code '%d'", response.StatusCode)
}

func SendConfiguration(client client.Client) error {
	version, err := GetConfigVersion(client)
	if err != nil {
		return errors.Wrap(err, "cannot get urrent config version for client")
	}

	if version != config.Config.Version {
		log.WithFields(log.Fields{
			"client": client.Name,
		}).Info("Send configuration to client")

		client.Version = config.Config.Version
		client.ServerAddress = config.Config.PublicAddress
		client.ServerPort = config.Config.Port

		response, err := utils.PerformRequest(config.Config, client.Address, client.Port, "POST", "/api/v1/config", client)
		if err != nil {
			return errors.Wrap(err, "error when performing api request")
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			return nil
		}
	}

	return nil
}
