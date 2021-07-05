package relique_job

import (
	"database/sql"
	"testing"

	"github.com/macarrie/relique/internal/types/job_status"
	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/macarrie/relique/internal/db"
	log "github.com/macarrie/relique/internal/logging"

	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/job_type"
	"github.com/macarrie/relique/internal/types/module"
)

func SetupTest(t *testing.T) {
	log.Setup(true, log.TEST_LOG_PATH)
	if err := db.InitTestDB(); err != nil {
		t.Errorf("cannot open db: '%s'", err)
	}
}

func TestReliqueJob_Update(t *testing.T) {
	SetupTest(t)

	newModule := module.Module{
		ModuleType:        "jobUpdateNewModule",
		Name:              "jobUpdateNewModule",
		BackupType:        backup_type.BackupType{Type: backup_type.Full},
		Schedules:         nil,
		ScheduleNames:     nil,
		BackupPaths:       []string{"/path1", "/path2"},
		PreBackupScript:   "",
		PostBackupScript:  "",
		PreRestoreScript:  "",
		PostRestoreScript: "",
	}
	newClient := client.Client{
		Name:          "jobUpdateNewClient",
		Address:       "jobUpdateNewClient",
		Port:          8434,
		Modules:       nil,
		Version:       "version",
		ServerAddress: "relique-server",
		ServerPort:    8433,
	}

	backupJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Backup})
	backupJobID, err := backupJob.Save()
	if backupJobID == 0 || err != nil {
		t.Errorf("Cannot save backupJob for TestReliqueJob_Update setup")
	}
	restoreJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Restore})
	restoreJobID, err := restoreJob.Save()
	if restoreJobID == 0 || err != nil {
		t.Errorf("Cannot save restoreJob for TestReliqueJob_Update setup")
	}

	missingModuleJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Backup})
	missingModuleJobID, err := missingModuleJob.Save()
	if missingModuleJobID == 0 || err != nil {
		t.Errorf("Cannot save missingModuleJobID for TestReliqueJob_Update setup")
	}
	missingModuleJob.Module = module.Module{}

	missingClientJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Backup})
	missingClientJobID, err := missingClientJob.Save()
	if missingClientJobID == 0 || err != nil {
		t.Errorf("Cannot save missingClientJobID for TestReliqueJob_Update setup")
	}
	missingClientJob.Client = &client.Client{}

	tests := []struct {
		name    string
		job     ReliqueJob
		inTx    bool
		want    int64
		wantErr bool
	}{
		{
			name:    "backupJob_update_regular_save",
			job:     backupJob,
			inTx:    false,
			want:    backupJobID,
			wantErr: false,
		},
		{
			name:    "restoreJob_update_regular_save",
			job:     restoreJob,
			inTx:    false,
			want:    restoreJobID,
			wantErr: false,
		},
		{
			name:    "missingModuleJob_update_regular_save",
			job:     missingModuleJob,
			inTx:    false,
			want:    0,
			wantErr: true,
		},
		{
			name:    "missingClientJob_update_regular_save",
			job:     missingClientJob,
			inTx:    false,
			want:    0,
			wantErr: true,
		},
		{
			name:    "backupJob_update_transaction_save",
			job:     backupJob,
			inTx:    true,
			want:    backupJobID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int64
			var tx *sql.Tx
			if tt.inTx {
				tx, _ = db.Write().Begin()
				defer db.Unlock()
				got, err = tt.job.Update(tx)
			} else {
				got, err = tt.job.Update(nil)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				if err != nil && tt.inTx {
					tx.Rollback()
				}
				return
			}

			if tt.inTx {
				tx.Commit()
			}

			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReliqueJob_Save(t *testing.T) {
	SetupTest(t)

	newModule := module.Module{
		ModuleType:        "jobSaveNewModule",
		Name:              "jobSaveNewModule",
		BackupType:        backup_type.BackupType{Type: backup_type.Full},
		Schedules:         nil,
		ScheduleNames:     nil,
		BackupPaths:       []string{"/path1", "/path2"},
		PreBackupScript:   "",
		PostBackupScript:  "",
		PreRestoreScript:  "",
		PostRestoreScript: "",
	}
	newClient := client.Client{
		Name:          "jobSaveNewClient",
		Address:       "jobSaveNewClient",
		Port:          8434,
		Modules:       nil,
		Version:       "version",
		ServerAddress: "relique-server",
		ServerPort:    8433,
	}

	missingModuleJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Backup})
	missingModuleJob.Module = module.Module{}

	missingClientJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Backup})
	missingClientJob.Client = &client.Client{}

	backupJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Backup})
	existingBackupJob := New(&newClient, newModule, job_type.JobType{Type: job_type.Backup})
	if _, err := existingBackupJob.Save(); err != nil {
		t.Error("Cannot save existingBackupJob during test setup")
	}

	restoreJob := ReliqueJob{
		Uuid: "673379a1-5d5d-4c87-bbff-65b3e733ab0a",
		Client: &client.Client{
			Name:    "local",
			Address: "localhost",
			Port:    8434,
			Modules: []module.Module{
				module.Module{
					ID:                0,
					ModuleType:        "relique",
					Name:              "relique-diff",
					BackupType:        backup_type.BackupType{Type: 1},
					Schedules:         []schedule.Schedule{},
					BackupPaths:       []string{"/var/lib/relique", "/var/log/relique"},
					PreBackupScript:   "/var/lib/relique/modules/relique/scripts/prebackup.sh",
					PostBackupScript:  "/var/lib/relique/modules/relique/scripts/postbackup.sh",
					PreRestoreScript:  "/var/lib/relique/modules/relique/scripts/prerestore.sh",
					PostRestoreScript: "/var/lib/relique/modules/relique/scripts/postrestore.sh",
				},
			},
			Version:       "c5a89efe-8217-4166-a9ed-39830d033a47",
			ServerAddress: "localhost",
			ServerPort:    8433,
		},
		Module: module.Module{
			ModuleType:        "relique",
			Name:              "ondemand-relique-restore",
			BackupType:        backup_type.BackupType{Type: 3},
			BackupPaths:       []string{"/var/lib/relique", "/var/log/relique", "/etc/relique"},
			PreBackupScript:   "/var/lib/relique/modules/relique/scripts/prebackup.sh",
			PostBackupScript:  "/var/lib/relique/modules/relique/scripts/postbackup.sh",
			PreRestoreScript:  "/var/lib/relique/modules/relique/scripts/prerestore.sh",
			PostRestoreScript: "/var/lib/relique/modules/relique/scripts/postrestore.sh",
		},
		Status:             job_status.JobStatus{Status: 1},
		Done:               false,
		BackupType:         backup_type.BackupType{Type: 3},
		JobType:            job_type.JobType{Type: 2},
		RestoreJobUuid:     "77f32aba-ed9e-4d83-9c4c-396f33c3305f",
		RestoreDestination: "/tmp/relique_dest",
	}

	tests := []struct {
		name    string
		job     ReliqueJob
		want    int64
		wantErr bool
	}{
		{
			name:    "backup_job_save_ok",
			job:     backupJob,
			want:    2,
			wantErr: false,
		},
		{
			name:    "restore_job_save_ok",
			job:     restoreJob,
			want:    3,
			wantErr: false,
		},
		// TODO: Enable after update test
		{
			name:    "save_existing_job",
			job:     existingBackupJob,
			want:    1,
			wantErr: false,
		},
		{
			name:    "save_missing_module",
			job:     missingModuleJob,
			want:    0,
			wantErr: true,
		},
		{
			name:    "save_missing_client",
			job:     missingClientJob,
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.job.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.job.Module.ID == 0 {
				t.Errorf("Save() error, got zero ID for module. Module should be saved and have a DB ID")
				return
			}
			if !tt.wantErr && tt.job.Client.ID == 0 {
				t.Errorf("Save() error, got zero ID for module. Module should be saved and have a DB ID")
				return
			}

			if got != tt.want {
				t.Errorf("Save() got = %v, want %v", got, tt.want)
			}
		})
	}
}
