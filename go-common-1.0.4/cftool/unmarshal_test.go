package cftool

import (
	"encoding/json"
	"github.com/eolinker/go-common/autowire"
	"testing"
)

func Test_configYamlUnmarshaler_UnmarshalYAML(t *testing.T) {

	type TestConfig struct {
		UserName string `yaml:"user_name"`
		Password string `yaml:"password"`
		IP       string `yaml:"ip"`
		Port     int    `yaml:"port"`
		DB       string `yaml:"db"`
	}

	yamlData := []byte(`
mysql:
  user_name: "root"
  password: "asdfqwer"
  ip: "172.23.112.55"
  port: 3306
  db: "apserver"
redis: ""`)

	Register[TestConfig]("mysql")
	var test *TestConfig
	autowire.Autowired(&test)

	if err := unmarshalConfig(yamlData, "mysql"); err != nil {
		t.Errorf("UnmarshalYAML() error = %v", err)
		return
	}
	data, err := json.Marshal(test)
	if err != nil {
		return
	}
	t.Log("unmarshal:", string(data))

}
