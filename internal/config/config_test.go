package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/macarrie/relique/internal/module"
)

func resetConfig(t *testing.T) {
	t.Helper()
	viper.Reset()
	Current = Configuration{}
	Loaded = false
	customConfigFile = false
	customConfigFilePath = ""
}

// writeTestConfig creates a config tree in a temp dir with a main config file, a clients folder and a repositories folder
func writeTestConfig(t *testing.T, mainCfg string) string {
	t.Helper()
	tmp := t.TempDir()

	if err := os.MkdirAll(filepath.Join(tmp, "clients"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "repositories"), 0755); err != nil {
		t.Fatal(err)
	}
	repoCfg := "name = \"local\"\ntype = \"local\"\npath = \"/tmp/relique_repo\"\ndefault = true\n"
	if err := os.WriteFile(filepath.Join(tmp, "repositories", "local.toml"), []byte(repoCfg), 0644); err != nil {
		t.Fatal(err)
	}
	clientCfg := "name = \"cl\"\naddress = \"127.0.0.1\"\n"
	if err := os.WriteFile(filepath.Join(tmp, "clients", "cl.toml"), []byte(clientCfg), 0644); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(tmp, "relique.toml")
	if err := os.WriteFile(cfgPath, []byte(mainCfg), 0644); err != nil {
		t.Fatal(err)
	}

	return cfgPath
}

func TestNew(t *testing.T) {
	resetConfig(t)
	New()

	checks := map[string][2]string{
		"ClientCfgPath":     {Current.ClientCfgPath, CLIENTS_DEFAULT_FOLDER},
		"RepoCfgPath":       {Current.RepoCfgPath, REPOS_DEFAULT_FOLDER},
		"ModuleInstallPath": {Current.ModuleInstallPath, MODULES_DEFAULT_FOLDER},
		"DBPath":            {Current.DBPath, DB_DEFAULT_FOLDER},
		"CatalogPath":       {Current.CatalogPath, CATALOG_DEFAULT_FOLDER},
		"WebUI.BindAddr":    {Current.WebUI.BindAddr, WEBUI_DEFAULT_BIND_ADDR},
		"WebUI.SSLCert":     {Current.WebUI.SSLCert, WEBUI_DEFAULT_SSL_CERT},
		"WebUI.SSLKey":      {Current.WebUI.SSLKey, WEBUI_DEFAULT_SSL_KEY},
	}
	for field, values := range checks {
		if values[0] != values[1] {
			t.Errorf("New() %s = %v, want %v", field, values[0], values[1])
		}
	}
	if Current.WebUI.Port != WEBUI_DEFAULT_PORT {
		t.Errorf("New() WebUI.Port = %v, want %v", Current.WebUI.Port, WEBUI_DEFAULT_PORT)
	}
	if module.MODULES_INSTALL_PATH != MODULES_DEFAULT_FOLDER {
		t.Errorf("New() module.MODULES_INSTALL_PATH = %v, want %v", module.MODULES_INSTALL_PATH, MODULES_DEFAULT_FOLDER)
	}
}

func TestSetDefaultValuesKeepsExisting(t *testing.T) {
	resetConfig(t)
	Current.DBPath = "custom_db"
	Current.WebUI.Port = 1234
	setDefaultValues()

	if Current.DBPath != "custom_db" {
		t.Errorf("setDefaultValues() overwrote DBPath: %v", Current.DBPath)
	}
	if Current.WebUI.Port != 1234 {
		t.Errorf("setDefaultValues() overwrote WebUI.Port: %v", Current.WebUI.Port)
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		mainCfg     string
		useFile     bool
		wantErr     bool
		wantClients int
		wantRepos   int
	}{
		{
			name:        "valid_config",
			mainCfg:     "module_install_path = \"/tmp/relique_modules\"\n",
			useFile:     true,
			wantErr:     false,
			wantClients: 1,
			wantRepos:   1,
		},
		{
			name:    "parse_error",
			mainCfg: "this is not = = toml",
			useFile: true,
			wantErr: true,
		},
		{
			name:    "missing_clients_folder",
			mainCfg: "client_cfg_path = \"does_not_exist\"\n",
			useFile: true,
			wantErr: true,
		},
		{
			name:    "missing_repos_folder",
			mainCfg: "repo_cfg_path = \"does_not_exist\"\n",
			useFile: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetConfig(t)
			defer resetConfig(t)

			cfgPath := writeTestConfig(t, tt.mainCfg)
			if tt.useFile {
				UseFile(cfgPath)
			}

			err := Load("relique")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if Loaded {
					t.Errorf("Load() set Loaded = true on error")
				}
				return
			}

			if !Loaded {
				t.Errorf("Load() did not set Loaded")
			}
			if len(Current.Clients) != tt.wantClients {
				t.Errorf("Load() loaded %d clients, want %d", len(Current.Clients), tt.wantClients)
			}
			if len(Current.Repositories) != tt.wantRepos {
				t.Errorf("Load() loaded %d repositories, want %d", len(Current.Repositories), tt.wantRepos)
			}
			if Current.CurrentFile != cfgPath {
				t.Errorf("Load() CurrentFile = %v, want %v", Current.CurrentFile, cfgPath)
			}
			if module.MODULES_INSTALL_PATH != "/tmp/relique_modules" {
				t.Errorf("Load() module.MODULES_INSTALL_PATH = %v, want /tmp/relique_modules", module.MODULES_INSTALL_PATH)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	resetConfig(t)
	defer resetConfig(t)
	New()

	tmp := t.TempDir()
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "valid_path", path: filepath.Join(tmp, "out.toml"), wantErr: false},
		{name: "invalid_path", path: filepath.Join(tmp, "missing", "out.toml"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Write(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			content, readErr := os.ReadFile(tt.path)
			if readErr != nil {
				t.Fatalf("cannot read written config: %v", readErr)
			}
			if !strings.Contains(string(content), `catalog_path = "catalog"`) {
				t.Errorf("Write() unexpected contents: %s", content)
			}
		})
	}
}

func TestGetCfgPaths(t *testing.T) {
	resetConfig(t)
	defer resetConfig(t)

	tmp := t.TempDir()
	viper.SetConfigFile(filepath.Join(tmp, "relique.toml"))

	// Empty values must fall back to defaults
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{name: "clients", got: GetClientsCfgPath, want: filepath.Join(tmp, CLIENTS_DEFAULT_FOLDER)},
		{name: "repos", got: GetReposCfgPath, want: filepath.Join(tmp, REPOS_DEFAULT_FOLDER)},
		{name: "catalog", got: GetCatalogCfgPath, want: filepath.Join(tmp, CATALOG_DEFAULT_FOLDER)},
		{name: "db", got: GetDBPath, want: filepath.Join(tmp, DB_DEFAULT_FOLDER)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filepath.Clean(tt.got()); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

	Current.DBPath = "custom_db"
	if got := filepath.Clean(GetDBPath()); got != filepath.Join(tmp, "custom_db") {
		t.Errorf("GetDBPath() with custom value = %v, want %v", got, filepath.Join(tmp, "custom_db"))
	}
}

func TestCheck(t *testing.T) {
	if err := Check(); err != nil {
		t.Errorf("Check() error = %v, want nil", err)
	}
}
