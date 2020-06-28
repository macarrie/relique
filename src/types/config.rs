use serde::{Serialize, Deserialize};
use std::fmt;

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(default)]
pub struct Client {
    pub name :String,
    pub address :String,
    pub port :Option<u32>,
    pub modules :Vec<String>,
}


#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(default)]
pub struct Config {
    pub clients: Option<Vec<Client>>,
    pub port :Option<u32>,
    pub clients_cfg_path :Option<String>,
}

impl Default for Client {
    fn default() -> Self {
        Client {
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
            clients: None,
            port: Some(8433),
            clients_cfg_path: Some(String::from("clients"))
        }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum ErrorLevel {
    Warning,
    Critical,
}

#[derive(Debug, Clone)]
pub struct Error {
    pub key :String,
    pub level :ErrorLevel,
    pub desc :String
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_fmt(format_args!("[{:?}] {} (key: '{}')", self.level, self.desc, self.key))
    }
}
