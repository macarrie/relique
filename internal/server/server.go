package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var srv *http.Server

func Start(debug bool, bindAddr string, port int, certPath string, keyPath string) error {
	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}

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
		Addr:         fmt.Sprintf("%s:%d", bindAddr, port),
		Handler:      router,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	slog.With(
		slog.Int("port", port),
	).Info("Starting HTTP server")

	if _, err := os.Lstat(certPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		slog.With(
			slog.Any("error", err),
			slog.String("file", certPath),
		).Error("Cannot find SSL certificate file")
		os.Exit(1)
	}
	if _, err := os.Lstat(keyPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		slog.With(
			slog.Any("error", err),
			slog.String("file", keyPath),
		).Error("Cannot find SSL certificate key file")
		os.Exit(1)
	}

	if err := srv.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
		slog.With(
			slog.Any("error", err),
		).Error("Cannot start HTTP server")
		os.Exit(1)
	}

	return nil
}

func Stop() error {
	slog.Info("Gracefully shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.With(slog.Any("error", err)).Error("Cannot gracefully shut down HTTP server. The server has been asked less politely to stop.")
		os.Exit(1)
	}

	return nil
}
