use crate::lib;
use crate::types::backup_module::BackupModule;
use crate::types::backup_type::BackupType;
use crate::types::config;
use crate::types::db;
use crate::types::job_status::JobStatus;
use anyhow::Result;
use log::*;
use rusqlite::named_params;
use serde::{Deserialize, Serialize};
use std::sync::{Arc, RwLock};
use std::{fmt, thread};
use uuid::Uuid;

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct BackupJob {
    pub uuid: Uuid,
    pub client: config::client::Client,
    pub module: BackupModule,
    pub status: JobStatus,
    pub backup_type: BackupType,
}

impl fmt::Display for BackupJob {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self.status {
            JobStatus::Pending => {
                write!(f, "Job {} pending for '{}' (module '{}')", self.uuid, self.module.name, self.module.module_type)
            },
            JobStatus::Active => {
                write!(f, "Job {} running for '{}' (module '{}')", self.uuid, self.module.name, self.module.module_type)
            },
            JobStatus::Done => {
                write!(f, "Job {} done for '{}' (module '{}'). The job will be cleaned at the end its schedule", self.uuid, self.module.name, self.module.module_type)
            },
            JobStatus::Incomplete => {
                write!(f, "Job {} done but some files could not be backup up for '{}' (module '{}'). The job will be cleaned at the end its schedule", self.uuid, self.module.name, self.module.module_type)
            },
            JobStatus::Error => {
                write!(f, "Job {} ended in error for '{}' (module '{}'). See previous logs for more details", self.uuid, self.module.name, self.module.module_type)
            }
        }
    }
}

impl BackupJob {
    pub fn new(client: config::client::Client, module: BackupModule) -> Self {
        BackupJob {
            uuid: Uuid::new_v4(),
            client,
            module: module.clone(),
            status: JobStatus::Pending,
            backup_type: module.backup_type.unwrap_or(BackupType::Full),
        }
    }

    pub fn set_status(&mut self, status: JobStatus) {
        self.status = status;
    }

    pub fn get_active_jobs(pool: db::Pool) -> Result<Vec<BackupJob>> {
        let conn = pool.get()?;
        let mut stmt = conn.prepare(
            "SELECT \
                uuid  \
             FROM jobs \
             WHERE status = :status \
             JOIN clients ON client_id = clients.id \
             JOIN modules ON module_id = modules.id",
        )?;
        let uuids: Vec<Uuid> = stmt
            .query_map_named(named_params![":status": JobStatus::Active], |row| {
                Ok(row.get(0)?)
            })?
            .map(Result::unwrap)
            .collect();

        if uuids.is_empty() {
            return Ok(vec![]);
        }

        let jobs = uuids
            .iter()
            .map(|uuid| BackupJob::find_by_uuid(*uuid, pool.clone()))
            .filter(|res| res.is_ok())
            .map(Result::unwrap)
            .filter(|opt| opt.is_some())
            .map(Option::unwrap)
            .collect();

        Ok(jobs)
    }

    pub fn find_by_uuid(uuid: Uuid, pool: db::Pool) -> Result<Option<BackupJob>> {
        let conn = pool.get()?;
        let mut stmt = conn.prepare(
            "SELECT \
                uuid, \
                status, \
                jobs.backup_type, \
                modules.module_type, \
                modules.name, \
                modules.backup_type, \
                modules.pre_backup_script, \
                modules.post_backup_script, \
                modules.pre_restore_script, \
                modules.post_restore_script, \
                clients.name, \
                clients.address, \
                clients.port, \
                clients.server_address, \
                clients.server_port \
            FROM jobs \
            JOIN clients ON jobs.client_id = clients.id \
            JOIN modules ON jobs.module_id = modules.id \
            WHERE uuid = :uuid",
        )?;
        let jobs: Vec<BackupJob> = stmt
            .query_map_named(named_params![":uuid": uuid], |row| {
                Ok(BackupJob {
                    uuid: row.get(0)?,
                    status: row.get(1)?,
                    backup_type: row.get(2)?,
                    module: BackupModule {
                        module_type: row.get(3)?,
                        name: row.get(3)?,
                        backup_type: row.get(5)?,
                        pre_backup_script: row.get(6)?,
                        post_backup_script: row.get(7)?,
                        pre_restore_script: row.get(8)?,
                        post_restore_script: row.get(9)?,
                        schedules: None,
                        backup_paths: None,
                    },
                    client: config::client::Client {
                        name: row.get(10)?,
                        address: row.get(11)?,
                        port: row.get(12)?,
                        config_version: None,
                        server_address: row.get(13)?,
                        server_port: row.get(14)?,
                        modules: vec![],
                        schedules: vec![],
                    },
                })
            })?
            .map(Result::unwrap)
            .collect();

        if jobs.is_empty() {
            return Ok(None);
        }

        if jobs.len() > 1 {
            anyhow::bail!("Found {len} backup jobs with unique UUID in database");
        }

        Ok(Some(jobs[0].clone()))
    }

    pub fn get_db_id(uuid: Uuid, pool: db::Pool) -> Result<i64> {
        let conn = pool.get()?;
        let mut stmt = conn.prepare("SELECT (id) FROM jobs WHERE uuid = :uuid")?;
        let ids: Vec<i64> = stmt
            .query_map_named(named_params![":uuid": uuid], |row| Ok(row.get(0)?))?
            .map(Result::unwrap)
            .collect();

        if ids.is_empty() {
            anyhow::bail!(
                "Cannot find backup job with unique UUID '{}' in database",
                uuid
            );
        }

        if ids.len() > 1 {
            anyhow::bail!("Found {len} backup jobs with unique UUID in database");
        }

        Ok(ids[0])
    }

    pub fn save(&self, pool: db::Pool) -> Result<i64> {
        debug!("Saving job '{}' into database", self.uuid);
        let module_id = self.module.save(pool.clone())?;
        let client_id = self.client.save(pool.clone())?;
        let params = named_params! {
            ":uuid": self.uuid,
            ":status": self.status,
            ":backup_type": self.backup_type,
            ":module_id": module_id,
            ":client_id": client_id,
        };

        let conn = pool.get()?;

        if let Ok(id) = BackupJob::get_db_id(self.uuid, pool.clone()) {
            conn.execute_named(
                "UPDATE jobs \
                SET status = :status, \
                    backup_type = :backup_type, \
                    module_id = :module_id, \
                    client_id = :client_id \
                WHERE uuid = :uuid",
                params,
            )?;
            return Ok(id);
        }

        conn.execute_named(
            "INSERT INTO jobs (uuid, status, backup_type, module_id, client_id) VALUES (:uuid, :status, :backup_type, :module_id, :client_id)",
            params,
        )?;

        Ok(conn.last_insert_rowid())
    }
}

pub async fn run_job(
    cfg: config::Config,
    client_cfg: config::client::Client,
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
                id = job.uuid,
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
                id = job.uuid,
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
                id = job.uuid,
                client = job.client,
                module = job.module
            );
        }
    });

    // TODO: Remove log
    error!("Launched job, exiting function");

    Ok(())
}
