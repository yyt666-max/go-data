package autowire

import (
	"fmt"
	"reflect"
)

func typeName[T any]() string {
	t := reflect.TypeOf(new(T)).Elem()

	return reflectTypeName(t)

}
func reflectTypeName(t reflect.Type) string {

	p, n := reflectName(t)
	return fmt.Sprintf("%s.%s", p, n)
}
func TypeName[T any](v ...T) string {
	return typeName[T]()
}
func reflectName(t reflect.Type) (path string, name string) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath(), t.Name()
}

func elemType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func elemValue(value reflect.Value) reflect.Value {
	if value.Kind() != reflect.Ptr {
		return value
	}
	e := value.Elem()
	if e.Kind() == reflect.Ptr {
		return elemValue(e)
	}
	return value
}
