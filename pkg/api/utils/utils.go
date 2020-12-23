package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	log "github.com/macarrie/relique/internal/logging"

	config "github.com/macarrie/relique/internal/types/config/common"
)

func getBaseURL(addr string, port uint32) string {
	return fmt.Sprintf("https://%s:%d/", addr, port)
}

func PerformRequest(config config.Configuration, addr string, port uint32, method string, path string, paramsInt interface{}) (http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.StrictSSLCertificateCheck},
	}
	httpClient := &http.Client{
		Transport: tr,
		// TODO: Put timeout value in const var
		Timeout: time.Duration(60 * time.Second),
	}

	var request *http.Request
	var err error

	urlObject, _ := url.ParseRequestURI(getBaseURL(addr, port))
	urlObject.Path = path
	if paramsInt == nil {
		paramsInt = map[string]string{}
	}

	if method == "POST" || method == "PUT" || method == "GET" {
		jsonParams, _ := json.Marshal(paramsInt)

		request, err = http.NewRequest(method, urlObject.String(), bytes.NewReader(jsonParams))
		if err != nil {
			return http.Response{}, errors.Wrap(err, "cannot create http request object")
		}
	} else {
		log.WithFields(log.Fields{
			"method": method,
		}).Error("HTTP method unknown")

		return http.Response{}, fmt.Errorf("unknown HTTP method '%s'", method)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, errors.Wrap(err, "cannot perform http request")
	}

	return *response, nil
}
func SendMultipart(config config.Configuration, addr string, port uint32, path string, contentType string, r io.Reader) (http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.StrictSSLCertificateCheck},
	}
	httpClient := &http.Client{
		Transport: tr,
		// TODO: Put timeout value in const var
		// Longer timeout for multipart file send
		Timeout: time.Duration(15 * time.Minute),
	}

	var request *http.Request
	var err error

	urlObject, _ := url.ParseRequestURI(getBaseURL(addr, port))
	urlObject.Path = path

	request, err = http.NewRequest("POST", urlObject.String(), r)
	if err != nil {
		return http.Response{}, errors.Wrap(err, "cannot create http request object")
	}

	request.Header.Set("Content-Type", contentType)
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, errors.Wrap(err, "cannot perform http request")
	}

	return *response, nil
}
