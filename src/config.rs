use anyhow::Result;
use std::fs::File;
use std::io::Read;
use std::path::{Path, PathBuf};

use crate::types::config;
use log::*;

use crate::types::backup_module::BackupModuleDef;
use crate::types::config::Client;
use crate::types::schedule::Schedule;
use serde::de::DeserializeOwned;
use std::collections::HashMap;
use uuid::Uuid;
use walkdir::WalkDir;

pub fn load(path: &Path, is_server: bool) -> Result<config::Config> {
    if !path.exists() {
        return Err(std::io::Error::new(
            std::io::ErrorKind::NotFound,
            format!("Configuration file '{}' not found", path.display()),
        )
        .into());
    }

    info!("Loading configuration file: '{}'", path.display());

    let mut file = File::open(path)?;
    let mut cfg_string = String::new();
    file.read_to_string(&mut cfg_string)?;

    let mut cfg: config::Config = toml::from_str(cfg_string.as_str())?;
    cfg.config_version = Some(Uuid::new_v4());

    if is_server {
        let clients_path = get_load_path(path, cfg.clone().clients_cfg_path, "clients");

        info!(
            "Looking for clients configuration in '{}'",
            clients_path.display()
        );
        let mut clients = load_folder_configuration::<Client>(clients_path.as_ref()).unwrap();
        let default_module_configuration = load_modules_default_config(&clients).unwrap_or_else(|err| {
            error!("{}", err);
            std::process::exit(exitcode::CONFIG);
        });

        info!(
            "Found {} clients declarations in configuration",
            clients.len()
        );

        let schedule_path = get_load_path(path, cfg.clone().schedules_cfg_path, "schedule");

        info!(
            "Looking for schedule configuration in '{}'",
            schedule_path.display()
        );
        let schedules = load_folder_configuration::<Schedule>(schedule_path.as_ref())
            .unwrap_or_else(|err| {
                error!("{}", err);
                std::process::exit(exitcode::CONFIG);
            });

        info!(
            "Found {} schedule declarations in configuration",
            schedules.len()
        );

        cfg.schedules = Some(schedules.clone());

        for client in &mut clients {
            (*client).schedules = schedules.clone();
            (*client).config_version = cfg.config_version;
            for module in &mut client.modules {
                let default_module_cfg = default_module_configuration.get(&module.module_type);
                if default_module_cfg.is_none() {
                    error!("Could not find default configuration values for module '{}'. Check that this module is installed correctly", module.module_type);
                    continue;
                }

                *module = merge_module_definitions(&mut *module, default_module_cfg.unwrap());
            }
        }

        cfg.clients = Some(clients);
    }

    info!("Successfully loaded configuration");
    Ok(cfg)
}

fn get_load_path(base_path: &Path, path: Option<String>, default: &str) -> PathBuf {
    let path_str = path.unwrap_or_else(|| String::from(default));
    let raw_cfg_path = Path::new(&path_str);

    let mut item_path: PathBuf;
    if raw_cfg_path.is_relative() {
        item_path = base_path.canonicalize().unwrap();
        item_path.pop();
        item_path.push(raw_cfg_path);
    } else {
        item_path = raw_cfg_path.to_path_buf();
    }

    item_path
}

fn merge_module_definitions(def :&mut BackupModuleDef, default_values :&BackupModuleDef) -> BackupModuleDef {
    let mut new = BackupModuleDef {
        ..def.to_owned()
    };

    if new.backup_paths.is_none() {
        new.backup_paths = default_values.backup_paths.clone();
    }

    if new.pre_backup_script.is_none() {
        new.pre_backup_script = default_values.pre_backup_script.clone();
    }

    if new.post_backup_script.is_none() {
        new.post_backup_script = default_values.post_backup_script.clone();
    }

    if new.pre_restore_script.is_none() {
        new.pre_restore_script = default_values.pre_restore_script.clone();
    }

    if new.post_restore_script.is_none() {
        new.post_restore_script = default_values.post_restore_script.clone();
    }

    new
}

fn load_modules_default_config(clients :&[Client]) -> Result<HashMap<String, BackupModuleDef>> {
    let mut defaults: HashMap<String, BackupModuleDef> = HashMap::new();

    for client in clients {
        for module in &client.modules {
            let path_str = format!("/var/lib/relique/modules/{}/default.toml", module.module_type);
            let path = Path::new(&path_str);

            if !path.exists() {
                return Err(std::io::Error::new(
                    std::io::ErrorKind::NotFound,
                    format!("Default module parameters file not found for module '{}' in '{}'", module.module_type, path.display()),
                ).into());
            }

            info!("Loading default parameters for module '{}': '{}'", module.module_type, path.display());

            let mut file = File::open(path)?;
            let mut cfg_string = String::new();
            file.read_to_string(&mut cfg_string)?;

            let def: BackupModuleDef = toml::from_str(cfg_string.as_str())?;
            if !defaults.contains_key(&module.module_type) {
                defaults.insert(module.module_type.clone(), def);
            }
        }
    }

    Ok(defaults)
}

fn load_folder_configuration<T>(path: &Path) -> Result<Vec<T>>
where
    T: DeserializeOwned + Clone,
{
    let mut item_list: Vec<T> = Vec::new();

    for entry in WalkDir::new(path) {
        let path_entry = entry.unwrap();
        let ext = path_entry.path().extension();
        match ext {
            Some(e) => {
                if e == "toml" {
                    let mut file = File::open(path_entry.path())?;
                    let mut cfg_string = String::new();
                    file.read_to_string(&mut cfg_string)?;

                    item_list.push(toml::from_str(&cfg_string)?)
                }
            }
            None => {
                continue;
            }
        }
    }

    Ok(item_list)
}

pub fn check(cfg: &config::Config) -> Vec<config::Error> {
    debug!("Starting configuration checks");

    let mut errors: Vec<config::Error> = Vec::new();

    // Check: No clients
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
        });
    }

    // Check: Client duplicates
    let clients = cfg.clients.clone().unwrap_or_else(Vec::new);
    let mut clients_name_map: HashMap<String, bool> = HashMap::new();
    let mut clients_addr_port_map: HashMap<(String, u32), bool> = HashMap::new();
    for client in clients {
        let port = client.port.unwrap_or(8434);
        // If key already exist, duplicate client name
        if clients_name_map.insert(client.name.clone(), true).is_some() {
            let msg = format!(
                "Client with duplicate names found in configuration: '{}'",
                client.name.clone()
            );
            error!("{}", msg);
            errors.push(config::Error {
                key: "clients.name".to_string(),
                level: config::ErrorLevel::Critical,
                desc: msg,
            });
        }
        if clients_addr_port_map
            .insert((client.address.clone(), port), true)
            .is_some()
        {
            let msg = format!(
                "Client with duplicate address/port found in configuration: '{}:{}'",
                client.address.clone(),
                port
            );
            error!("{}", msg);
            errors.push(config::Error {
                key: "clients.address, clients.port".to_string(),
                level: config::ErrorLevel::Critical,
                desc: msg,
            });
        }
    }
    // TODO: Check client modules empty
    // TODO: Check unknown client modules
    // TODO: Check schedule ranges coherence
    // TODO: Check unknown schedules in clients config

    errors
}
