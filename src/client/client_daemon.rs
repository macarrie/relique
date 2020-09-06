use anyhow::Result;
use log::*;

use crate::types::app::{ReliqueApp, Stopping};
use crate::types::backup_job::run_job;
use crate::types::backup_job::{BackupJob, JobStatus};
use crate::types::config;
use std::sync::{Arc, RwLock};

pub struct ClientDaemon {
    pub config: config::Config,
    pub client_config: Option<config::Client>,
    pub jobs: Vec<Arc<RwLock<BackupJob>>>,
}

impl ReliqueApp for ClientDaemon {
    fn new(config: config::Config) -> Result<Self> {
        Ok(ClientDaemon {
            config,
            client_config: None,
            jobs: Vec::new(),
        })
    }

    fn loop_func(&mut self) -> Result<Stopping> {
        if self.client_config.is_none() {
            info!("Waiting for configuration from relique server");
            return Ok(Stopping::No);
        }

        if has_active_schedule(&self.client_config) {
            start_backup_jobs(self).unwrap();
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

fn create_backup_jobs(cfg: &Option<config::Client>) -> Vec<Arc<RwLock<BackupJob>>> {
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

    // TODO: Check module schedules before starting jobs
    for module in &cfg.modules {
        jobs.push(Arc::new(RwLock::new(BackupJob::new(
            cfg.clone(),
            module.clone(),
        ))));
    }

    jobs
}

fn start_backup_jobs(state: &mut ClientDaemon) -> Result<()> {
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
        let res = run_job(thread_job);
        if let Err(e) = res {
            error!("Error encountered when running backup job: '{}", e);
        }

        state.jobs.push(push_job);
    }

    Ok(())
}

fn has_active_schedule(client: &Option<config::Client>) -> bool {
    if client.is_none() {
        return false;
    }

    let cfg_schedules = client.clone().unwrap().schedules;

    // Check if at least one schedule is active
    let active_schedules: Vec<String> = cfg_schedules
        .iter()
        .filter(|sched| sched.is_active())
        .map(|sched| sched.name.clone())
        .collect();

    if active_schedules.is_empty() {
        debug!("No active schedules");
        return false;
    }

    info!("Active schedules: {}", active_schedules.join(", "));

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
