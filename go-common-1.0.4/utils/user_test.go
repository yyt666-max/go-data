package utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestUserId(t *testing.T) {

	ctx := gin.CreateTestContextOnly(nil, gin.Default())

	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{
			name: "ginCtx",
			ctx:  SetUserId(ctx, "test_USerID"),
			want: "test_USerID",
		},
		{
			name: "widthValue",
			ctx:  context.WithValue(ctx, "TestForWithValue", "test_USerID"),
			want: "test_USerID",
		},
		{
			name: "test_USerID",
			ctx:  context.Background(),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserId(tt.ctx); got != tt.want {
				t.Errorf("UserId() = %v, want %v", got, tt.want)
			}
		})
	}
}
