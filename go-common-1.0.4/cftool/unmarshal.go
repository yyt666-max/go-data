package cftool

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type yamlUnmarshaler struct {
	data map[string][]*tFieldInfo
	name string
}

func newConfigYamlUnmarshaler(fields map[string][]*tFieldInfo, name string) *yamlUnmarshaler {
	return &yamlUnmarshaler{
		data: fields,
		name: fmt.Sprintf("root:%s", name),
	}
}

func (c *yamlUnmarshaler) ResolveValue(name string, v *yaml.Node) error {
	datas[name] = v
	targets, has := c.data[name]
	if !has {
		return nil
	}

	for _, target := range targets {

		err := v.Decode(target.target)

		if err != nil {
			return err
		}

	}
	return nil
}

func (c *yamlUnmarshaler) UnmarshalYAML(value *yaml.Node) error {
	err := c.ResolveValue(c.name, value)
	if err != nil {
		return err
	}

	for i, v := range value.Content {
		if v.Kind == yaml.ScalarNode {
			continue
		}
		if v.Kind == yaml.MappingNode || v.Kind == yaml.SequenceNode {
			name := value.Content[i-1].Value

			err := c.ResolveValue(name, v)
			if err != nil {
				return err
			}

		}
	}

	return nil
}
