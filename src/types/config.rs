use serde::{Deserialize, Serialize};
use std::fmt;
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(default)]
pub struct Client {
    pub name: String,
    pub address: String,
    pub port: Option<u32>,
    pub modules: Vec<String>,
    pub config_version: Option<Uuid>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(default)]
pub struct Config {
    pub config_version: Option<Uuid>,
    pub clients: Option<Vec<Client>>,
    pub bind_addr: Option<String>,
    pub port: Option<u32>,
    pub ssl_cert: Option<String>,
    pub ssl_key: Option<String>,
    pub strict_ssl_certificate_check: Option<bool>,
    pub clients_cfg_path: Option<String>,
}

#[derive(Serialize, Deserialize)]
pub struct ConfigVersion {
    pub version: Option<Uuid>,
}

impl Default for Client {
    fn default() -> Self {
        Client {
            config_version: None,
            name: "".to_string(),
            address: "".to_string(),
            port: Some(8433),
            modules: vec![],
        }
    }
}

impl Default for Config {
    fn default() -> Self {
        Config {
            config_version: None,
            clients: None,
            bind_addr: Some("0.0.0.0".to_string()),
            port: Some(8433),
            ssl_cert: Some(String::from("/etc/relique/cert.pem")),
            ssl_key: Some(String::from("/etc/relique/key.pem")),
            strict_ssl_certificate_check: Some(false),
            clients_cfg_path: Some(String::from("clients")),
        }
    }
}

impl Default for ConfigVersion {
    fn default() -> Self {
        ConfigVersion { version: None }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum ErrorLevel {
    Warning,
    Critical,
}

#[derive(Debug, Clone)]
pub struct Error {
    pub key: String,
    pub level: ErrorLevel,
    pub desc: String,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_fmt(format_args!(
            "[{:?}] {} (key: '{}')",
            self.level, self.desc, self.key
        ))
    }
}
