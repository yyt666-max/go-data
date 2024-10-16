package register

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	locker   sync.Mutex
	handlers = make(map[string]any)
	values   = make(map[string]any)
)

func Handle[T any](h func(v T)) {
	name := typeName[T]()
	if !locker.TryLock() {
		panic(fmt.Sprintf("Handler:%s deadlock", name))
	}

	defer locker.Unlock()
	v, has := values[name]
	if has {
		h(v.(T))
		return
	}
	list, has := handlers[name]
	if !has {
		list = make([]func(v T), 0)
		//handlers[name]=list
	}
	hl := list.([]func(v T))
	handlers[name] = append(hl, h)
}

func typeName[T any]() string {
	t := reflect.TypeOf(new(T)).Elem()
	pkg := t
	for pkg.Kind() == reflect.Ptr {
		pkg = pkg.Elem()
	}

	return fmt.Sprintf("%s.%s", pkg.PkgPath(), t.String())

}

func Call[T any](v T) {
	name := typeName[T]()
	if !locker.TryLock() {
		panic(fmt.Sprintf("Handler:%s deadlock", name))
	}

	defer locker.Unlock()
	values[name] = v
	list, has := handlers[name]
	if !has {
		return
	}
	delete(handlers, name)
	hl, ok := list.([]func(v T))
	if !ok {
		return
	}
	for _, h := range hl {
		h(v)
	}

}
