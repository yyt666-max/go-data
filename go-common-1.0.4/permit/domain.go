package permit

import (
	"github.com/gin-gonic/gin"
)

type DomainHandler func(c *gin.Context) ([]string, []string, bool)

var (
	//domainHandlers []*domainHandlerItem
	domainHandlers = make(map[string]DomainHandler)
)

func AddDomainHandler(group string, handler DomainHandler) {
	locker.Lock()
	defer locker.Unlock()

	if _, has := domainHandlers[group]; has {
		panic("domain handler already exists:" + group)
	}

	domainHandlers[group] = handler

	//domainHandlers = append(domainHandlers, &domainHandlerItem{
	//	prefix:  prefix,
	//	handler: handler,
	//})
	//
	//utils.Sort(domainHandlers, func(i, j *domainHandlerItem) bool {
	//	if len(i.prefix) != len(j.prefix) {
	//		return len(i.prefix) < len(j.prefix)
	//	}
	//	return i.prefix < j.prefix
	//})
}
func SelectDomain(group string) (DomainHandler, bool) {

	locker.RLock()
	defer locker.RUnlock()
	h, has := domainHandlers[group]
	return h, has
	//for _, item := range domainHandlers {
	//	if strings.HasPrefix(path, item.prefix) {
	//		return item.handler, true
	//	}
	//}
	//return nil, false
}
