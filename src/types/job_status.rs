use anyhow::Result;
use rusqlite::types::{FromSql, FromSqlError, FromSqlResult, ToSql, ToSqlOutput, ValueRef};
use serde::{Deserialize, Serialize};

#[derive(PartialEq, Clone, Copy, Debug, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum JobStatus {
    Pending = 0,
    Active,
    Done,
    Incomplete,
    Error,
}

impl JobStatus {
    fn from_i64(value: i64) -> Result<JobStatus> {
        match value {
            0 => Ok(JobStatus::Pending),
            1 => Ok(JobStatus::Active),
            2 => Ok(JobStatus::Done),
            3 => Ok(JobStatus::Incomplete),
            4 => Ok(JobStatus::Error),
            _ => Err(anyhow::anyhow!(
                "Cannot coerce '{}' into enum variant",
                value
            )),
        }
    }
}

impl FromSql for JobStatus {
    fn column_result(value: ValueRef) -> FromSqlResult<Self> {
        JobStatus::from_i64(value.as_i64()?)
            .map_err(|_e| FromSqlError::OutOfRange(value.as_i64().unwrap()))
    }
}

impl ToSql for JobStatus {
    fn to_sql(&self) -> rusqlite::Result<ToSqlOutput> {
        Ok((*self as i64).into())
    }
}
