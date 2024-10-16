package pm3

import "github.com/gin-gonic/gin"

type simpleMiddleware struct {
	checkFunc func(method, path string) bool
	handler   gin.HandlerFunc
	sort      int
}

func (s *simpleMiddleware) Sort() int {
	return s.sort
}

func (s *simpleMiddleware) Check(method string, path string) (bool, []gin.HandlerFunc) {
	if s.handler == nil {
		return false, nil
	}
	if s.checkFunc == nil || s.checkFunc(method, path) {
		return true, []gin.HandlerFunc{s.handler}

	}
	return false, nil
}

func (s *simpleMiddleware) Handler(ginCtx *gin.Context) {
	s.handler(ginCtx)
}

func CreateMiddle(checkHandler func(method string, path string) bool, handler gin.HandlerFunc, sort int) IMiddleware {
	return &simpleMiddleware{
		checkFunc: checkHandler,
		handler:   handler,
		sort:      sort,
	}
}
