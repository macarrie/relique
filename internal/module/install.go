package module

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-multierror"
)

func GetLocallyInstalled(modulesInstallPath string) ([]Module, error) {
	slog.With(
		slog.String("path", modulesInstallPath),
	).Debug("Using modules install folder")

	items, err := os.ReadDir(modulesInstallPath)
	if err != nil {
		return []Module{}, fmt.Errorf("cannot list installed modules from filesystem: %w", err)
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
			slog.With(
				slog.Any("error", err),
			).Error("Cannot read module default configuration")
			errorList = multierror.Append(errorList, err)
			continue
		}
		if err := mod.GetAvailableVariants(); err != nil {
			mod.GetLog().With(
				slog.Any("error", err),
			).Error("Cannot get module available variants")
			errorList = multierror.Append(errorList, err)
		}
		modulesList = append(modulesList, mod)
	}
	return modulesList, errorList.ErrorOrNil()
}

func IsInstalled(modulesInstallPath string, moduleName string) (bool, error) {
	installedModules, err := GetLocallyInstalled(modulesInstallPath)
	if err != nil {
		return false, fmt.Errorf("cannot get locally installed modules: %w", err)
	}

	for _, mod := range installedModules {
		if mod.Name == moduleName {
			return true, nil
		}
	}

	return false, nil
}

func gitClone(source string, dest string) error {
	slog.With(
		slog.String("source_path", source),
		slog.String("destination", dest),
	).Debug("Cloning git repository")

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
		URL: fmt.Sprintf("%s%s%s", protocolPrefix, source, gitSuffix),
		//Progress: os.Stdout,
	}); err != nil {
		return fmt.Errorf("cannot clone git repository: %w", err)
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
		return fmt.Errorf("cannot perform download request: %w", err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer out.Close()

	// Write the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("cannot write downloaded archive to destination file: %w", err)
	}

	return nil
}

func extractArchive(source string, dest string) error {
	slog.With(
		slog.String("source_path", source),
		slog.String("destination", dest),
	).Debug("Extracting module archive")

	if _, err := os.Lstat(source); os.IsNotExist(err) {
		return fmt.Errorf("cannot find archive file to extract")
	}

	extractCmd := exec.Command("tar", "-xvf", source, "-C", dest)

	var out bytes.Buffer
	extractCmd.Stdout = &out
	extractCmd.Stderr = &out

	err := extractCmd.Run()
	if err != nil {
		return fmt.Errorf("cannot extract module archive: %w, %s", err, out.String())
	}

	// TODO: Chmod/chown files

	return nil
}

func installModuleFiles(source string, dest string, skipChown bool) error {
	slog.With(
		slog.String("source_path", source),
		slog.String("destination", dest),
	).Debug("Installing module files")

	copyCmd := exec.Command("cp", "-r", source, dest)

	var out bytes.Buffer
	copyCmd.Stdout = &out
	copyCmd.Stderr = &out

	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("cannot copy module files: %s, %w", out.String(), err)
	}

	chmodCmd := exec.Command("chmod", "-R", "755", dest)

	var out2 bytes.Buffer
	chmodCmd.Stdout = &out2
	chmodCmd.Stderr = &out2

	if err := chmodCmd.Run(); err != nil {
		return fmt.Errorf("cannot chmod installed module files: %s, %w", out2.String(), err)
	}

	if !skipChown {
		chownCmd := exec.Command("chown", "-R", ":relique", dest)

		var out3 bytes.Buffer
		chownCmd.Stdout = &out3
		chownCmd.Stderr = &out3

		if err := chownCmd.Run(); err != nil {
			return fmt.Errorf("cannot chown installed module files: %s, %w", out3.String(), err)
		}
	}

	return nil
}

func Install(modulesInstallPath string, path string, local bool, archive bool, force bool, skipChown bool) error {
	slog.With(
		slog.String("install_path", modulesInstallPath),
		slog.Bool("local", local),
		slog.Bool("archive", archive),
		slog.String("source_path", path),
	).Info("Installing relique module")

	if _, err := os.Lstat(modulesInstallPath); os.IsNotExist(err) {
		return fmt.Errorf("module install path does not exist. Please check that relique is correctly installed and that provided module install path is correct")
	}

	tempInstallFolder, err := os.MkdirTemp("", "relique-module-install-*")
	if err != nil {
		return fmt.Errorf("cannot create temporary install folder: %w", err)
	}
	defer os.RemoveAll(tempInstallFolder)

	archivePath := path
	if !local {
		if archive {
			tempDownloadFolder, err := os.MkdirTemp("", "relique-module-download-*")
			if err != nil {
				return fmt.Errorf("cannot create temporary download folder: %w", err)
			}
			defer os.RemoveAll(tempDownloadFolder)

			// TODO: If archive, download archive
			archiveDownloadDestination := filepath.Clean(fmt.Sprintf("%s/module.tar.gz", tempDownloadFolder))
			archivePath = archiveDownloadDestination
			if err := downloadArchive(archiveDownloadDestination, path); err != nil {
				return fmt.Errorf("cannot download module archive: %w", err)
			}
		} else {
			if err := gitClone(path, tempInstallFolder); err != nil {
				return fmt.Errorf("cannot clone git repository: %w", err)
			}
		}
	}

	if archive {
		if err := extractArchive(archivePath, tempInstallFolder); err != nil {
			return fmt.Errorf("cannot extract module archive: %w", err)
		}
	}

	parsedModule, err := LoadFromFile(fmt.Sprintf("%s/default.toml", tempInstallFolder))
	if err != nil {
		return fmt.Errorf("cannot read default.toml for module, please make sure the module to install contains a correctly formatted default.toml file: %w", err)
	}

	installPath := filepath.Clean(fmt.Sprintf("%s/%s", modulesInstallPath, strings.ToLower(parsedModule.Name)))

	if _, err := os.Lstat(installPath); !os.IsNotExist(err) && !force {
		return fmt.Errorf("module install path already exists but -f/--force not used, skipping module install")
	}

	// Module already exists, removing before reinstall
	if _, err := os.Lstat(installPath); err == nil {
		if err := Remove(modulesInstallPath, parsedModule.Name); err != nil {
			return fmt.Errorf("cannot remove existing module version before install: %w", err)
		}
	}

	if err := installModuleFiles(tempInstallFolder, installPath, skipChown); err != nil {
		return fmt.Errorf("cannot install module files to their final destination: %w", err)
	}

	return nil
}

func Remove(modulesInstallPath string, moduleName string) error {
	slog.With(
		slog.String("install_path", modulesInstallPath),
		slog.String("name", moduleName),
	).Info("Uninstalling relique module")

	// Control path to remove to avoid making huge mistakes
	if modulesInstallPath == "/" {
		return fmt.Errorf("modules install path has been set to /. Refusing to remove module to avoid accidental disastrous file deletions")
	}

	// Do not trust moduleName directly. Search for installed modules matching with provided name and remove it
	installedModules, err := GetLocallyInstalled(modulesInstallPath)
	if err != nil {
		return fmt.Errorf("cannot get locally modules: %w", err)
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
	fullPath := filepath.Clean(fmt.Sprintf("/%s/%s", modulesInstallPath, foundModule.Name))
	if fullPath == "/" {
		return fmt.Errorf("modules installed path to remove has been computed to /. Refusing to remove module to avoid accidental disastrous file deletions")
	}

	if err := os.RemoveAll(fullPath); err != nil {
		return fmt.Errorf("cannot remove module: %w", err)
	}

	return nil
}
