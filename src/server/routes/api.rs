use crate::lib::backup;
use crate::lib::backup::get_diff_reference_file_path;
use crate::lib::rsync;
use crate::server::server_daemon::ServerDaemon;
use crate::types::backup_file::BackupFile;
use crate::types::backup_job::BackupJob;
use crate::types::job_status::JobStatus;
use actix_web::web::Bytes;
use actix_web::{get, post, put, web, HttpResponse, Responder};
use futures::StreamExt;
use log::*;
use std::path::{Path, PathBuf};
use std::sync::RwLock;
use uuid::Uuid;

#[post("/backup/jobs/register")]
pub async fn post_backup_jobs_register(
    state: web::Data<RwLock<ServerDaemon>>,
    job: web::Json<BackupJob>,
) -> impl Responder {
    let state = state.read().unwrap();
    info!(
        "Registering job '{job}' from client '{client}'",
        job = job,
        client = job.client
    );
    if BackupJob::find_by_uuid(job.uuid, state.db_pool.clone())
        .unwrap_or(None)
        .is_some()
    {
        return HttpResponse::Conflict().body("Job already registered in relique server");
    }

    let save_res = job.save(state.db_pool.clone()).map_err(|e| {
        let msg = format!("Could not save job '{}' into database: '{}'", job.uuid, e);
        error!("{}", msg);
        HttpResponse::InternalServerError().body(msg)
    });
    if let Err(http_err) = save_res {
        return http_err;
    }

    HttpResponse::Ok().body("Job registered")
}

#[put("/backup/jobs/{id}/status")]
pub async fn put_backup_jobs_id_status(
    state: web::Data<RwLock<ServerDaemon>>,
    path: web::Path<String>,
    status: web::Json<JobStatus>,
) -> impl Responder {
    let state = state.read().unwrap();
    let status = status.into_inner();
    let id = path.into_inner();
    let uuid = Uuid::parse_str(&id);
    if let Err(err) = uuid {
        return HttpResponse::BadRequest()
            .body(format!("Cannot parse valid UUID from '{}': '{}'", id, err));
    }
    let uuid = uuid.unwrap();

    info!(
        "Updating job '{id}' status to {status:?}",
        id = id,
        status = status
    );

    let job = BackupJob::find_by_uuid(uuid, state.db_pool.clone());
    if let Err(err) = job {
        error!(
            "An error occurred when querying job with uuid '{}' in database: '{}'",
            id, err
        );
        return HttpResponse::InternalServerError().body("");
    }
    let job = job.unwrap();
    if job.is_none() {
        return HttpResponse::NotFound().body("Job not found");
    }
    let mut job = job.unwrap();
    job.status = status;
    job.save(state.db_pool.clone()).unwrap_or_else(|err| {
        error!("Could not save status changes for job: '{}'", err);
        0
    });

    HttpResponse::Ok().body("Job status updated")
}

#[get("/backup/jobs/{id}/signature")]
pub async fn get_backup_jobs_id_signature(
    state: web::Data<RwLock<ServerDaemon>>,
    path: web::Path<String>,
    mut payload: web::Payload,
) -> impl Responder {
    let state = state.read().unwrap();
    let id = path.into_inner();
    let uuid = Uuid::parse_str(&id);
    if let Err(err) = uuid {
        return HttpResponse::BadRequest()
            .body(format!("Cannot parse valid UUID from '{}': '{}'", id, err));
    }
    let uuid = uuid.unwrap();

    let job_opt = BackupJob::find_by_uuid(uuid, state.db_pool.clone()).unwrap_or(None);
    if job_opt.is_none() {
        return HttpResponse::NotFound().body("Job not found");
    }
    let mut job = job_opt.unwrap();

    let mut bytes = web::BytesMut::new();
    while let Some(item) = payload.next().await {
        let item = item.unwrap();
        bytes.extend_from_slice(&item);
    }

    let backup_file: serde_json::Result<BackupFile> = serde_json::from_slice(bytes.as_ref());
    if let Err(e) = backup_file {
        error!(
            "Could not parse backup file information from client: '{err}'",
            err = e
        );
        return HttpResponse::BadRequest().body("Could not parse backup file information");
    }
    let backup_file = backup_file.unwrap();

    // TODO: If full backup, /dev/null is OK. Else, get real file path to compute signature
    let diff_reference_file_path = get_diff_reference_file_path(
        &state.config,
        state.db_pool.clone(),
        &mut job,
        backup_file.path,
    );

    let signature = rsync::get_signature(&Path::new(&diff_reference_file_path));
    if let Err(e) = signature {
        let err_msg = format!(
            "Could not get signature for file '{file}': '{err}'",
            file = diff_reference_file_path,
            err = e
        );
        error!("{}", err_msg);
        return HttpResponse::InternalServerError().body(err_msg);
    }
    let signature = signature.unwrap();
    let serialized_signature = serde_json::to_string(&signature).unwrap();
    let (tx, rx) = actix_utils::mpsc::channel::<Result<Bytes, actix_web::Error>>();
    let send_res = tx.send(Ok(Bytes::copy_from_slice(&serialized_signature.as_bytes())));
    if let Err(e) = send_res {
        let err_msg = format!("Could not send signature through channel: '{err}'", err = e);
        error!("{}", err_msg);
        return HttpResponse::InternalServerError().body(err_msg);
    }

    HttpResponse::Ok()
        .content_type("application/json")
        .streaming(rx)
}

#[post("/backup/jobs/{id}/delta")]
pub async fn post_backup_jobs_id_delta(
    state: web::Data<RwLock<ServerDaemon>>,
    path: web::Path<String>,
    mut payload: web::Payload,
) -> impl Responder {
    let state = state.read().unwrap();
    let id = path.into_inner();
    let uuid = Uuid::parse_str(&id);
    if let Err(err) = uuid {
        return HttpResponse::BadRequest()
            .body(format!("Cannot parse valid UUID from '{}': '{}'", id, err));
    }
    let uuid = uuid.unwrap();

    let job_opt = BackupJob::find_by_uuid(uuid, state.db_pool.clone()).unwrap_or(None);
    if job_opt.is_none() {
        return HttpResponse::NotFound().body("Job not found");
    }
    let job = job_opt.unwrap();

    let mut bytes = web::BytesMut::new();
    while let Some(item) = payload.next().await {
        let item = item.unwrap();
        bytes.extend_from_slice(&item);
    }

    let backup_file: serde_json::Result<BackupFile> = serde_json::from_slice(bytes.as_ref());
    if let Err(e) = backup_file {
        error!(
            "Could not parse backup file information from client: '{err}'",
            err = e
        );
        return HttpResponse::BadRequest().body("Could not parse backup file information");
    }

    let backup_file = backup_file.unwrap();

    let mut folder_path_buf = PathBuf::from(backup_file.path.clone());
    folder_path_buf.pop();
    let folder_path = folder_path_buf.as_path();

    let folder_backup_file = BackupFile {
        job_id: job.uuid,
        path: String::from(folder_path.to_str().unwrap()),
        signature: None,
        delta: None,
        // TODO: Remove
        is_dir: false,
    };
    let dir_tree_res =
        backup::create_directory_tree(state.config.clone(), job.clone(), &folder_backup_file).await;

    if let Err(e) = dir_tree_res {
        let msg = format!(
            "[Job '{job_id}'] Could not create directory structure for '{dir}': '{err}'",
            job_id = job.uuid,
            dir = backup_file.path,
            err = e
        );
        error!("{}", msg);
        return HttpResponse::InternalServerError().body(msg);
    }

    let save_from_delta_res =
        backup::save_file_from_remote_delta(state.config.clone(), job.clone(), &backup_file).await;

    if let Err(e) = save_from_delta_res {
        let msg = format!(
            "[Job '{job_id}'] Could not create backup file from client delta for path '{dir}': '{err}'",
            job_id = job.uuid,
            dir = backup_file.path,
            err = e
        );
        error!("{}", msg);
        return HttpResponse::InternalServerError().body(msg);
    }

    HttpResponse::Ok().body("Delta applied")
}
