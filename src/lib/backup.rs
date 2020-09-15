use crate::lib;

use crate::types::backup_job::{BackupJob, JobStatus};
use crate::types::config::{Client, Config};
use anyhow::{anyhow, Result};
use log::*;
use std::sync::{Arc, RwLock};
use std::{thread, time};
use walkdir::WalkDir;

async fn register_job_to_server(
    arc_cfg: Arc<RwLock<Config>>,
    arc_client_cfg: Arc<RwLock<Client>>,
    job_arc: &Arc<RwLock<BackupJob>>,
) -> Result<()> {
    // TODO: Register backup job into server (for monitoring and interface display)

    let cfg = arc_cfg.read().unwrap();
    let client_cfg = arc_client_cfg.read().unwrap();
    let job = job_arc.read().unwrap();
    info!(
        "Register job start to server for {backup_type} backup with job ID '{id}' on client {client}, module {module}",
        backup_type = job.module.backup_type.as_ref().unwrap_or(&String::from("")),
        id = job.id,
        client = job.client,
        module = job.module
    );

    let url = format!(
        "https://{address}:{port}/api/v1/backup/jobs/register",
        address = client_cfg.server_address,
        port = client_cfg.server_port.unwrap_or(8433),
    );

    let http_client = lib::web::get_http_client(cfg.strict_ssl_certificate_check).unwrap();
    let res = http_client.post(&url).send_json(&job_arc).await;

    res.map_err(|e| anyhow!("{}", e)).and(Ok(()))
}

async fn update_job_status_to_server(
    arc_cfg: Arc<RwLock<Config>>,
    arc_client_cfg: Arc<RwLock<Client>>,
    job_arc: &Arc<RwLock<BackupJob>>,
) -> Result<()> {
    let cfg = arc_cfg.read().unwrap();
    let client_cfg = arc_client_cfg.read().unwrap();
    let job = job_arc.read().unwrap();
    info!(
        "Update job status in server: job '{id}' with status '{status:?}' on client '{client}'",
        id = job.id,
        status = job.status,
        client = job.client,
    );

    let url = format!(
        "https://{address}:{port}/api/v1/backup/jobs/{id}/status",
        address = client_cfg.server_address,
        port = client_cfg.server_port.unwrap_or(8433),
        id = job.id,
    );

    let http_client = lib::web::get_http_client(cfg.strict_ssl_certificate_check).unwrap();
    let res = http_client.put(&url).send_json(&job.status).await;

    res.map_err(|e| anyhow!("{}", e)).and(Ok(()))
}

fn perform_backup(job_arc: &Arc<RwLock<BackupJob>>) {
    let job = job_arc.read().unwrap();
    for backup_path in job.module.backup_paths.clone().unwrap_or_default() {
        for entry in WalkDir::new(backup_path.clone()) {
            if let Err(error) = entry {
                error!(
                    "[Job {job_id}] Could not browse path '{path}' for backup: {err}",
                    job_id = job.id,
                    path = backup_path.clone(),
                    err = error,
                );
                continue;
            }
            let path_entry = entry.unwrap();
            let file_path = path_entry.path();
            warn!("Entry: {:?}", file_path.display());

            //let tmp_path :&Path = file_path.strip_prefix(source_dir).unwrap();
            //let tmp_bkp_path = format!("{}/{}", dest_dir, tmp_path.display());

            //let bkp_file_path = Path::new(&tmp_bkp_path);
            //println!("{}", file_path.display());
            //println!("{}", bkp_file_path.display());

            //if file_path.is_dir() {
            //std::fs::create_dir_all(bkp_file_path).unwrap();
            //} else {
            //hard_link(file_path, bkp_file_path).unwrap();
            //}
        }
    }
}

#[actix_rt::main]
pub async fn start(
    cfg: Arc<RwLock<Config>>,
    client_cfg: Arc<RwLock<Client>>,
    job_arc: Arc<RwLock<BackupJob>>,
) -> Result<()> {
    // TODO: Start job
    {
        let job = job_arc.read().unwrap();
        info!(
            "Launching {backup_type} backup for job ID '{id}' on client {client} for module {module}",
            backup_type = job.module.backup_type.as_ref().unwrap_or(&String::from("")),
            id = job.id,
            client = job.client,
            module = job.module
        );
    }

    // TODO: Handle error
    register_job_to_server(Arc::clone(&cfg), Arc::clone(&client_cfg), &job_arc).await?;

    // TODO: Handle error
    perform_backup(&job_arc);

    for _i in 0..35 {
        {
            let job = job_arc.read().unwrap();
            warn!("Job {id} run loop iter", id = job.id);
        }
        thread::sleep(time::Duration::from_secs(1))
    }

    job_arc.write().unwrap().set_status(JobStatus::Done);
    // TODO: Handle error
    update_job_status_to_server(Arc::clone(&cfg), Arc::clone(&client_cfg), &job_arc).await?;

    Ok(())
}
