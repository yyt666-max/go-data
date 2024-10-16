package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server interface {
	http.Handler
	Permits() map[string][]string
}

var (
	_ Server = (*server)(nil)
)

type server struct {
	*gin.Engine
	permits map[string][]string
}

func (s *server) Permits() map[string][]string {
	return s.permits
}
