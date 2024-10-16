package cftool

import "gopkg.in/yaml.v3"

var (
	fields = make(map[string][]*tFieldInfo)
	datas  = make(map[string]*yaml.Node)
)

type tFieldInfo struct {
	typeName string
	target   interface{}
}
