use anyhow::Result;
use log::*;
use uuid::Uuid;

use crate::types::backup_module;
use crate::types::config;
use std::sync::{Arc, RwLock};
use std::time;
use std::{fmt, thread};

#[derive(Clone, Debug)]
pub struct BackupJob {
    pub id: Uuid,
    pub client: config::Client,
    pub module: backup_module::BackupModuleDef,
    pub status: JobStatus,
}

#[derive(Clone, Debug)]
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
    pub fn new(client: config::Client, module: backup_module::BackupModuleDef) -> Self {
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

pub fn run_job(job_arc: Arc<RwLock<BackupJob>>) -> Result<()> {
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
        // TODO: Perform backup
        for _i in 0..130 {
            {
                let job = job_arc.read().unwrap();
                warn!("Job {id} run loop iter", id = job.id);
            }
            thread::sleep(time::Duration::from_secs(1))
        }
        // TODO: Launch postbackup script

        job_arc.write().unwrap().set_status(JobStatus::Done);
    });

    // TODO: Remove log
    error!("Launched job, exiting function");

    Ok(())
}
