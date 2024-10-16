package repo

import (
	"reflect"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    Repository
		wantErr bool
	}{
		{
			name: "example",
			args: args{file: "../../test/repo/example.toml"},
			want: &RepositoryLocal{
				Name:    "test_repo",
				Type:    "local",
				Path:    "/tmp/test_repo",
				Default: false,
			},
			wantErr: false,
		},
		{
			name:    "default_values",
			args:    args{file: "../../test/repo/empty.toml"},
			want:    &GenericRepository{},
			wantErr: true,
		},
		{
			name:    "not_found",
			args:    args{file: "../../test/repo/not_found.toml"},
			want:    &GenericRepository{},
			wantErr: true,
		},
		{
			name:    "invalid_toml",
			args:    args{file: "../../test/repo/parse_error.toml"},
			want:    &GenericRepository{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFromFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []Repository
		wantErr bool
	}{
		{
			name: "example",
			args: args{path: "../../test/repo/"},
			want: []Repository{&RepositoryLocal{
				Name:    "test_repo",
				Type:    "local",
				Path:    "/tmp/test_repo",
				Default: false,
			},
			},
			wantErr: false,
		},
		{
			name:    "path_does_not_exist",
			args:    args{path: "/does_not_exist"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFromPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadFromPath() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetByName(t *testing.T) {
	type args struct {
		list             []Repository
		searchModuleName string
	}
	repoList := []Repository{
		&GenericRepository{
			Name: "test1",
		},
		&GenericRepository{
			Name: "test2",
		},
		&GenericRepository{
			Name: "test3",
		},
	}
	tests := []struct {
		name    string
		args    args
		want    Repository
		wantErr bool
	}{
		{
			name: "repo_found",
			args: args{
				list:             repoList,
				searchModuleName: "test1",
			},
			want:    &GenericRepository{Name: "test1"},
			wantErr: false,
		},
		{
			name: "repo_not_found",
			args: args{
				list:             repoList,
				searchModuleName: "totodesbois",
			},
			want:    &GenericRepository{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetByName(tt.args.list, tt.args.searchModuleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefault(t *testing.T) {
	type args struct {
		list []Repository
	}
	repoList := []Repository{
		&GenericRepository{
			Name:    "test1",
			Default: true,
		},
		&GenericRepository{
			Name: "test2",
		},
		&GenericRepository{
			Name: "test3",
		},
	}
	tests := []struct {
		name    string
		args    args
		want    Repository
		wantErr bool
	}{
		{
			name: "default_found",
			args: args{
				list: repoList,
			},
			want:    &GenericRepository{
				Name: "test1",
				Default: true,
			},
			wantErr: false,
		},
		{
			name: "default_not_found",
			args: args{
				list: []Repository{},
			},
			want:    &GenericRepository{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDefault(tt.args.list)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefault() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDefault() got = %v, want %v", got, tt.want)
			}
		})
	}
}
