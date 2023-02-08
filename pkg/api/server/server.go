// API Methods used by server daemon
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	consts "github.com/macarrie/relique/internal/types"

	config "github.com/macarrie/relique/internal/types/config/server_daemon_config"
	"github.com/macarrie/relique/pkg/api/utils"

	log "github.com/macarrie/relique/internal/logging"
	client "github.com/macarrie/relique/internal/types/client"
	"github.com/pkg/errors"
)

func GetConfigVersion(cl *client.Client) (string, error) {
	log.WithFields(log.Fields{
		"client": cl.Name,
	}).Debug("Checking client configuration version")

	response, err := utils.PerformRequest(config.Config, cl.Address, cl.Port, "GET", "/api/v1/config/version", nil)
	if err != nil {
		cl.APIAlive = consts.CRITICAL
		return "", errors.Wrap(err, "error when performing api request")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		cl.APIAlive = consts.CRITICAL
		return "", errors.Wrap(err, "cannot read response body from api requets")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var configVersion struct {
			Version string
		}
		if err := json.Unmarshal(body, &configVersion); err != nil {
			cl.APIAlive = consts.UNKNOWN
			return "", errors.Wrap(err, "cannot parse config version returned from client")
		}

		cl.APIAlive = consts.OK
		return configVersion.Version, nil
	}

	cl.APIAlive = consts.CRITICAL
	return "", fmt.Errorf("cannot get client version, status code '%d'", response.StatusCode)
}

func SendConfiguration(cl *client.Client) error {
	version, err := GetConfigVersion(cl)
	if err != nil {
		return errors.Wrap(err, "cannot get current config version for client")
	}

	if version == config.Config.Version {
		cl.APIAlive = consts.OK
		return nil
	}
	log.WithFields(log.Fields{
		"client": cl.Name,
	}).Info("Send configuration to client")

	cl.Version = config.Config.Version

	response, err := utils.PerformRequest(
		config.Config,
		cl.Address,
		cl.Port,
		"POST",
		"/api/v1/config",
		cl)
	if err != nil {
		cl.APIAlive = consts.CRITICAL
		return errors.Wrap(err, "error when performing api request")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		cl.APIAlive = consts.OK
		return nil
	}

	cl.APIAlive = consts.CRITICAL
	return nil
}
