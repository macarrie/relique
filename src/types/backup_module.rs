use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BackupModule {
    pub module_type: String,
    pub name: String,
    pub backup_type: Option<String>, // Transform to enum
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
