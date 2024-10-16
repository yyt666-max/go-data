package store_mysql

import (
	"fmt"
)

type DBConfig struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Db       string `yaml:"db"`
}

func (c *DBConfig) getDBNS() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.Ip, c.Port, c.Db)
}
