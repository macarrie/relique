package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/macarrie/relique/internal/types/config/client_daemon_config"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
)

var srv *http.Server

func Start() error {
	gin.SetMode(gin.ReleaseMode)

	router := getRoutes()
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", client_daemon_config.Config.BindAddr, client_daemon_config.Config.Port),
		Handler:      router,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.WithFields(log.Fields{
		"port": client_daemon_config.Config.Port,
	}).Info("Starting HTTP server")

	if _, err := os.Lstat(client_daemon_config.Config.SSLCert); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.WithFields(log.Fields{
			"error": err,
			"file":  client_daemon_config.Config.SSLCert,
		}).Fatal("Cannot find SSL certificate file")
		return err
	}
	if _, err := os.Lstat(client_daemon_config.Config.SSLKey); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.WithFields(log.Fields{
			"error": err,
			"file":  client_daemon_config.Config.SSLKey,
		}).Fatal("Cannot find SSL certificate file")
		return err
	}

	if err := srv.ListenAndServeTLS(client_daemon_config.Config.SSLCert, client_daemon_config.Config.SSLKey); err != nil && err != http.ErrServerClosed {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot start HTTP server")
		return err
	}

	return nil
}

func Stop() error {
	log.Info("Gracefully shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Cannot gracefully shut down HTTP server. The server has been asked less politely to stop.")
	}

	return nil
}
