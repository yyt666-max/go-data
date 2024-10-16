package pm3

import (
	"github.com/eolinker/go-common/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Driver interface {
	Create() (IPlugin, error)
}
type AccessConfig interface {
	Access() map[string][]string
}
type IPlugin interface {
	Name() string
}

type IPluginApis interface {
	APis() []Api
}
type IPluginMiddleware interface {
	Middlewares() []IMiddleware
}

type IMiddleware interface {
	Check(method string, path string) (bool, []gin.HandlerFunc)
	Sort() int
	//Handler(ginCtx *gin.Context)
}
type MiddlewareList []IMiddleware

func (ml MiddlewareList) Sort() {
	utils.Sort(ml, func(a, b IMiddleware) bool {
		return a.Sort() < b.Sort()
	})
}
func (ml MiddlewareList) Check(method string, path string) []gin.HandlerFunc {
	rl := make([]gin.HandlerFunc, 0, len(ml))
	for _, m := range ml {
		if ok, h := m.Check(method, path); ok {
			rl = append(rl, h...)
		}
	}
	if len(rl) == 0 {
		return nil
	}
	return rl
}

type IFrontendFiles interface {
	Files() []FrontendFiles
}
type FrontendFiles struct {
	Path       string
	FileSystem http.FileSystem
}
type ContentHandler struct {
	Path    string
	Handler gin.HandlerFunc
}
