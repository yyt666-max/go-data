package pm3

import "github.com/gin-gonic/gin"

type Api interface {
	Method() string
	Path() string
	Handler(*gin.Context)
}
