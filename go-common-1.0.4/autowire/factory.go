package autowire

import (
	"log"
	"reflect"
)

func Auto[T any](handler func() reflect.Value) {
	name := TypeName[T]()

	locker.Lock()
	defer locker.Unlock()

	factory[name] = &factoryItem{
		fh: handler,
		et: elemType(reflect.TypeOf(new(T))),
	}
}
func Auto2[T any, S any]() {
	Auto[T](func() reflect.Value {
		return reflect.ValueOf(new(S))
	})
}
func doFactory() {
	locker.Lock()
	list := factory
	factory = make(map[string]*factoryItem)

	locker.Unlock()

	isContinue := true
	for isContinue {
		isContinue = false
		for n, i := range list {
			isContinue = isContinue || tryCreate(n, i)
		}
	}

}

type factoryItem struct {
	fh func() reflect.Value
	et reflect.Type
}

func tryCreate(name string, item *factoryItem) bool {
	if _, has := beans[name]; has {
		return false
	}

	if _, has := requires[name]; !has {
		return false
	}

	value := item.fh()
	if elemType(item.et).Kind() == reflect.Interface {
		if !value.Type().Implements(item.et) {

			log.Fatalf("not support inject type %s as %s", reflectTypeName(value.Type()), reflectTypeName(item.et))
		}
	}

	autowired(item.et, value)
	return true
}
