use crate::types::backup_type::BackupType;
use crate::types::db;
use anyhow::Result;
use log::*;
use rusqlite::named_params;
use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BackupModule {
    pub module_type: String,
    pub name: String,
    pub backup_type: Option<BackupType>,
    pub schedules: Option<Vec<String>>,
    pub backup_paths: Option<Vec<String>>,
    pub pre_backup_script: Option<String>,
    pub post_backup_script: Option<String>,
    pub pre_restore_script: Option<String>,
    pub post_restore_script: Option<String>,
}

impl fmt::Display for BackupModule {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{} (type {})", self.name, self.module_type)
    }
}

impl BackupModule {
    pub fn find_by_id(id: i64, pool: db::Pool) -> Result<BackupModule> {
        let conn = pool.get()?;
        let mut stmt = conn.prepare(
            "SELECT ( module_type, name, backup_type, pre_backup_script, post_backup_script, pre_restore_script, post_restore_script ) \
                 FROM modules \
                 WHERE id = :id"
        )?;
        let modules: Vec<BackupModule> = stmt
            .query_map_named(named_params![":id": id], |row| {
                Ok(BackupModule {
                    module_type: row.get(1)?,
                    name: row.get(2)?,
                    backup_type: row.get(3)?,
                    pre_backup_script: row.get(4)?,
                    post_backup_script: row.get(5)?,
                    pre_restore_script: row.get(6)?,
                    post_restore_script: row.get(7)?,
                    // TODO: Load schedules
                    schedules: None,
                    backup_paths: None,
                })
            })?
            .map(Result::unwrap)
            .collect();

        if modules.is_empty() {
            anyhow::bail!(
                "Cannot find backup module with unique ID '{}' in database",
                id
            );
        }

        if modules.len() > 1 {
            anyhow::bail!(
                "Found {len} backup modules with unique ID '{id}' in database",
                len = modules.len(),
                id = id,
            );
        }

        Ok(modules[0].clone())
    }

    pub fn get_db_id(name: &String, pool: db::Pool) -> Result<i64> {
        let conn = pool.get()?;
        let mut stmt = conn.prepare("SELECT (id) FROM modules WHERE name = :name")?;
        let ids: Vec<i64> = stmt
            .query_map_named(named_params![":name": name], |row| Ok(row.get(0)?))?
            .map(Result::unwrap)
            .collect();

        if ids.is_empty() {
            anyhow::bail!("Cannot find backup module with name '{}' in database", name);
        }

        if ids.len() > 1 {
            anyhow::bail!("Found {len} backup modules with name '{}' in database",);
        }

        Ok(ids[0])
    }

    pub fn save(&self, pool: db::Pool) -> Result<i64> {
        debug!("Saving backup module '{}' into database", self.name);
        let params = named_params! {
            ":module_type": self.module_type,
            ":name": self.name,
            ":backup_type": self.backup_type,
            ":pre_backup_script": self.pre_backup_script,
            ":post_backup_script": self.post_backup_script,
            ":pre_restore_script": self.pre_restore_script,
            ":post_restore_script": self.post_restore_script,
        };

        let conn = pool.get()?;

        if let Ok(id) = BackupModule::get_db_id(&self.name, pool.clone()) {
            let update_params = named_params! {
                ":id": id,
                ":module_type": self.module_type,
                ":name": self.name,
                ":backup_type": self.backup_type,
                ":pre_backup_script": self.pre_backup_script,
                ":post_backup_script": self.post_backup_script,
                ":pre_restore_script": self.pre_restore_script,
                ":post_restore_script": self.post_restore_script,
            };
            conn.execute_named(
                "UPDATE modules \
                SET module_type = :module_type, \
                    name = :name, \
                    backup_type = :backup_type, \
                    pre_backup_script = :pre_backup_script, \
                    post_backup_script = :post_backup_script, \
                    pre_restore_script = :pre_restore_script, \
                    post_restore_script = :post_restore_script  \
                WHERE id = :id",
                update_params,
            )?;
            return Ok(id);
        }

        conn.execute_named(
            "INSERT INTO modules ( module_type, name, backup_type, pre_backup_script, post_backup_script, pre_restore_script, post_restore_script ) \
            VALUES ( :module_type, :name, :backup_type, :pre_backup_script, :post_backup_script, :pre_restore_script, :post_restore_script )",
            params,
        )?;

        Ok(conn.last_insert_rowid())
    }
}
