CREATE TABLE jobs (
	id 			INTEGER PRIMARY KEY,
	uuid 		TEXT NOT NULL UNIQUE,
	status 		INTEGER NOT NULL,
	backup_type INTEGER NOT NULL,
	job_type 	INTEGER NOT NULL,
	done 		INTEGER NOT NULL,
	start_time 	TIMESTAMP,
	end_time 	TIMESTAMP,
    module_type TEXT,
    client_name TEXT,
    repo_name 	TEXT
);

CREATE TABLE images (
	id 					INTEGER PRIMARY KEY,
	uuid 				TEXT NOT NULL UNIQUE,
	created_at 			TIMESTAMP,
    module_type 		TEXT,
    client_name 		TEXT,
    repo_name 			TEXT,
	number_of_elements 	INTEGER,
	number_of_files    	INTEGER,
	number_of_folders  	INTEGER,
	size_on_disk       	INTEGER
);