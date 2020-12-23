package module

import (
	"database/sql"
	"fmt"

	"github.com/macarrie/relique/internal/db"
	"github.com/pkg/errors"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/spf13/viper"
)

type Module struct {
	ID               int64
	ModuleType       string `mapstructure:"module_type"`
	Name             string `mapstructure:"name"`
	BackupTypeString string `mapstructure:"backup_type"`
	// TODO:Load BackupType const directly from "diff" and "full" strings
	// TODO: Compute BackupType when loading configuration
	BackupType backup_type.BackupType
	// TODO:Load schedule struct
	Schedules         []string
	BackupPaths       []string `mapstructure:"backup_paths"`
	PreBackupScript   string   `mapstructure:"pre_backup_script"`
	PostBackupScript  string   `mapstructure:"post_backup_script"`
	PreRestoreScript  string   `mapstructure:"pre_restore_script"`
	PostRestoreScript string   `mapstructure:"post_restore_script"`
}

func (m *Module) String() string {
	return m.Name
}

func (m *Module) ComputeBackupTypeFromString() {
	var t backup_type.BackupType
	t, err := backup_type.FromString(m.BackupTypeString)
	if err != nil {
		log.WithFields(log.Fields{
			"field": m.BackupTypeString,
			"err":   err,
		}).Error("Cannot parse backup type from configuration file. Setting backup type to full")
		t.Type = backup_type.Full
	}

	m.BackupType = t
}

func (m *Module) LoadDefaultConfiguration() error {
	defaults, err := LoadFromFile(fmt.Sprintf("/var/lib/relique/modules/%s/default.toml", m.ModuleType))
	if err != nil {
		return err
	}

	if len(m.BackupPaths) == 0 {
		m.BackupPaths = defaults.BackupPaths
	}

	if m.PreBackupScript == "" {
		m.PreBackupScript = defaults.PreBackupScript
	}

	if m.PostBackupScript == "" {
		m.PostBackupScript = defaults.PostBackupScript
	}

	if m.PreRestoreScript == "" {
		m.PreRestoreScript = defaults.PreRestoreScript
	}

	if m.PostRestoreScript == "" {
		m.PostRestoreScript = defaults.PostRestoreScript
	}

	return nil
}

func LoadFromFile(file string) (Module, error) {
	log.WithFields(log.Fields{
		"path": file,
	}).Debug("Loading module configuration parameters from file")

	modViper := viper.New()
	modViper.SetConfigType("toml")
	modViper.SetConfigFile(file)

	if err := modViper.ReadInConfig(); err != nil {
		return Module{}, err
	}

	var module Module
	if err := modViper.Unmarshal(&module); err != nil {
		return Module{}, err
	}

	return module, nil
}

func GetID(name string) (int64, error) {
	row := db.Write().QueryRow("SELECT id FROM modules WHERE name = $1", name)
	defer db.Unlock()

	var id int64
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, errors.Wrap(err, "cannot search retrieve module ID in db")
	}

	return id, nil
}

func GetByID(id int64) (Module, error) {
	log.WithFields(log.Fields{
		"id": id,
	}).Trace("Looking for module in database")

	request := `SELECT  
		id,
		module_type, 
		name, 
		backup_type, 
		pre_backup_script, 
		post_backup_script, 
		pre_restore_script, 
		post_restore_script 
	FROM modules
	WHERE id = $1`
	row := db.Read().QueryRow(request, id)
	defer db.RUnlock()

	var mod Module
	err := row.Scan(&mod.ID,
		&mod.ModuleType,
		&mod.Name,
		&mod.BackupType.Type,
		&mod.PreBackupScript,
		&mod.PostBackupScript,
		&mod.PreRestoreScript,
		&mod.PostRestoreScript,
	)
	if err == sql.ErrNoRows {
		return Module{}, nil
	} else if err != nil {
		return Module{}, errors.Wrap(err, "cannot retrieve module from db")
	}

	return mod, nil
}

func (m *Module) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"id":          m.ID,
		"name":        m.Name,
		"type":        m.ModuleType,
		"backup_type": m.BackupType.String(),
	})
}

func (m *Module) Save() (int64, error) {
	id, err := GetID(m.Name)
	if err != nil {
		return 0, errors.Wrap(err, "cannot search for possibly existing module ID")
	}

	if id != 0 {
		m.ID = id
		return m.Update()
	}

	m.GetLog().Debug("Saving module into database")

	sql := `INSERT INTO modules ( module_type, name, backup_type, pre_backup_script, post_backup_script, pre_restore_script, post_restore_script )
VALUES ( 
	$1, 
	$2, 
	$3, 
	$4, 
	$5, 
	$6, 
	$7
)`

	result, err := db.Write().Exec(
		sql,
		m.ModuleType,
		m.Name,
		m.BackupType.Type,
		m.PreBackupScript,
		m.PostBackupScript,
		m.PreRestoreScript,
		m.PostRestoreScript,
	)
	defer db.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "cannot save module into db")
	}

	m.ID, err = result.LastInsertId()
	if m.ID == 0 {
		return 0, errors.Wrap(err, "cannot get last insert id for job")
	}

	return m.ID, nil
}

func (m *Module) Update() (int64, error) {
	m.GetLog().Debug("Updating module details into database")

	if m.ID == 0 {
		return 0, fmt.Errorf("cannot update module with ID 0")
	}

	sql := `UPDATE modules
	SET module_type = $1,
		name = $2,
		backup_type = $3,
		pre_backup_script = $4,
		post_backup_script = $5,
		pre_restore_script = $6,
		post_restore_script = $7 
	WHERE id = $8`

	_, err := db.Write().Exec(
		sql,
		m.ModuleType,
		m.Name,
		m.BackupType.Type,
		m.PreBackupScript,
		m.PostBackupScript,
		m.PreRestoreScript,
		m.PostRestoreScript,
		m.ID,
	)
	defer db.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "cannot update module into db")
	}

	return m.ID, nil
}
