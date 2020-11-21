use anyhow::{anyhow, Result};
use log::*;

use crate::lib;
use crate::types::app::{ReliqueApp, Stopping};
use crate::types::backup_job::BackupJob;
use crate::types::config;
use crate::types::config::config_version::ConfigVersion;
use crate::types::db;
use futures::executor::block_on;
use uuid::Uuid;

pub struct ServerDaemon {
    pub config: config::Config,
    pub db_pool: db::Pool,
}

impl ReliqueApp for ServerDaemon {
    fn new(config: config::Config, db_pool: Option<db::Pool>) -> Result<Self> {
        Ok(ServerDaemon {
            config,
            db_pool: db_pool.unwrap(),
        })
    }

    fn loop_func(&mut self) -> Result<Stopping> {
        error!("APP LOOP FUNC");
        if self.config.clients.is_none() {
            info!("No clients found in configuration");
            return Ok(Stopping::No);
        }

        let client_list = self.config.clients.as_ref().unwrap();
        let config_send_future = send_configuration_to_clients(&self.config, &client_list);
        block_on(config_send_future).unwrap_or_else(|err| {
            warn!("An error occurred when sending configuration to some clients ({}). See previous log entries for more details", err);
        });

        let active_jobs = BackupJob::get_active_jobs(self.db_pool.clone()).unwrap_or(vec![]);
        if active_jobs.is_empty() {
            info!("No active backup jobs on clients");
        } else {
            info!("{} backup jobs active on clients", active_jobs.len());
            for job in active_jobs {
                debug!(
                    "Backup job '{job}' active on client '{client}'",
                    job = job,
                    client = job.client
                );
            }
        }

        Ok(Stopping::No)
    }
}

async fn send_configuration_to_clients(
    cfg: &config::Config,
    clients: &[config::client::Client],
) -> Result<()> {
    for client in clients {
        let client_cfg_version = get_config_version(cfg.clone(), client.clone());
        let client_cfg_version = match client_cfg_version {
            Ok(ver) => ver,
            Err(err) => {
                error!(
                    "Could not get client version for client '{client},{addr}': {e}",
                    client = client.name,
                    addr = client.address,
                    e = err
                );
                continue;
            }
        };

        if client_cfg_version != cfg.config_version {
            send_config_to_client(cfg.clone(), client.clone()).unwrap_or_else(|err| {
                error!(
                    "Could not send configuration to client '{client},{addr}': {e}",
                    client = client.name,
                    addr = client.address,
                    e = err
                );
            });
        }
    }

    Ok(())
}

#[actix_rt::main]
async fn get_config_version(
    cfg: config::Config,
    client: config::client::Client,
) -> Result<Option<Uuid>> {
    let url = format!(
        "https://{address}:{port}/api/v1/config/version",
        address = client.address,
        port = client.port.unwrap()
    );

    let client = lib::web::get_http_client(cfg.strict_ssl_certificate_check)?;
    let res = client.get(&url).send().await;

    match res {
        Ok(mut response) => {
            let config_version = response.json::<ConfigVersion>().await?;
            Ok(config_version.version)
        }
        Err(e) => Err(anyhow!("{}", e)),
    }
}

#[actix_rt::main]
async fn send_config_to_client(cfg: config::Config, client: config::client::Client) -> Result<()> {
    let url = format!(
        "https://{address}:{port}/api/v1/config",
        address = client.address,
        port = client.port.unwrap()
    );

    let http_client = lib::web::get_http_client(cfg.strict_ssl_certificate_check).unwrap();
    let res = http_client.post(&url).send_json(&client).await;

    res.map_err(|e| anyhow!("{}", e)).and(Ok(()))
}
