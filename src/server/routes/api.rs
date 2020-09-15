use crate::server::server_daemon::ServerDaemon;
use crate::types::backup_job::{BackupJob, JobStatus};
use actix_web::{post, put, web, HttpResponse, Responder};
use log::*;
use std::sync::RwLock;

#[post("/backup/jobs/register")]
pub async fn backup_jobs_register(
    state: web::Data<RwLock<ServerDaemon>>,
    job: web::Json<BackupJob>,
) -> impl Responder {
    let mut state = state.write().unwrap();
    info!(
        "Registering job '{job}' from client '{client}'",
        job = job,
        client = job.client
    );
    let job_exists = state
        .active_jobs
        .clone()
        .into_iter()
        .filter(|j| j.id == job.id && j.status == JobStatus::Active)
        .count()
        > 0;
    if job_exists {
        return HttpResponse::Conflict().body("Job already registered in relique server");
    }

    state.active_jobs.push(job.clone());
    HttpResponse::Ok().body("Job registered")
}

#[put("/backup/jobs/{id}/status")]
pub async fn update_backup_jobs_status(
    state: web::Data<RwLock<ServerDaemon>>,
    path: web::Path<String>,
    status: web::Json<JobStatus>,
) -> impl Responder {
    let mut state = state.write().unwrap();
    let status = status.into_inner();
    let id = path.into_inner();

    info!(
        "Updating job '{id}' status to {status:?}",
        id = id,
        status = status
    );

    let job_index = state
        .active_jobs
        .iter()
        .position(|j| j.id.to_string() == id);
    if job_index.is_none() {
        return HttpResponse::NotFound().body("Job not found");
    }
    let job_index = job_index.unwrap();

    state.active_jobs[job_index].status = status;
    HttpResponse::Ok().body("Job status updated")
}
