use crate::types::backup_job::BackupJob;
use crate::types::backup_type::BackupType;
use crate::types::db;
use crate::types::job_status::JobStatus;
use anyhow::{anyhow, Result};
use log::*;
use r2d2_sqlite::SqliteConnectionManager;
use rusqlite::{named_params, Connection, NO_PARAMS};
use std::fs::create_dir_all;
use std::path::Path;
use uuid::Uuid;

pub async fn get_pool() -> Result<db::Pool> {
    info!("Connecting to database");

    let base_path = Path::new("/var/lib/relique/db/");
    create_dir_all(base_path)?;

    let db_file = "server.db";
    let mut pathbuf = base_path.to_path_buf();
    pathbuf.push(db_file);
    let path = pathbuf.as_path();

    init_schema(path).await?;

    let manager = SqliteConnectionManager::file(path);

    // TODO: Move path into settings
    let pool = r2d2::Pool::new(manager).map_err(|e| {
        let msg = format!("Could not create database connection pool: '{}'", e);
        error!("{}", msg);
        anyhow!("{}", msg)
    })?;

    Ok(pool)
}

pub async fn init_schema(db_path: &Path) -> Result<()> {
    debug!("Initializing database schema");

    let conn = Connection::open(db_path)?;

    conn.execute(
        "CREATE TABLE IF NOT EXISTS modules (
            id INTEGER PRIMARY KEY,
            module_type TEXT NOT NULL,
            name TEXT NOT NULL UNIQUE,
            backup_type INTEGER NOT NULL,
            pre_backup_script TEXT,
            post_backup_script TEXT,
            pre_restore_script TEXT,
            post_restore_script TEXT
        )",
        NO_PARAMS,
    )?;

    conn.execute(
        "CREATE TABLE IF NOT EXISTS clients (
             id INTEGER PRIMARY KEY,
             config_version TEXT,
             name TEXT NOT NULL UNIQUE,
             address TEXT NOT NULL,
             port INTEGER NOT NULL,
             server_address INTEGER NOT NULL,
             server_port INTEGER NOT NULL
         )",
        NO_PARAMS,
    )?;

    // Many to many relationship storage table
    conn.execute(
        "CREATE TABLE IF NOT EXISTS modules_schedules (
            schedule_id INTEGER,
            module_id INTEGER,
            FOREIGN KEY(schedule_id) REFERENCES schedules(id),
            FOREIGN KEY(module_id) REFERENCES modules(id)
        )",
        NO_PARAMS,
    )?;

    conn.execute(
        "CREATE TABLE IF NOT EXISTS jobs (
            id INTEGER PRIMARY KEY,
            uuid TEXT NOT NULL UNIQUE,
            status INTEGER NOT NULL,
            backup_type INTEGER NOT NULL,
            module_id INTEGER NOT NULL,
            client_id INTEGER NOT NULL,
            FOREIGN KEY(module_id) REFERENCES modules(id) ON DELETE CASCADE ON UPDATE CASCADE,
            FOREIGN KEY(client_id) REFERENCES clients(id) ON DELETE CASCADE ON UPDATE CASCADE
        )",
        NO_PARAMS,
    )?;

    Ok(())
}

pub fn get_previous_full_backup_job(job: &BackupJob, pool: db::Pool) -> Result<Option<BackupJob>> {
    let params = named_params! {
        ":module_type": job.module.module_type,
        ":client_name": job.client.name,
        ":backup_type": BackupType::Full,
        ":status": JobStatus::Done,
    };

    let conn = pool.get()?;

    let mut stmt = conn.prepare(
        "SELECT \
            uuid \
        FROM jobs \
        JOIN clients ON jobs.client_id = clients.id \
        JOIN modules ON jobs.module_id = modules.id \
        WHERE modules.module_type = :module_type \
            AND jobs.backup_type = :backup_type \
            AND jobs.status = :status \
            AND clients.name = :client_name \
        ORDER BY jobs.id DESC",
    )?;
    let ids: Vec<Uuid> = stmt
        .query_map_named(params, |row| Ok(row.get(0)?))?
        .map(Result::unwrap)
        .collect();

    if ids.is_empty() {
        return Ok(None);
    }

    let id = ids[0];

    BackupJob::find_by_uuid(id, pool)
}
