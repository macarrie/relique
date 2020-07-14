use anyhow::Result;
use log::*;

use crate::lib;
use crate::types;
use crate::types::app::{ReliqueApp, Stopping};
use crate::types::config::ConfigVersion;
use futures::executor::block_on;
use uuid::Uuid;

pub struct ServerDaemon {
    pub config: types::config::Config,
}

impl ReliqueApp for ServerDaemon {
    fn new(config: types::config::Config) -> Result<Self> {
        Ok(ServerDaemon { config })
    }

    fn loop_func(&mut self) -> Result<Stopping> {
        if self.config.clients.is_none() {
            info!("No clients found in configuration");
            return Ok(Stopping::No);
        }

        let client_list = self.config.clients.as_ref().unwrap();
        let config_send_future = send_configuration_to_clients(&self.config, &client_list);
        block_on(config_send_future).unwrap_or_else(|err| {
            warn!("An error occurred when sending configuration to some clients ({}). See previous log entries for more details", err);
        });

        Ok(Stopping::No)
    }
}

async fn send_configuration_to_clients(
    cfg: &types::config::Config,
    clients: &[types::config::Client],
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
    cfg: types::config::Config,
    client: types::config::Client,
) -> Result<Option<Uuid>> {
    let url = format!(
        "https://{address}:{port}/api/v1/config/version",
        address = client.address,
        port = client.port.unwrap()
    );

    let client = lib::web::get_http_client(cfg.ssl_cert, cfg.strict_ssl_certificate_check)?;
    let res = client.get(&url).send().await?;
    let config_version = res.json::<ConfigVersion>().await?;

    Ok(config_version.version)
}

#[actix_rt::main]
async fn send_config_to_client(
    cfg: types::config::Config,
    client: types::config::Client,
) -> Result<()> {
    let url = format!(
        "https://{address}:{port}/api/v1/config",
        address = client.address,
        port = client.port.unwrap()
    );

    let http_client =
        lib::web::get_http_client(cfg.ssl_cert, cfg.strict_ssl_certificate_check).unwrap();
    let res = http_client
        .post(&url)
        .json(&client)
        .send()
        .await?
        .error_for_status();
    if let Err(e) = res {
        return Err(e.into());
    }

    Ok(())
}
