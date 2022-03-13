package common

import (
	"fmt"

	"github.com/macarrie/relique/internal/types/module"

	"github.com/macarrie/relique/internal/types/displayable"
)

type ConfigDisplay struct {
	Clients                   []displayable.Displayable `json:"clients"`
	Schedules                 []displayable.Displayable `json:"schedules"`
	BindAddr                  string                    `json:"bind_addr"`
	PublicAddress             string                    `json:"public_address"`
	Port                      uint32                    `json:"port"`
	SSLCert                   string                    `json:"ssl_cert"`
	SSLKey                    string                    `json:"ssl_key"`
	StrictSSLCertificateCheck bool                      `json:"strict_ssl_certificate_check"`
	ClientCfgPath             string                    `json:"client_cfg_path"`
	SchedulesCfgPath          string                    `json:"schedules_cfg_path"`
	BackupStoragePath         string                    `json:"backup_storage_path"`
	ModuleInstallPath         string                    `json:"module_install_path"`
	RetentionPath             string                    `json:"retention_path"`
}

func (c Configuration) Display() displayable.Struct {
	var clientDisplay []displayable.Displayable
	for _, mod := range c.Clients {
		clientDisplay = append(clientDisplay, mod)
	}

	var scheduleDisplay []displayable.Displayable
	for _, mod := range c.Schedules {
		scheduleDisplay = append(scheduleDisplay, mod)
	}

	var d displayable.Struct = ConfigDisplay{
		Clients:                   clientDisplay,
		Schedules:                 scheduleDisplay,
		BindAddr:                  c.BindAddr,
		PublicAddress:             c.PublicAddress,
		Port:                      c.Port,
		SSLCert:                   c.SSLCert,
		SSLKey:                    c.SSLKey,
		StrictSSLCertificateCheck: c.StrictSSLCertificateCheck,
		ClientCfgPath:             c.ClientCfgPath,
		SchedulesCfgPath:          c.SchedulesCfgPath,
		BackupStoragePath:         c.BackupStoragePath,
		ModuleInstallPath:         c.ModuleInstallPath,
		RetentionPath:             c.RetentionPath,
	}

	return d
}

func (d ConfigDisplay) Summary() string {
	installedModules, _ := module.GetLocallyInstalled()
	return fmt.Sprintf("%d clients, %d schedules, %d installed modules", len(d.Clients), len(d.Schedules), len(installedModules))
}

func (d ConfigDisplay) Details() string {
	var clientsSummary string
	var schedulesSummary string
	var installedModulesSummary string

	for _, cl := range d.Clients {
		clientsSummary = fmt.Sprintf("%s\t- %s\n", clientsSummary, cl.Display().Summary())
	}

	for _, sch := range d.Schedules {
		schedulesSummary = fmt.Sprintf("%s\t- %s\n", schedulesSummary, sch.Display().Summary())
	}

	installedModules, _ := module.GetLocallyInstalled()
	for _, mod := range installedModules {
		installedModulesSummary = fmt.Sprintf("%s\t- %s\n", installedModulesSummary, mod.Display().Summary())
	}
	return fmt.Sprintf(
		`
GLOBAL CONFIGURATION
--------------------

Public address: 	%s
Bind address: 		%s
Port: 			%d

SSL certificate: 	%s
SSL key: 		%s
Strict SSL check: 	%v

Client config path: 	%s
Schedules config path: 	%s
Backup storage path: 	%s
Modules install path: 	%s
Retention file path: 	%s

CLIENTS
-------
%s

SCHEDULES
---------
%s

INSTALLED MODULES TYPES
-----------------------
%s`,
		d.PublicAddress,
		d.BindAddr,
		d.Port,
		d.SSLCert,
		d.SSLKey,
		d.StrictSSLCertificateCheck,
		d.ClientCfgPath,
		d.SchedulesCfgPath,
		d.BackupStoragePath,
		d.ModuleInstallPath,
		d.RetentionPath,
		clientsSummary,
		schedulesSummary,
		installedModulesSummary,
	)
}

func (d ConfigDisplay) TableHeaders() []string {
	// No table display for config
	return []string{}
}

func (d ConfigDisplay) TableRow() []string {
	// No table display for config
	return []string{}
}
