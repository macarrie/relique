package repo

import (
	"reflect"
	"testing"
	"os"
)

func TestRepoLocalNew(t *testing.T) {
	type args struct {
		Name    string
		Path    string
		Default bool
	}
	tests := []struct {
		name    string
		args    args
		want    RepositoryLocal
	}{
		{
			name: "new_local_repo",
			args: args{
				Name:    "test_local_repo",
				Path:    "/tmp/test_local_repo",
				Default: true,
			},
			want: RepositoryLocal{
				Name:    "test_local_repo",
				Type:    "local",
				Path:    "/tmp/test_local_repo",
				Default: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RepoLocalNew(tt.args.Name, tt.args.Path, tt.args.Default)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoLocal_GetLog(t *testing.T) {
	tests := []struct {
		name    string
		repo RepositoryLocal
		wantNil bool
	}{
		{
			name: "get_log",
			repo: RepositoryLocal{
				Name:    "test_repo",
				Type: "local",
				Path: "/tmp/test_repo",
				Default:false,
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.repo.GetLog(); (got == nil) != tt.wantNil {
				t.Errorf("GetLog() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func TestRepoLocal_Write(t *testing.T) {
	type args struct {
		Repo    RepositoryLocal
		Path    string
	}
	tests := []struct {
		name    string
		args    args
		want    RepositoryLocal
		wantErr bool
	}{
		{
			name: "write",
			args: args{
				Repo: RepositoryLocal{
					Name:    "test_local_repo",
					Type:    "local",
					Path:    "/tmp/test_local_repo",
					Default: true,
				},
			},
			want: RepositoryLocal{
				Name:    "test_local_repo",
				Type:    "local",
				Path:    "/tmp/test_local_repo",
				Default: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpTestFolder, err := os.MkdirTemp("", "relique-test-repolocal-write-*")
			defer os.RemoveAll(tmpTestFolder)
			if err != nil {
				t.Errorf("Write() cannot create test folder, error = %v", err)
				return
			}
			
			writeErr := tt.args.Repo.Write(tmpTestFolder)
			if (writeErr != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", writeErr, tt.wantErr)
				return
			}

			got, err := LoadFromPath(tmpTestFolder)
			if len(got) != 1 {
			if !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("Load after Write() got = %v, want %v", got, tt.want)
			}
			}
		})
	}
}