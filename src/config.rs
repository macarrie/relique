use std::path::{Path, PathBuf};
use std::error::Error;
use std::fs::File;
use std::io::{Read};

use crate::types::config;

use walkdir::WalkDir;

pub fn load(path :&Path) -> Result<config::Config, Box<dyn Error>> {
    if !path.exists() {
        return Err(Box::new(std::io::Error::new(
            std::io::ErrorKind::NotFound,
            format!("Configuration file '{}' not found", path.display())
        )));
    }

    let mut file = File::open(path)?;
    let mut cfg_string = String::new();
    file.read_to_string(&mut cfg_string)?;

    let mut cfg :config::Config = toml::from_str(cfg_string.as_str())?;

    let clients_cfg_path_str = cfg.clone().clients_cfg_path.unwrap_or(String::from("clients"));
    let clients_cfg_path = Path::new(&clients_cfg_path_str);
    println!("Clients cfg path: {}", clients_cfg_path.display());

    let mut clients_path :PathBuf;
    if clients_cfg_path.is_relative() {
        clients_path = path.clone().canonicalize().unwrap();
        clients_path.pop();
        clients_path.push(clients_cfg_path);
    } else {
        clients_path = clients_cfg_path.to_path_buf();
    }
    println!("Final clients path: {}", clients_path.display());

    let clients = load_clients_configuration(clients_path.as_ref()).unwrap();
    println!("Clients path: {:?}", clients);

    cfg.clients = Some(clients);

    Ok(cfg)
}

fn load_clients_configuration(path :&Path) -> Result<Vec<config::Client>, Box<dyn Error>> {
    let mut clients :Vec<config::Client> = Vec::new();

    for entry in WalkDir::new(path) {
        let path_entry = entry.unwrap();
        let ext = path_entry.path().extension();
        match ext {
            Some(e) => {
                if e == "toml" {
                    let mut file = File::open(path_entry.path())?;
                    let mut cfg_string = String::new();
                    file.read_to_string(&mut cfg_string)?;

                    let client :config::Client = toml::from_str(cfg_string.as_str())?;
                    clients.push(client)
                }
            },
            None => {
                continue;
            }
        }
    }

    Ok(clients)
}

pub fn check(cfg :config::Config) -> Vec<config::Error> {
    let mut errors :Vec<config::Error> = Vec::new();

    if cfg.clients.is_none() || cfg.clients.unwrap_or_default().is_empty() {
        // TODO: Log
        errors.push(config::Error {
            key: "clients".to_string(),
            level: config::ErrorLevel::Critical,
            desc: "No clients defined".to_string(),
        })
    }

    errors
}
