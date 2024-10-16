package pm3

import "github.com/gin-gonic/gin"

type formApi struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

func (f *formApi) Method() string {
	return f.method
}

func (f *formApi) Path() string {
	return f.path
}

func (f *formApi) Handler(context *gin.Context) {
	f.handler(context)
}
func Gen(a *apiDoc) Api {
	return &formApi{
		method:  a.Method,
		path:    a.Path,
		handler: a.Handler(),
	}
}
