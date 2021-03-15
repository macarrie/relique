package module

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/hashicorp/go-multierror"

	"github.com/macarrie/relique/internal/types/custom_errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/macarrie/relique/internal/db"
	"github.com/pkg/errors"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/pelletier/go-toml"
)

var MODULES_INSTALL_PATH = "/var/lib/relique/modules"

type Module struct {
	ID                int64
	ModuleType        string                 `json:"module_type" toml:"module_type"`
	Name              string                 `json:"name" toml:"name"`
	BackupType        backup_type.BackupType `json:"backup_type" toml:"backup_type"`
	Schedules         []schedule.Schedule    `json:"schedules" toml:"-"`
	ScheduleNames     []string               `json:"-" toml:"schedules"`
	BackupPaths       []string               `json:"backup_paths" toml:"backup_paths"`
	PreBackupScript   string                 `json:"pre_backup_script" toml:"pre_backup_script"`
	PostBackupScript  string                 `json:"post_backup_script" toml:"post_backup_script"`
	PreRestoreScript  string                 `json:"pre_restore_script" toml:"pre_restore_script"`
	PostRestoreScript string                 `json:"post_restore_script" toml:"post_restore_script"`
}

func (m *Module) String() string {
	return m.Name
}

func (m *Module) LoadDefaultConfiguration() error {
	defaults, err := LoadFromFile(fmt.Sprintf("%s/%s/default.toml", MODULES_INSTALL_PATH, m.ModuleType))
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

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return Module{}, errors.Wrap(err, "cannot open file")
	}

	content, _ := ioutil.ReadAll(f)

	var module Module
	if err := toml.Unmarshal(content, &module); err != nil {
		return Module{}, errors.Wrap(err, "cannot parse toml file")
	}

	// TODO: Load schedules by name

	if err := module.Valid(); err != nil {
		return Module{}, errors.Wrap(err, "invalid module loaded from file")
	}

	return module, nil
}

func GetID(name string, tx *sql.Tx) (int64, error) {
	request := sq.Select(
		"id",
	).From(
		"modules",
	).Where(
		"name = ?",
		name,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	var row *sql.Row
	if tx == nil {
		row = db.Read().QueryRow(query, args...)
		defer db.RUnlock()
	} else {
		row = tx.QueryRow(query, args...)
	}

	var id int64
	if err := row.Scan(&id); err == sql.ErrNoRows {
		return 0, &custom_errors.DBNotFoundError{
			ID:       0,
			ItemType: "module",
		}
	} else if err != nil {
		return 0, errors.Wrap(err, "cannot search retrieve module ID in db")
	}

	return id, nil
}

func GetByID(id int64) (Module, error) {
	log.WithFields(log.Fields{
		"id": id,
	}).Trace("Looking for module in database")

	request := sq.Select(
		"id",
		"module_type",
		"name",
		"backup_type",
		"backup_paths",
		"pre_backup_script",
		"post_backup_script",
		"pre_restore_script",
		"post_restore_script",
	).From(
		"modules",
	).Where(
		"id = ?",
		id,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return Module{}, errors.Wrap(err, "cannot build sql query")
	}

	row := db.Read().QueryRow(query, args...)
	defer db.RUnlock()

	var rawBackupPaths string
	var mod Module
	if err := row.Scan(&mod.ID,
		&mod.ModuleType,
		&mod.Name,
		&mod.BackupType.Type,
		&rawBackupPaths,
		&mod.PreBackupScript,
		&mod.PostBackupScript,
		&mod.PreRestoreScript,
		&mod.PostRestoreScript,
	); err == sql.ErrNoRows {
		return Module{}, &custom_errors.DBNotFoundError{
			ID:       id,
			ItemType: "module",
		}
	} else if err != nil {
		return Module{}, errors.Wrap(err, "cannot retrieve module from db")
	}

	mod.BackupPaths = strings.Split(rawBackupPaths, ":")
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

func (m *Module) Save(tx *sql.Tx) (int64, error) {
	id, err := GetID(m.Name, tx)
	if err != nil && !custom_errors.IsDBNotFoundError(err) {
		return 0, errors.Wrap(err, "cannot search for possibly existing module ID")
	}

	if id != 0 {
		m.ID = id
		return m.Update(tx)
	}

	m.GetLog().Debug("Saving module into database")

	request := sq.Insert("modules").Columns(
		"module_type",
		"name",
		"backup_type",
		"backup_paths",
		"pre_backup_script",
		"post_backup_script",
		"pre_restore_script",
		"post_restore_script",
	).Values(
		db.GetNullString(m.ModuleType),
		db.GetNullString(m.Name),
		db.GetNullInt32(uint32(m.BackupType.Type)),
		strings.Join(m.BackupPaths, ":"),
		m.PreBackupScript,
		m.PostBackupScript,
		m.PreRestoreScript,
		m.PostRestoreScript,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	var result sql.Result
	if tx == nil {
		result, err = db.Write().Exec(query, args...)
		defer db.Unlock()
	} else {
		result, err = tx.Exec(query, args...)
	}
	if err != nil {
		return 0, errors.Wrap(err, "cannot save module into db")
	}
	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 || err != nil {
		return 0, errors.Wrap(err, "no rows affected when saving item")
	}

	m.ID, err = result.LastInsertId()
	if m.ID == 0 || err != nil {
		return 0, errors.Wrap(err, "cannot get last insert ID")
	}

	return m.ID, nil
}

func (m *Module) Update(tx *sql.Tx) (int64, error) {
	m.GetLog().Debug("Updating module details into database")

	if m.ID == 0 {
		return 0, fmt.Errorf("cannot update module with ID 0")
	}

	request := sq.Update("modules").SetMap(sq.Eq{
		"module_type":         db.GetNullString(m.ModuleType),
		"name":                db.GetNullString(m.Name),
		"backup_type":         db.GetNullInt32(uint32(m.BackupType.Type)),
		"backup_paths":        strings.Join(m.BackupPaths, ":"),
		"pre_backup_script":   m.PreBackupScript,
		"post_backup_script":  m.PostBackupScript,
		"pre_restore_script":  m.PreRestoreScript,
		"post_restore_script": m.PostRestoreScript,
	}).Where(
		" id = ?",
		m.ID,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	var result sql.Result
	if tx == nil {
		result, err = db.Write().Exec(query, args...)
		defer db.Unlock()
	} else {
		result, err = tx.Exec(query, args...)
	}
	if err != nil {
		return 0, errors.Wrap(err, "cannot update module into db")
	}
	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 || err != nil {
		return 0, errors.Wrap(err, "no rows affected when updating item")
	}

	return m.ID, nil
}

func (m *Module) Valid() error {
	var objErrors *multierror.Error
	if m.ModuleType == "" {
		objErrors = multierror.Append(objErrors, fmt.Errorf("empty module type"))
	}
	if m.Name == "" {
		objErrors = multierror.Append(objErrors, fmt.Errorf("empty module name"))
	}
	if m.BackupType.Type == backup_type.Unknown {
		objErrors = multierror.Append(objErrors, fmt.Errorf("unknown backup type"))
	}

	return objErrors.ErrorOrNil()
}
