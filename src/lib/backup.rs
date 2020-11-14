use crate::lib;

use crate::db;
use crate::lib::rsync;
use crate::types::backup_file::BackupFile;
use crate::types::backup_job::BackupJob;
use crate::types::backup_type::BackupType;
use crate::types::config;
use crate::types::config::Config;
use crate::types::job_status::JobStatus;
use crate::types::rsync::Signature;
use anyhow::{anyhow, Result};
use futures::StreamExt;
use log::*;
use std::fs::create_dir_all;
use std::path::Path;
use std::sync::{Arc, RwLock};
use walkdir::WalkDir;

async fn register_job_to_server(
    arc_cfg: Arc<RwLock<Config>>,
    arc_client_cfg: Arc<RwLock<config::client::Client>>,
    job_arc: &Arc<RwLock<BackupJob>>,
) -> Result<()> {
    let cfg = arc_cfg.read().unwrap();
    let client_cfg = arc_client_cfg.read().unwrap();
    let job = job_arc.read().unwrap();
    info!(
        "Register job start to server for {backup_type:?} backup with job ID '{id}' on client {client}, module {module}",
        backup_type = job.backup_type,
        id = job.uuid,
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
    arc_client_cfg: Arc<RwLock<config::client::Client>>,
    job_arc: &Arc<RwLock<BackupJob>>,
) -> Result<()> {
    let cfg = arc_cfg.read().unwrap();
    let client_cfg = arc_client_cfg.read().unwrap();
    let job = job_arc.read().unwrap();
    info!(
        "Update job status in server: job '{id}' with status '{status:?}' on client '{client}'",
        id = job.uuid,
        status = job.status,
        client = job.client,
    );

    let url = format!(
        "https://{address}:{port}/api/v1/backup/jobs/{id}/status",
        address = client_cfg.server_address,
        port = client_cfg.server_port.unwrap_or(8433),
        id = job.uuid,
    );

    let http_client = lib::web::get_http_client(cfg.strict_ssl_certificate_check).unwrap();
    let res = http_client.put(&url).send_json(&job.status).await;

    res.map_err(|e| anyhow!("{}", e)).and(Ok(()))
}

async fn send_files(
    arc_cfg: Arc<RwLock<Config>>,
    arc_client_cfg: Arc<RwLock<config::client::Client>>,
    job_arc: &Arc<RwLock<BackupJob>>,
) -> Result<JobStatus> {
    // TODO: Handle incomplete and error status for job
    let cfg = arc_cfg.read().unwrap();
    let client_cfg = arc_client_cfg.read().unwrap();
    let job = job_arc.read().unwrap();
    let mut return_status = JobStatus::Active;
    for backup_path in job.module.backup_paths.clone().unwrap_or_default() {
        info!(
            "[Job {id}] Performing backup of path {path:?}",
            id = job.uuid,
            path = backup_path.clone()
        );

        let http_client = lib::web::get_http_client(cfg.strict_ssl_certificate_check).unwrap();

        for entry in WalkDir::new(backup_path.clone()) {
            if let Err(error) = entry {
                error!(
                    "[Job {job_id}] Could not browse path '{path}' for backup: {err}",
                    job_id = job.uuid,
                    path = backup_path.clone(),
                    err = error,
                );
                return_status = JobStatus::Incomplete;
                continue;
            }
            let path_entry = entry.unwrap();
            let file_path = path_entry.path();

            if file_path.is_dir() {
                continue;
            }

            let mut backup_file = BackupFile {
                job_id: job.uuid,
                path: file_path.to_str().unwrap_or("").to_string(),
                signature: None,
                delta: None,
                is_dir: file_path.is_dir(),
            };

            // TODO: Handle permissions
            let url = format!(
                "https://{address}:{port}/api/v1/backup/jobs/{id}/signature",
                address = client_cfg.server_address,
                port = client_cfg.server_port.unwrap_or(8433),
                id = job.uuid,
            );

            let res = http_client
                .get(&url)
                .send_json(&backup_file)
                .await
                .map_err(|err| anyhow!("{}", err));
            if let Err(ref e) = res {
                error!(
                    "Could not get signature from server for file '{}': '{}'",
                    backup_file.path, e
                );
                return_status = JobStatus::Incomplete;
                continue;
            }

            let mut sig_result_bytes = actix_web::web::BytesMut::new();
            let mut res = res.unwrap();
            while let Some(item) = res.next().await {
                let item = item.unwrap();
                sig_result_bytes.extend_from_slice(&item);
            }

            let server_signature: serde_json::Result<Signature> =
                serde_json::from_slice(sig_result_bytes.as_ref());
            if let Err(e) = server_signature {
                let err_msg = format!(
                    "Could not parse signature received from server for file '{file}': '{err}'",
                    file = backup_file.path,
                    err = e
                );
                error!("{}", err_msg);
                return_status = JobStatus::Incomplete;
                continue;
            }
            backup_file.signature = Some(server_signature.unwrap());

            let delta = rsync::get_delta(
                backup_file.path.clone(),
                backup_file.signature.clone().unwrap(),
            );
            if let Err(e) = delta {
                let err_msg = format!(
                    "Could not compute delta for file '{file}': '{err}'",
                    file = backup_file.path,
                    err = e
                );
                error!("{}", err_msg);
                return_status = JobStatus::Incomplete;
                continue;
            }

            let url = format!(
                "https://{address}:{port}/api/v1/backup/jobs/{id}/delta",
                address = client_cfg.server_address,
                port = client_cfg.server_port.unwrap_or(8433),
                id = job.uuid,
            );

            backup_file.delta = Some(delta.unwrap());
            let res = http_client.post(&url).send_json(&backup_file).await;
            if let Err(ref e) = res {
                // TODO: Better err message
                error!(
                    "Could not send delta to server for file '{}': '{}'",
                    backup_file.path, e
                );
                return_status = JobStatus::Incomplete;
            }
        }
    }

    if return_status == JobStatus::Active {
        return_status = JobStatus::Done;
    }

    // TODO: Handle error
    Ok(return_status)
}

// TODO: Merge with send delta func
#[actix_rt::main]
pub async fn start(
    cfg: Arc<RwLock<Config>>,
    client_cfg: Arc<RwLock<config::client::Client>>,
    job_arc: Arc<RwLock<BackupJob>>,
) -> Result<()> {
    // TODO: Start job
    {
        let job = job_arc.read().unwrap();
        info!(
            "Launching {backup_type:?} backup for job ID '{id}' on client {client} for module {module}",
            backup_type = job.backup_type,
            id = job.uuid,
            client = job.client,
            module = job.module
        );
    }

    // TODO: Handle error
    register_job_to_server(Arc::clone(&cfg), Arc::clone(&client_cfg), &job_arc).await?;
    // TODO: Handle error
    let send_files_res = send_files(Arc::clone(&cfg), Arc::clone(&client_cfg), &job_arc).await;
    let status = match send_files_res {
        Err(_e) => {
            error!("Could not send all files to relique master server. Check previous logs for more information");
            JobStatus::Error
        }
        Ok(status) => status,
    };

    job_arc.write().unwrap().set_status(status);
    // TODO: Handle error
    update_job_status_to_server(Arc::clone(&cfg), Arc::clone(&client_cfg), &job_arc).await?;

    Ok(())
}

pub fn get_local_backup_file_path(
    cfg: &Config,
    job: &BackupJob,
    backup_file: &BackupFile,
) -> String {
    format!(
        "{main}/{client_name}/{id}/{file_path}",
        main = cfg.backup_storage_path,
        client_name = job.client.name,
        id = job.uuid,
        file_path = backup_file.path
    )
}

pub fn get_diff_reference_file_path(
    config: &Config,
    pool: crate::types::db::Pool,
    job: &mut BackupJob,
    file_path: String,
) -> String {
    match job.backup_type {
        BackupType::Diff => {
            let full_job = db::server::get_previous_full_backup_job(&job, pool.clone());
            if let Err(err) = full_job {
                error!(
                    "Error encountered when querying previous full backup job: '{}'",
                    err
                );

                return String::from("/dev/null");
            }
            let full_job = full_job.unwrap();

            if full_job.is_none() {
                info!("[Job '{id}'] No previous full backup found. Performing full backup instead of diff", id = job.uuid);
                job.backup_type = BackupType::Full;
                // TODO: Handle error
                job.save(pool.clone()).unwrap_or_else(|err| {
                    error!("Could not save status changes for job: '{}'", err);
                    0
                });

                return String::from("/dev/null");
            }
            let full_job = full_job.unwrap();
            return get_local_backup_file_path(
                config,
                &full_job,
                &BackupFile {
                    job_id: full_job.uuid,
                    path: file_path,
                    signature: None,
                    delta: None,
                    is_dir: false,
                },
            );
        }
        BackupType::Full => {
            return String::from("/dev/null");
        }
    }
}

pub async fn save_file_from_remote_delta(
    cfg: Config,
    job: BackupJob,
    backup_file: &BackupFile,
) -> Result<()> {
    info!(
        "[Job '{id}'] Creating backup for file '{file}' from client delta",
        id = job.uuid,
        file = backup_file.path
    );

    // TODO: If null delta, create hard link instead
    let path_str = get_local_backup_file_path(&cfg, &job, backup_file);

    let path = Path::new(&path_str);
    let source_diff_path: String;
    if path.exists() {
        source_diff_path = path_str.clone();
    } else {
        source_diff_path = "/dev/null".to_string();
    }

    let apply_delta_res = rsync::apply_delta(
        source_diff_path.clone(),
        backup_file.delta.clone().unwrap(),
        path,
    );
    if let Err(e) = apply_delta_res {
        let err_msg = format!(
            "Could not apply delta for file '{file}': '{err}'",
            file = source_diff_path,
            err = e
        );
        error!("{}", err_msg);
        return Err(anyhow!("{}", err_msg));
    }

    Ok(())
}

pub async fn create_directory_tree(
    cfg: Config,
    job: BackupJob,
    backup_file: &BackupFile,
) -> Result<()> {
    info!(
        "[Job '{id}'] Creating directory structure for '{file}'",
        id = job.uuid,
        file = backup_file.path
    );

    let path_str = get_local_backup_file_path(&cfg, &job, backup_file);

    let path = Path::new(&path_str);
    create_dir_all(path).map_err(|e| anyhow!("{}", e))
}
