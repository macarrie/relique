use crate::types::config::client::Client;
use crate::types::schedule::Schedule;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

pub mod client;
pub mod config_version;
pub mod error;

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(default)]
pub struct Config {
    pub config_version: Option<Uuid>,
    pub clients: Option<Vec<Client>>,
    pub schedules: Option<Vec<Schedule>>,
    pub bind_addr: Option<String>,
    pub public_address: String,
    pub port: Option<u32>,
    pub ssl_cert: Option<String>,
    pub ssl_key: Option<String>,
    pub strict_ssl_certificate_check: Option<bool>,
    pub clients_cfg_path: Option<String>,
    pub schedules_cfg_path: Option<String>,
    pub backup_storage_path: String,
}

impl Default for Config {
    fn default() -> Self {
        Config {
            config_version: None,
            clients: None,
            schedules: None,
            bind_addr: Some("0.0.0.0".to_string()),
            public_address: "localhost".to_string(),
            port: Some(8433),
            ssl_cert: Some(String::from("/etc/relique/cert.pem")),
            ssl_key: Some(String::from("/etc/relique/key.pem")),
            strict_ssl_certificate_check: Some(false),
            clients_cfg_path: Some(String::from("clients")),
            schedules_cfg_path: Some(String::from("schedules")),
            backup_storage_path: String::from("/opt/relique/"),
        }
    }
}
