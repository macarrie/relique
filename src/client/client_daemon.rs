use anyhow::Result;
use log::*;

use crate::types::app::{ReliqueApp, Stopping};
use crate::types::backup_job::run_job;
use crate::types::backup_job::BackupJob;
use crate::types::config;
use crate::types::db;
use crate::types::job_status::JobStatus;
use crate::types::schedule::Schedule;
use futures::executor::block_on;
use std::sync::{Arc, RwLock};

// TODO: Add last ping and alert if no ping from server
pub struct ClientDaemon {
    pub config: config::Config,
    pub client_config: Option<config::client::Client>,
    pub jobs: Vec<Arc<RwLock<BackupJob>>>,
    pub db_pool: Option<db::Pool>,
}

impl ReliqueApp for ClientDaemon {
    fn new(config: config::Config, _db_pool: Option<db::Pool>) -> Result<Self> {
        Ok(ClientDaemon {
            config,
            client_config: None,
            jobs: Vec::new(),
            db_pool: None,
        })
    }

    fn loop_func(&mut self) -> Result<Stopping> {
        error!("APP LOOP FUNC");
        if self.client_config.is_none() {
            info!("Waiting for configuration from relique server");
            return Ok(Stopping::No);
        }

        if has_active_schedule(&self.client_config) {
            let jobs_future = start_backup_jobs(self).await?;

            for job_arc in &self.jobs {
                let job = job_arc.read().unwrap();
                info!("{}", job);
            }
        } else {
            clean_done_jobs(self).unwrap();
        }

        Ok(Stopping::No)
    }
}

fn create_backup_jobs(cfg: &Option<config::client::Client>) -> Vec<Arc<RwLock<BackupJob>>> {
    let cfg = match cfg {
        None => return Vec::new(),
        Some(c) => c,
    };

    let mut jobs: Vec<Arc<RwLock<BackupJob>>> = Vec::new();
    if cfg.modules.is_empty() {
        warn!(
            "No modules defined for client {}. No backup jobs to launch",
            cfg
        );
    }

    for module in &cfg.modules {
        let schedules = cfg
            .schedules
            .clone()
            .into_iter()
            .filter(|s| {
                module
                    .schedules
                    .as_ref()
                    .unwrap_or(&vec![])
                    .contains(&s.name)
            })
            .collect();
        let active_schedules = get_active_schedules(schedules);
        if !active_schedules.is_empty() {
            jobs.push(Arc::new(RwLock::new(BackupJob::new(
                cfg.clone(),
                module.clone(),
            ))));
        }
    }

    jobs
}

async fn start_backup_jobs(state: &mut ClientDaemon) -> Result<()> {
    let jobs = create_backup_jobs(&state.client_config);
    for job_arc in jobs {
        let job = job_arc.read().unwrap();
        let job_already_exist = state
            .jobs
            .iter()
            .filter(|j| j.read().unwrap().module.name == job.module.name)
            .count()
            != 0;
        if !state.jobs.is_empty() && job_already_exist {
            continue;
        }

        let push_job = Arc::clone(&job_arc);
        let thread_job = Arc::clone(&job_arc);
        // TODO: Handle error
        let res = run_job(
            state.config.clone(),
            state.client_config.as_ref().unwrap().clone(),
            thread_job,
        )
        .await;
        if let Err(e) = res {
            error!("Error encountered when running backup job: '{}", e);
        }

        state.jobs.push(push_job);
    }

    Ok(())
}

fn get_active_schedules(schedules: Vec<Schedule>) -> Vec<String> {
    // Check if at least one schedule is active
    let active_schedules: Vec<String> = schedules
        .iter()
        .filter(|sched| sched.is_active())
        .map(|sched| sched.name.clone())
        .collect();

    active_schedules
}

fn has_active_schedule(client: &Option<config::client::Client>) -> bool {
    if client.is_none() {
        return false;
    }

    let cfg_schedules = client.clone().unwrap().schedules;
    let active_schedules = get_active_schedules(cfg_schedules);

    if active_schedules.is_empty() {
        debug!("No active schedules");
        return false;
    }

    debug!("Active schedules: {}", active_schedules.join(", "));

    true
}

fn clean_done_jobs(state: &mut ClientDaemon) -> Result<()> {
    let initial_job_count = state.jobs.len();
    state.jobs.retain(|j| {
        if let JobStatus::Done = j.read().unwrap().status {
            return false;
        }

        true
    });

    if state.jobs.len() != initial_job_count {
        debug!(
            "Cleaned {} finished backup jobs from job pool",
            initial_job_count - state.jobs.len()
        );
    }

    Ok(())
}
