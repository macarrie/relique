use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BackupModuleDef {
    pub module_type: String,
    pub name: String,
    pub schedules: Option<Vec<String>>,
    pub backup_paths: Option<Vec<String>>,
    pub pre_backup_script: Option<String>,
    pub post_backup_script: Option<String>,
    pub pre_restore_script: Option<String>,
    pub post_restore_script: Option<String>,
}

impl fmt::Display for BackupModuleDef {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{} (type {})", self.name, self.module_type)
    }
}

pub trait BackupModule: Sized {
    fn new() -> Result<Self>;

    fn backup(self) -> Result<()>;

    fn restore(self) -> Result<()>;
}
