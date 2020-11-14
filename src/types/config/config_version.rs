use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
pub struct ConfigVersion {
    pub version: Option<Uuid>,
}

impl Default for ConfigVersion {
    fn default() -> Self {
        ConfigVersion { version: None }
    }
}
