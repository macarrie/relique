# Relique

Relique backs up and restores files from remote clients over rsync and SSH, and keeps track of what each run did.

## Language

### Jobs

**Job**:
One backup or restore run for a client and module, with a single final status.
_Avoid_: Run, task, execution

**Backup** / **Restore**:
The two job types: a backup copies data from a client into a repository, a restore copies an image's data back to a client.
_Avoid_: Using "backup" to mean the file transfer phase of any job

**Backup Path**:
One path on the client listed by a module; a job transfers each backup path independently.

**Sync**:
The file transfer of one backup path within a job, in either direction.
_Avoid_: Backup (for the transfer itself), rsync task, file sync

**Retry**:
A new job started with the same parameters as a finished backup job; it is a distinct job, not a continuation.

### Job Events

**Job Event**:
A timestamped record of a state change in exactly one job, with a type and a human-readable description.
_Avoid_: Event (unqualified), log entry, job log

**Job Event Type**:
The kind of a job event; it alone conveys how serious the event is.
_Avoid_: Level, severity
