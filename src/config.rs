use std::error::Error;
use std::fs::File;
use std::io::Read;
use std::path::{Path, PathBuf};

use crate::types::config;
use log::*;

use walkdir::WalkDir;

pub fn load(path: &Path, load_clients_conf: bool) -> Result<config::Config, Box<dyn Error>> {
    if !path.exists() {
        return Err(Box::new(std::io::Error::new(
            std::io::ErrorKind::NotFound,
            format!("Configuration file '{}' not found", path.display()),
        )));
    }

    info!("Loading configuration file: '{}'", path.display());

    let mut file = File::open(path)?;
    let mut cfg_string = String::new();
    file.read_to_string(&mut cfg_string)?;

    let mut cfg: config::Config = toml::from_str(cfg_string.as_str())?;

    if load_clients_conf {
        let clients_cfg_path_str = cfg
            .clone()
            .clients_cfg_path
            .unwrap_or(String::from("clients"));
        let clients_cfg_path = Path::new(&clients_cfg_path_str);

        let mut clients_path: PathBuf;
        if clients_cfg_path.is_relative() {
            clients_path = path.clone().canonicalize().unwrap();
            clients_path.pop();
            clients_path.push(clients_cfg_path);
        } else {
            clients_path = clients_cfg_path.to_path_buf();
        }

        info!(
            "Looking for clients configuration in '{}'",
            clients_path.display()
        );
        let clients = load_clients_configuration(clients_path.as_ref()).unwrap();
        info!(
            "Found {} client declarations in configuration",
            clients.len()
        );

        cfg.clients = Some(clients);
    }

    info!("Successfully loaded configuration");
    Ok(cfg)
}

fn load_clients_configuration(path: &Path) -> Result<Vec<config::Client>, Box<dyn Error>> {
    let mut clients: Vec<config::Client> = Vec::new();

    for entry in WalkDir::new(path) {
        let path_entry = entry.unwrap();
        let ext = path_entry.path().extension();
        match ext {
            Some(e) => {
                if e == "toml" {
                    let mut file = File::open(path_entry.path())?;
                    let mut cfg_string = String::new();
                    file.read_to_string(&mut cfg_string)?;

                    let client: config::Client = toml::from_str(cfg_string.as_str())?;
                    clients.push(client)
                }
            }
            None => {
                continue;
            }
        }
    }

    Ok(clients)
}

pub fn check(cfg: &config::Config) -> Vec<config::Error> {
    debug!("Starting configuration checks");

    let mut errors: Vec<config::Error> = Vec::new();

    let clients_count = match &cfg.clients {
        None => 0,
        Some(clients) => clients.len(),
    };
    if clients_count == 0 {
        warn!("No client declaration found in configuration");
        errors.push(config::Error {
            key: "clients".to_string(),
            level: config::ErrorLevel::Warning,
            desc: "No clients defined".to_string(),
        })
    }

    // TODO: Check client modules empty
    // TODO: Check unknown client modules

    errors
}
