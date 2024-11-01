package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/dustin/go-humanize"
	"github.com/manifoldco/promptui"
	"github.com/pelletier/go-toml"

	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/repo"
)

func GetStoragePath(r repo.Repository, repoName string, uuid string) (string, error) {
	var repository repo.Repository

	if r == nil {
		if repoName == "" {
			return "", fmt.Errorf("cannot get storage path because used repository is unknown")
		}

		repoFromConfig, err := repo.GetByName(config.Current.Repositories, repoName)
		if err != nil {
			return "", fmt.Errorf("cannot get repository from configuration: %w", err)
		}

		repository = repoFromConfig
	} else {
		repository = r
	}

	if repository.GetType() == "local" {
		localRepo := repository.(*repo.RepositoryLocal)
		return filepath.Clean(fmt.Sprintf("%s/%s/", localRepo.Path, uuid)), nil
	}

	return "", fmt.Errorf("repo type not implemented")
}

func GetCatalogPath(uuid string) (string) {
	return filepath.Clean(fmt.Sprintf("%s/%s", config.GetCatalogCfgPath(), uuid))
}

func FormatDatetime(t time.Time) string {
	if t.IsZero() {
		return "---"
	}

	return t.Format("2006/01/02 15:04:05")
}

func FormatDuration(d time.Duration) string {
	return d.String()
}

func SerializeToFile[S interface{}](data S, destinationFile string) error {
	tomlContents, serializeErr := toml.Marshal(data)
	if serializeErr != nil {
		return fmt.Errorf("cannot serialize data to toml: %w", serializeErr)
	}
	if err := os.WriteFile(destinationFile, tomlContents, 0644); err != nil {
		return fmt.Errorf("cannot export data to file: %w", err)
	}

	return nil
}

func RenderTemplateToMarkdown(templateName string, tpl string, data interface{}) (string, error) {
	var tplRender bytes.Buffer
	t := template.Must(
		template.New(templateName).Funcs(template.FuncMap{
			"join":      strings.Join,
			"file_size": humanize.Bytes,
			"datetime":  FormatDatetime,
		}).Parse(tpl),
	)
	if err := t.Execute(&tplRender, data); err != nil {
		return "", fmt.Errorf("cannot render template: %w", err)
	}

	out, err := glamour.Render(tplRender.String(), "auto")
	if err != nil {
		return "", fmt.Errorf("cannot render markdown: %w", err)
	}

	return out, err
}

func GenerateCustomRestorePaths(raw []string, modBackupPaths []string) map[string]string {
	restorePaths := make(map[string]string)
	for _, p := range raw {
		source, dest, _ := strings.Cut(p, ":")
		if source != "" {
			// --path toto:tata -> Restore backup path 'toto' to 'tata' on client
			restorePaths[source] = dest
			// --path test -> Select 'test' backup path and restore it to its own path
			if dest == "" {
				restorePaths[source] = source
			}
		} else {
			// :dest -> Restore everything to a different location
			restorePaths = make(map[string]string)
			for _, bp := range modBackupPaths {
				restorePaths[bp] = filepath.Clean(fmt.Sprintf("%s/%s", dest, bp))
			}
			break
		}
	}

	return restorePaths
}

func Confirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	return err == nil
}
