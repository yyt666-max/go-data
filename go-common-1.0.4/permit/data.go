package permit

import (
	"net/http"
	"strings"
	"sync"
)

var (
	locker       sync.RWMutex
	permitDefine = make(map[string]map[string][]string)
)

type Paths map[string][]string

func AddPermitRule(access string, paths ...string) {
	group, _ := ReadAccessKey(access)
	locker.Lock()
	defer locker.Unlock()
	for _, p := range paths {
		method, path := readPath(p)
		path = FormatPath(method, path)
		if _, ok := permitDefine[path]; !ok {
			permitDefine[path] = make(map[string][]string)
		}
		permitDefine[path][group] = append(permitDefine[path][group], access)
	}
}
func GetPathRule(method string, path string) (map[string][]string, bool) {
	locker.RLock()
	defer locker.RUnlock()
	m, has := permitDefine[FormatPath(method, path)]
	return m, has
}
func readPath(path string) (string, string) {
	ps := strings.SplitN(path, ":", 2)
	if len(ps) == 1 {
		return http.MethodGet, ps[0]
	}
	return strings.ToUpper(ps[0]), ps[1]
}
func FormatPath(method string, path string) string {
	return strings.ToUpper(method) + ":" + "/" + strings.TrimLeft(path, "/")
}
func ReadAccessKey(access string) (string, string) {
	vs := strings.SplitN(access, ".", 2)
	if len(vs) != 2 {
		return "unknown", access
	}
	return vs[0], vs[1]
}
func FormatAccess(group, access string) string {
	return group + "." + access

}

func All() map[string]map[string][]string {
	return permitDefine
}
