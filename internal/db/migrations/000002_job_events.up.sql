CREATE TABLE job_events (
	id 				INTEGER PRIMARY KEY,
	job_uuid 		TEXT NOT NULL,
	time 			TIMESTAMP NOT NULL,
	type 			INTEGER NOT NULL,
	description 	TEXT
);

CREATE INDEX idx_job_events_job_uuid ON job_events(job_uuid);
