use crate::types::rsync::Delta;
use crate::types::rsync::Signature;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
pub struct BackupFile {
    pub job_id: Uuid,
    pub path: String,
    pub signature: Option<Signature>,
    pub delta: Option<Delta>,
    // TODO: Remove
    pub is_dir: bool,
}
