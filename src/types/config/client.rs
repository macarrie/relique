use crate::types::backup_module::BackupModule;
use crate::types::db;
use crate::types::schedule::Schedule;
use anyhow::Result;
use log::*;
use rusqlite::named_params;
use serde::{Deserialize, Serialize};
use std::fmt;
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(default)]
pub struct Client {
    pub name: String,
    pub address: String,
    pub port: Option<u32>,
    pub modules: Vec<BackupModule>,
    pub schedules: Vec<Schedule>,
    pub config_version: Option<Uuid>,
    pub server_address: String,
    pub server_port: Option<u32>,
}

impl Default for Client {
    fn default() -> Self {
        Client {
            config_version: None,
            name: "".to_string(),
            address: "".to_string(),
            port: Some(8434),
            server_address: "".to_string(),
            server_port: Some(8433),
            modules: vec![],
            schedules: vec![],
        }
    }
}

impl fmt::Display for Client {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{},{}", self.name, self.address)
    }
}

impl Client {
    pub fn find_by_id(id: i64, pool: db::Pool) -> Result<Client> {
        let conn = pool.get()?;
        let mut stmt = conn.prepare(
            "SELECT ( config_version, name, address, port, server_address, server_port )
                 FROM clients \
                 WHERE id = :id",
        )?;
        let clients: Vec<Client> = stmt
            .query_map_named(named_params![":id": id], |row| {
                Ok(Client {
                    config_version: row.get(0)?,
                    name: row.get(1)?,
                    address: row.get(2)?,
                    port: row.get(3)?,
                    server_address: row.get(4)?,
                    server_port: row.get(5)?,
                    // TODO: Load modules
                    modules: vec![],
                    // TODO: Load schedules
                    schedules: vec![],
                })
            })?
            .map(Result::unwrap)
            .collect();

        if clients.is_empty() {
            anyhow::bail!("Cannot find client with unique ID '{}' in database", id);
        }

        if clients.len() > 1 {
            anyhow::bail!(
                "Found {len} clients with unique ID '{id}' in database",
                len = clients.len(),
                id = id,
            );
        }

        Ok(clients[0].clone())
    }

    pub fn get_db_id(name: &String, pool: db::Pool) -> Result<i64> {
        let conn = pool.get()?;
        let mut stmt = conn.prepare("SELECT (id) FROM clients WHERE name = :name")?;
        let ids: Vec<i64> = stmt
            .query_map_named(named_params![":name": name], |row| Ok(row.get(0)?))?
            .map(Result::unwrap)
            .collect();

        if ids.is_empty() {
            anyhow::bail!("Cannot find client with name '{}' in database", name);
        }

        if ids.len() > 1 {
            anyhow::bail!("Found {len} clients with name '{}' in database",);
        }

        Ok(ids[0])
    }

    pub fn save(&self, pool: db::Pool) -> Result<i64> {
        debug!("Saving client '{}' into database", self.name);
        let params = named_params! {
            ":config_version": self.config_version,
            ":name": self.name,
            ":address": self.address,
            ":port": self.port,
            ":server_address": self.server_address,
            ":server_port": self.server_port,
        };

        let conn = pool.get()?;

        if let Ok(id) = Client::get_db_id(&self.name, pool.clone()) {
            // UPDATE
            let update_params = named_params! {
                ":id": id,
                ":config_version": self.config_version,
                ":name": self.name,
                ":address": self.address,
                ":port": self.port,
                ":server_address": self.server_address,
                ":server_port": self.server_port,
            };
            conn.execute_named(
                "UPDATE clients \
                SET config_version = :config_version, \
                    name = :name, \
                    address = :address, \
                    port = :port, \
                    server_address = :server_address, \
                    server_port = :server_port \
                WHERE id = :id",
                update_params,
            )?;
            return Ok(id);
        }

        conn.execute_named(
            "INSERT INTO clients ( config_version, name, address, port, server_address, server_port ) \
            VALUES ( :config_version, :name, :address, :port, :server_address, :server_port )",
            params,
        )?;

        Ok(conn.last_insert_rowid())
    }
}
