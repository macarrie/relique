use anyhow::Result;
use log::*;

use crate::types::app::{ReliqueApp, Stopping};
use crate::types::config;

pub struct ClientDaemon {
    pub config: config::Config,
    pub client_config: Option<config::Client>,
}

impl ReliqueApp for ClientDaemon {
    fn new(config: config::Config) -> Result<Self> {
        Ok(ClientDaemon {
            config,
            client_config: None,
        })
    }

    fn loop_func(&mut self) -> Result<Stopping> {
        if self.client_config.is_none() {
            info!("Waiting for configuration from relique server");
            return Ok(Stopping::No);
        }

        Ok(Stopping::No)
    }
}
