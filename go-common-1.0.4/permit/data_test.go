package permit

import (
	"net/http"
	"testing"
)

func Test_readPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "get",
			args: args{
				path: "GET:/api/v1/test",
			},
			want:  http.MethodGet,
			want1: "/api/v1/test",
		},
		{
			name: "post",
			args: args{
				path: "POST:/api/v1/test",
			},
			want:  http.MethodPost,
			want1: "/api/v1/test",
		},
		{
			name: "put",
			args: args{
				path: "PUT:/api/v1/test",
			},
			want:  http.MethodPut,
			want1: "/api/v1/test",
		},
		{
			name: "delete",
			args: args{
				path: "DELETE:/api/v1/test",
			},
			want:  http.MethodDelete,
			want1: "/api/v1/test",
		},
		{
			name: "default",
			args: args{
				path: "/api/v1/test",
			},
			want:  http.MethodGet,
			want1: "/api/v1/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readPath(tt.args.path)
			if got != tt.want {
				t.Errorf("readPath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
