use anyhow::Result;
use log::*;
use uuid::Uuid;

use crate::lib;
use crate::types::backup_module;
use crate::types::config;
use serde::{Deserialize, Serialize};
use std::sync::{Arc, RwLock};
use std::{fmt, thread};

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct BackupJob {
    pub id: Uuid,
    pub client: config::Client,
    pub module: backup_module::BackupModule,
    pub status: JobStatus,
}

#[derive(PartialEq, Clone, Debug, Serialize, Deserialize)]
pub enum JobStatus {
    Pending,
    Active,
    Done,
    // TODO: Handle backup job errors
    //Error,
}

impl fmt::Display for BackupJob {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self.status {
            JobStatus::Pending => {
                write!(f, "Job {} pending for '{}' (module '{}')", self.id, self.module.name, self.module.module_type)
            },
            JobStatus::Active => {
                write!(f, "Job {} running for '{}' (module '{}')", self.id, self.module.name, self.module.module_type)
            },
            JobStatus::Done => {
                write!(f, "Job {} done for '{}' (module '{}'). The job will be cleaned at the end its schedule", self.id, self.module.name, self.module.module_type)
            },
            //JobStatus::Error => {
                //write!(f, "Job {} in error for '{}' (module '{}'). See previous logs for more details", self.id, self.module.name, self.module.module_type)
            //}
        }
    }
}

impl BackupJob {
    pub fn new(client: config::Client, module: backup_module::BackupModule) -> Self {
        BackupJob {
            id: Uuid::new_v4(),
            client,
            module,
            status: JobStatus::Pending,
        }
    }

    pub fn set_status(&mut self, status: JobStatus) {
        self.status = status;
    }
}

pub async fn run_job(
    cfg: config::Config,
    client_cfg: config::Client,
    job_arc: Arc<RwLock<BackupJob>>,
) -> Result<()> {
    // TODO: Handle thread crash
    let arc_config = Arc::new(RwLock::new(cfg));
    let arc_client_config = Arc::new(RwLock::new(client_cfg));
    thread::spawn(move || {
        {
            let job = job_arc.read().unwrap();
            info!(
                "Starting backup job ID '{id}' on client {client} for module {module}",
                id = job.id,
                client = job.client,
                module = job.module
            );
        }

        job_arc.write().unwrap().set_status(JobStatus::Active);

        // TODO: Launch prebackup script
        {
            let job = job_arc.read().unwrap();
            info!(
                "Launching prebackup script '{path}' for job ID '{id}' on client {client} for module {module}",
                path = job.module.pre_backup_script.as_ref().unwrap_or(&String::from("")),
                id = job.id,
                client = job.client,
                module = job.module
            );
        }

        // TODO: Perform backup
        let lib_job = Arc::clone(&job_arc);
        // TODO: Handle error
        let _start_res = lib::backup::start(arc_config, arc_client_config, lib_job);
        // TODO: Clean async
        //block_on(start_res);
        // TODO: Launch postbackup script
        {
            let job = job_arc.read().unwrap();
            info!(
                "Launching postbackup script '{path}' for job ID '{id}' on client {client} for module {module}",
                path = job.module.post_backup_script.as_ref().unwrap_or(&String::from("")),
                id = job.id,
                client = job.client,
                module = job.module
            );
        }
    });

    // TODO: Remove log
    error!("Launched job, exiting function");

    Ok(())
}
