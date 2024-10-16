package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eolinker/go-common/pm3"
	"github.com/eolinker/go-common/register"
	"github.com/eolinker/go-common/utils"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	frontendTop     = make(map[string]struct{})
	notFoundContent = []byte("404 page not found")

	indexHtmlHandler gin.HandlerFunc

	systemPlugin []pm3.IPlugin
)

func AddSystemPlugin(ps ...pm3.IPlugin) {
	systemPlugin = append(systemPlugin, ps...)
}
func SetIndexHtmlHandler(indexHtml gin.HandlerFunc) {
	indexHtmlHandler = indexHtml
}

type ServiceBuilder interface {
	Build() Server
	Detail() *ServerDetail
}
type imlServiceBuilder struct {
	pluginAll []pm3.IPlugin
}

func (i *imlServiceBuilder) Build() Server {
	engine := gin.Default()
	s := &server{
		Engine:  engine,
		permits: make(map[string][]string),
	}

	handlers := make(map[string][]pm3.Api)
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.GET("/", indexHtmlHandler)

	middlewareList := make(pm3.MiddlewareList, 0, len(i.pluginAll))
	for _, p := range i.pluginAll {
		if ac, ok := p.(pm3.AccessConfig); ok {
			for k, r := range ac.Access() {
				s.permits[k] = append(s.permits[k], r...)
			}
		}
		if mh, ok := p.(pm3.IPluginMiddleware); ok {
			middlewareList = append(middlewareList, mh.Middlewares()...)
		}
	}
	middlewareList.Sort()
	for _, p := range i.pluginAll {

		if fh, ok := p.(pm3.IFrontendFiles); ok {

			for _, file := range fh.Files() {
				middlewareHandlers := middlewareList.Check(http.MethodGet, file.Path)
				engine.Group("/", middlewareHandlers...).StaticFS(file.Path, file.FileSystem)
				root := strings.Split(strings.Trim(file.Path, "/"), "/")[0]
				frontendTop[root] = struct{}{}

			}
		}
		if ai, ok := p.(pm3.IPluginApis); ok {
			for _, a := range ai.APis() {
				middlewareHandlers := middlewareList.Check(a.Method(), a.Path())

				handlers[p.Name()] = append(handlers[p.Name()], a)
				engine.Group("/", middlewareHandlers...).Handle(a.Method(), a.Path(), a.Handler)
			}
		}
	}
	for i := range handlers {
		utils.Sort(handlers[i], func(i, j pm3.Api) bool {
			if i.Path() == j.Path() {
				return i.Method() < j.Method()
			}
			return i.Path() < j.Path()
		})
	}
	handlerPath := make(map[string][]string)
	for k, v := range handlers {
		handlerPath[k] = utils.SliceToSlice(v, func(i pm3.Api) string {
			return fmt.Sprint(i.Method(), " ", i.Path())
		})

	}
	engine.GET("/_system/apis", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/plain; charset=utf-8")
		ctx.YAML(http.StatusOK, handlerPath)
	})
	engine.NoRoute(apiNoRoute, assetsNoRoute, indexHtmlHandler)
	register.Call[Server](s)
	return s
}

func newServiceBuilder(plugins ...string) ServiceBuilder {
	if len(plugins) == 0 {
		plugins = pm3.All()
	}
	pluginList := make([]pm3.IPlugin, 0, len(systemPlugin)+len(plugins))
	pluginList = append(pluginList, systemPlugin...)

	pluginList = append(pluginList, pm3.Create(plugins...)...)
	return &imlServiceBuilder{
		pluginAll: pluginList,
	}
}

func CreateServer(plugins ...string) ServiceBuilder {
	return newServiceBuilder(plugins...)

}

func assetsNoRoute(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	root := strings.FieldsFunc(uri, func(r rune) bool {
		return r == '/'
	})
	if _, has := frontendTop[root[0]]; has {
		ctx.Data(http.StatusNotFound, "text/html; charset=utf-8", notFoundContent)
		ctx.Abort()
	}

}
func apiNoRoute(ctx *gin.Context) {
	if strings.HasPrefix(ctx.Request.RequestURI, "/api") {

		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
	}

}
