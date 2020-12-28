package custom_errors

import (
	"fmt"
	"testing"
)

func TestDBNotFoundError_Error(t *testing.T) {
	type fields struct {
		ID       int64
		ItemType string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string",
			fields: fields{
				ID:       10,
				ItemType: "test",
			},
			want: "no item 'test' with ID '10' found in DB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DBNotFoundError{
				ID:       tt.fields.ID,
				ItemType: tt.fields.ItemType,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDBNotFoundError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "error_not_found",
			args: args{e: &DBNotFoundError{
				ID:       2,
				ItemType: "test",
			}},
			want: true,
		},
		{
			name: "other_error",
			args: args{e: fmt.Errorf("other error")},
			want: false,
		},
		{
			name: "nil error",
			args: args{e: nil},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDBNotFoundError(tt.args.e); got != tt.want {
				t.Errorf("IsDBNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}
