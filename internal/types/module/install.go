package module

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/pkg/errors"
)

func GetLocallyInstalled() ([]Module, error) {
	SetModulePathDefaultValue()

	log.WithFields(log.Fields{
		"path": MODULES_INSTALL_PATH,
	}).Info("Using modules install folder")

	items, err := ioutil.ReadDir(MODULES_INSTALL_PATH)
	if err != nil {
		return []Module{}, errors.Wrap(err, "cannot list installed modules from filesystem")
	}

	var modulesList []Module
	var errorList *multierror.Error
	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		mod := Module{
			Name:       item.Name(),
			ModuleType: item.Name(),
		}
		if err := mod.LoadDefaultConfiguration(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot read module default configuration")
			errorList = multierror.Append(errorList, err)
			continue
		}
		modulesList = append(modulesList, mod)
	}
	return modulesList, errorList.ErrorOrNil()
}

func IsInstalled(moduleName string) (bool, error) {
	installedModules, err := GetLocallyInstalled()
	if err != nil {
		return false, errors.Wrap(err, "cannot get locally installed modules")
	}

	for _, mod := range installedModules {
		if mod.Name == moduleName {
			return true, nil
		}
	}

	return false, nil
}

func gitClone(source string, dest string) error {
	log.WithFields(log.Fields{
		"source_path": source,
		"destination": dest,
	}).Info("Cloning git repository")

	var gitSuffix string
	if strings.HasSuffix(source, ".git") {
		gitSuffix = ""
	} else {
		gitSuffix = ".git"
	}

	var protocolPrefix string
	if strings.HasPrefix(source, "http") || strings.HasPrefix(source, "git@") {
		protocolPrefix = ""
	} else {
		protocolPrefix = "https://"
	}

	if _, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:      fmt.Sprintf("%s%s%s", protocolPrefix, source, gitSuffix),
		Progress: os.Stdout,
	}); err != nil {
		return errors.Wrap(err, "cannot clone git repository")
	}

	return nil
}

func downloadArchive(filepath string, url string) error {
	// Get the data
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		// TODO: Put timeout value in const var
		Timeout: time.Duration(60 * time.Second),
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return errors.Wrap(err, "cannot perform download request")
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "cannot create destination file")
	}
	defer out.Close()

	// Write the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return errors.Wrap(err, "cannot write downloaded archive to destination file")
	}

	return nil
}

func extractArchive(source string, dest string) error {
	log.WithFields(log.Fields{
		"source_path": source,
		"destination": dest,
	}).Info("Extracting module archive")

	if _, err := os.Lstat(source); os.IsNotExist(err) {
		return fmt.Errorf("cannot find archive file to extract")
	}

	extractCmd := exec.Command("tar", "-xvf", source, "-C", dest)

	var out bytes.Buffer
	extractCmd.Stdout = &out
	extractCmd.Stderr = &out

	err := extractCmd.Run()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot extract module archive: %s", out.String()))
	}

	// TODO: Chmod/chown files

	return nil
}

func installModuleFiles(source string, dest string, skipChown bool) error {
	log.WithFields(log.Fields{
		"source_path": source,
		"destination": dest,
	}).Info("Installing module files")

	copyCmd := exec.Command("cp", "-r", source, dest)

	var out bytes.Buffer
	copyCmd.Stdout = &out
	copyCmd.Stderr = &out

	if err := copyCmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot copy module files: %s", out.String()))
	}

	chmodCmd := exec.Command("chmod", "-R", "755", dest)

	var out2 bytes.Buffer
	chmodCmd.Stdout = &out2
	chmodCmd.Stderr = &out2

	if err := chmodCmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot chmod installed module files: %s", out2.String()))
	}

	if !skipChown {
		chownCmd := exec.Command("chown", "-R", ":relique", dest)

		var out3 bytes.Buffer
		chownCmd.Stdout = &out3
		chownCmd.Stderr = &out3

		if err := chownCmd.Run(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("cannot chown installed module files: %s", out3.String()))
		}
	}

	return nil
}

func Install(path string, local bool, archive bool, force bool, skipChown bool) error {
	SetModulePathDefaultValue()

	log.WithFields(log.Fields{
		"install_path": MODULES_INSTALL_PATH,
		"local":        local,
		"archive":      archive,
		"source_path":  path,
	}).Info("Installing relique module")

	if _, err := os.Lstat(MODULES_INSTALL_PATH); os.IsNotExist(err) {
		return errors.Wrap(err, "module install path does not exist. Please check that relique is correctly installed and that provided module install path is correct")
	}

	tempInstallFolder, err := ioutil.TempDir("", "relique-module-install-*")
	if err != nil {
		return errors.Wrap(err, "cannot create temporary install folder")
	}
	defer os.RemoveAll(tempInstallFolder)

	archivePath := path
	if !local {
		if archive {
			tempDownloadFolder, err := ioutil.TempDir("", "relique-module-download-*")
			if err != nil {
				return errors.Wrap(err, "cannot create temporary download folder")
			}
			defer os.RemoveAll(tempDownloadFolder)

			// TODO: If archive, download archive
			archiveDownloadDestination := filepath.Clean(fmt.Sprintf("%s/module.tar.gz", tempDownloadFolder))
			archivePath = archiveDownloadDestination
			if err := downloadArchive(archiveDownloadDestination, path); err != nil {
				return errors.Wrap(err, "cannot download module archive")
			}
		} else {
			if err := gitClone(path, tempInstallFolder); err != nil {
				return errors.Wrap(err, "cannot clone git repository")
			}
		}
	}

	if archive {
		if err := extractArchive(archivePath, tempInstallFolder); err != nil {
			return errors.Wrap(err, "cannot extract module archive")
		}
	}

	parsedModule, err := LoadFromFile(fmt.Sprintf("%s/default.toml", tempInstallFolder))
	if err != nil {
		return errors.Wrap(err, "cannot read default.toml for module, please make sure the module to install contains a correctly formatted default.toml file")
	}

	installPath := filepath.Clean(fmt.Sprintf("%s/%s", MODULES_INSTALL_PATH, strings.ToLower(parsedModule.Name)))

	if _, err := os.Lstat(installPath); !os.IsNotExist(err) && !force {
		return fmt.Errorf("module install path already exists but -f/--force not used, skipping module install")
	}

	// Module already exists, removing before reinstall
	if _, err := os.Lstat(installPath); err == nil {
		if err := Remove(parsedModule.Name); err != nil {
			return errors.Wrap(err, "cannot remove existing module version before install")
		}
	}

	if err := installModuleFiles(tempInstallFolder, installPath, skipChown); err != nil {
		return errors.Wrap(err, "cannot install module files to their final destination")
	}

	return nil
}

func Remove(moduleName string) error {
	SetModulePathDefaultValue()

	log.WithFields(log.Fields{
		"install_path": MODULES_INSTALL_PATH,
		"name":         moduleName,
	}).Info("Uninstalling relique module")

	// Control path to remove to avoid making huge mistakes
	if MODULES_INSTALL_PATH == "/" {
		return fmt.Errorf("modules install path has been set to /. Refusing to remove module to avoid accidental disastrous file deletions")
	}

	// Do not trust moduleName directly. Search for installed modules matching with provided name and remove it
	installedModules, err := GetLocallyInstalled()
	if err != nil {
		return errors.Wrap(err, "cannot get locally modules")
	}

	found := false
	var foundModule Module
	for _, mod := range installedModules {
		if mod.Name == moduleName {
			found = true
			foundModule = mod
		}
	}

	if !found {
		return fmt.Errorf("cannot find '%v' in installed modules", moduleName)
	}

	// Add a / in front of path to make sure we have an absolute path
	fullPath := filepath.Clean(fmt.Sprintf("/%s/%s", MODULES_INSTALL_PATH, foundModule.Name))
	if fullPath == "/" {
		return fmt.Errorf("modules installed path to remove has been computed to /. Refusing to remove module to avoid accidental disastrous file deletions")
	}

	if err := os.RemoveAll(fullPath); err != nil {
		return errors.Wrap(err, "cannot remove module")
	}

	return nil
}
