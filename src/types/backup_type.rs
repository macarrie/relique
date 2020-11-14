use rusqlite::types::{FromSql, FromSqlError, FromSqlResult, ToSql, ToSqlOutput, ValueRef};
use rusqlite::Result;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone, Copy, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum BackupType {
    Full = 0,
    Diff,
}

impl BackupType {
    fn from_i64(value: i64) -> anyhow::Result<BackupType> {
        match value {
            0 => Ok(BackupType::Full),
            1 => Ok(BackupType::Diff),
            _ => Err(anyhow::anyhow!(
                "Cannot coerce '{}' into BackupType enum variant",
                value
            )),
        }
    }
}

impl FromSql for BackupType {
    fn column_result(value: ValueRef) -> FromSqlResult<Self> {
        BackupType::from_i64(value.as_i64()?)
            .map_err(|_e| FromSqlError::OutOfRange(value.as_i64().unwrap()))
    }
}

impl ToSql for BackupType {
    fn to_sql(&self) -> Result<ToSqlOutput> {
        Ok((*self as i64).into())
    }
}
