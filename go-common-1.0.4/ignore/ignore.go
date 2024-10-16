package ignore

import (
	"github.com/eolinker/go-common/utils"
	"sync"
)

var (
	rwlock      = sync.RWMutex{}
	ignorePaths = make(map[string]map[string]utils.Set[string])
)

func IgnorePath(name, method, path string) {
	rwlock.Lock()
	defer rwlock.Unlock()
	if ignorePaths == nil {
		ignorePaths = make(map[string]map[string]utils.Set[string])
	}
	if ignorePaths[name] == nil {
		ignorePaths[name] = make(map[string]utils.Set[string])
	}
	if ignorePaths[name][method] == nil {
		ignorePaths[name][method] = utils.NewSet(path)
	} else {
		ignorePaths[name][method].Set(path)
	}

}

func IsIgnorePath(name, method, path string) bool {
	rwlock.RLock()
	defer rwlock.RUnlock()
	if isIgnorePath(name, method, path) {
		return true
	}
	if method != "*" {
		return isIgnorePath(name, "*", path)
	}
	return false
}

func isIgnorePath(name, method, path string) bool {

	if ignorePaths == nil {
		return false
	}
	if ignorePaths[name] == nil {
		return false
	}
	if ignorePaths[name][method] == nil {
		return false
	}
	return ignorePaths[name][method].Has(path)
}
