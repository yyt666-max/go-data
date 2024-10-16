package cftool

import (
	"github.com/eolinker/go-common/autowire"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"sync"
)

var (
	locker sync.Mutex
)

func Register[T any](name string, vs ...T) {

	t := reflect.TypeOf(new(T))
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("config field must struct")
	}
	locker.Lock()
	defer locker.Unlock()

	fi := &tFieldInfo{
		typeName: autowire.TypeName[T](),
	}
	if len(vs) == 0 || reflect.ValueOf(vs[0]).Kind() != reflect.Pointer {
		fi.target = reflect.New(t).Interface()
	} else {
		fi.target = vs[0]
	}
	if node, has := datas[name]; has {
		if err := node.Decode(fi.target); err != nil {
			return
		}

		autowire.Inject(fi.target, fi.typeName)
		return
	}

	fields[name] = append(fields[name], fi)
}

func InitFor(name string, data []byte) {
	err := unmarshalConfig(data, name)
	if err != nil {
		panic(err)
	}
}
func ReadFile(path string) {

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = unmarshalConfig(data, path)
	if err != nil {
		panic(err)
	}
}
func unmarshalConfig(data []byte, name string) error {
	locker.Lock()
	defer locker.Unlock()

	c := newConfigYamlUnmarshaler(fields, name)
	fields = make(map[string][]*tFieldInfo)
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	for _, vs := range c.data {
		for _, v := range vs {
			autowire.Inject(v.target, v.typeName)
		}
	}

	return nil
}
