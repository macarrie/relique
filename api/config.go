package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/module"
)

var ReliqueVersion string

func ConfigGet() (config.Configuration, error) {
	if !config.Loaded {
		if err := config.Load("relique"); err != nil {
			return config.Configuration{}, fmt.Errorf("cannot load relique configuration: %w", err)
		}
	}
	return config.Current, nil
}

func ConfigGetVersion() string {
	if ReliqueVersion == "" {
		return "unknown"
	}
	return ReliqueVersion
}

func createSelfSignedSSLCerts(certPath string, keyPath string) error {
	var priv *rsa.PrivateKey
	var err error

	var KEY_SIZE int = 2048
	var CERT_HOSTS string = "relique"

	priv, err = rsa.GenerateKey(rand.Reader, KEY_SIZE)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
	// the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.

	notBefore := time.Now()
	notAfter := notBefore.Add(3650 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"relique"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(CERT_HOSTS, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &(priv.PublicKey), priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("error closing cert file: %v", err)
	}
	slog.With(
		slog.String("path", certPath),
	).Info("Wrote certificate file")

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key file for writing: %v", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("error closing file: %v", err)
	}
	slog.With(
		slog.String("path", keyPath),
	).Info("Wrote certificate key file")

	return nil
}

func ConfigInit(cfgPath string, modPath string, repoPath string, catalogPath string) error {
	configPath := cfgPath
	moduleInstallPath := modPath
	repoStoragePath := repoPath
	catalogStoragePath := catalogPath
	if cfgPath == "" {
		configPath = "/etc/relique"
		if modPath == "" {
			moduleInstallPath = "/var/lib/relique/modules"
		}
		if repoPath == "" {
			repoStoragePath = "/var/lib/relique/storage"
		}
		if catalogPath == "" {
			catalogStoragePath = "/var/lib/relique/catalog"
		}
	}
	if moduleInstallPath == "" {
		moduleInstallPath = filepath.Clean(fmt.Sprintf("%s/modules", configPath))
	}
	if repoStoragePath == "" {
		repoStoragePath = filepath.Clean(fmt.Sprintf("%s/storage", configPath))
	}
	if catalogStoragePath == "" {
		catalogStoragePath = filepath.Clean(fmt.Sprintf("%s/catalog", configPath))
	}

	configPath = filepath.Clean(configPath)
	moduleInstallPath = filepath.Clean(moduleInstallPath)
	certsPath := filepath.Clean(fmt.Sprintf("%s/certs/", configPath))
	certFilePath := filepath.Clean(fmt.Sprintf("%s/cert.pem", certsPath))
	keyFilePath := filepath.Clean(fmt.Sprintf("%s/key.pem", certsPath))

	// Check if config folder already exists
	if _, err := os.Stat(configPath); err == nil && !os.IsNotExist(err) {
		return fmt.Errorf("specified folder '%s' already exists, aborting config init to avoid overwriting existing configuration", configPath)
	}

	// Check if module install folder already exists
	if _, err := os.Stat(moduleInstallPath); err == nil && !os.IsNotExist(err) {
		return fmt.Errorf("specified folder '%s' already exists, aborting config init to avoid overwriting existing module install folder", moduleInstallPath)
	}

	// Create config folder
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("cannot create folder '%s': %w", configPath, err)
	}
	slog.With(
		slog.String("path", configPath),
	).Info("Created default relique configuration folder")

	// Create db folder
	dbPath := fmt.Sprintf("%s/%s", configPath, config.DB_DEFAULT_FOLDER)
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return fmt.Errorf("cannot create folder '%s': %w", dbPath, err)
	}
	slog.With(
		slog.String("path", dbPath),
	).Info("Created default relique database folder")

	// Create module config folder
	if err := os.MkdirAll(moduleInstallPath, 0755); err != nil {
		return fmt.Errorf("cannot create folder '%s': %w", moduleInstallPath, err)
	}
	slog.With(
		slog.String("path", moduleInstallPath),
	).Info("Created default relique module install folder")

	configFilePath := filepath.Clean(fmt.Sprintf("%s/relique.toml", configPath))
	config.New()
	config.Current.ModuleInstallPath = moduleInstallPath
	config.Current.WebUI.SSLCert = certFilePath
	config.Current.WebUI.SSLKey = keyFilePath
	module.MODULES_INSTALL_PATH = moduleInstallPath
	if err := config.Write(configFilePath); err != nil {
		return fmt.Errorf("cannot create default configuration file: %w", err)
	}
	slog.With(
		slog.String("path", configFilePath),
	).Info("Created relique configuration file")

	// Create self signed certificates for webui
	if err := os.MkdirAll(certsPath, 0755); err != nil {
		return fmt.Errorf("cannot create certificates folder '%s': %w", certsPath, err)
	}
	slog.With(
		slog.String("path", certsPath),
	).Debug("Created certificates folder")
	if err := createSelfSignedSSLCerts(certFilePath, keyFilePath); err != nil {
		return fmt.Errorf("cannot create self signed certificates: %w", err)
	}

	// Create modules folder
	if err := os.MkdirAll(moduleInstallPath, 0755); err != nil {
		return fmt.Errorf("cannot create module install folder '%s': %w", moduleInstallPath, err)
	}
	slog.With(
		slog.String("path", moduleInstallPath),
	).Debug("Created modules install folder")

	// Install default modules
	if err := ModuleInstall(moduleInstallPath, "https://github.com/macarrie/relique-module-generic", false, false, false); err != nil {
		return fmt.Errorf("cannot install default generic module: %w", err)
	}
	slog.Debug("Installed default generic module")

	// Create clients folder
	clientsFolder := config.GetClientsCfgPath()
	if err := os.Mkdir(clientsFolder, 0755); err != nil {
		return fmt.Errorf("cannot create clients folder '%s': %w", clientsFolder, err)
	}
	slog.With(
		slog.String("path", clientsFolder),
	).Info("Created clients configuration folder")

	if err := ClientCreate("local", "localhost"); err != nil {
		return fmt.Errorf("cannot create default client: %w", err)
	}
	// TODO: Add example module to local client

	// Create catalog config folder
	catalogFolder := config.GetCatalogCfgPath()
	if err := os.Mkdir(catalogFolder, 0755); err != nil {
		return fmt.Errorf("cannot create catalog folder '%s': %w", catalogFolder, err)
	}
	slog.With(
		slog.String("path", catalogFolder),
	).Info("Created catalog folder")

	// Create repositories config folder
	reposFolder := config.GetReposCfgPath()
	if err := os.Mkdir(reposFolder, 0755); err != nil {
		return fmt.Errorf("cannot create repositories folder '%s': %w", reposFolder, err)
	}
	slog.With(
		slog.String("path", reposFolder),
	).Info("Created repositories configuration folder")

	// Create default repo
	if err := RepoCreateLocal("local", repoStoragePath, true); err != nil {
		return fmt.Errorf("cannot create default repository: %w", err)
	}

	return nil
}
