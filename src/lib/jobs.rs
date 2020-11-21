use crate::lib;
use crate::types::backup_job::{BackupJob, BackupJobSearchParameters};
use crate::types::backup_type::BackupType;
use crate::types::config;

use log::*;

use anyhow::{anyhow, Result};

pub async fn list(
    cfg: config::Config,
    client_name: Option<String>,
    module_name: Option<String>,
    backup_type_str: Option<String>,
) -> Result<Vec<BackupJob>> {
    let backup_type: Option<BackupType>;
    if backup_type_str.is_some() {
        let parse_res =
            serde_json::from_str::<BackupType>(&backup_type_str.unwrap_or("full".to_string()));
        backup_type = parse_res
            .map_err(|err| {
                error!(
                    "Could not parse backup type value from parameter: '{err}'",
                    err = err
                );
            })
            .ok();
    } else {
        backup_type = None;
    }

    println!("TODO: List jobs on client '{client:?}', module'{module:?}', backup type '{type:?}'",
        client = client_name,
         module = module_name,
         type = backup_type,
    );

    let url = format!(
        "https://localhost:{port}/api/v1/jobs/search",
        port = cfg.port.unwrap()
    );

    let client = lib::web::get_http_client(cfg.clone().strict_ssl_certificate_check)?;
    let mut res = client
        .get(&url)
        .send_json(&BackupJobSearchParameters {
            client_name,
            module_name,
            backup_type,
            limit: None,
        })
        .await
        .map_err(|e| {
            let msg = format!("Could not get job search results from server: '{}'", e);
            error!("{}", msg);
            anyhow!("{}", msg)
        })?;

    let jobs = res
        .json::<Vec<BackupJob>>()
        .await
        .map_err(|err| anyhow!("{}", err));

    jobs
}
