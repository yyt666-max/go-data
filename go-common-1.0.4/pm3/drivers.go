package pm3

import (
	"fmt"
	"github.com/eolinker/eosc"
)

var (
	drivers = eosc.BuildUntyped[string, Driver]()
)

func Register(name string, driver Driver) {
	drivers.Set(name, driver)
}

func List() []Driver {
	return drivers.List()
}
func All() []string {
	return drivers.Keys()
}
func Create(plugins ...string) []IPlugin {

	notFound := make([]string, 0, len(plugins))
	pl := make([]IPlugin, 0, len(plugins))
	dl := make(map[string]Driver, len(plugins))

	for _, n := range plugins {
		d, h := drivers.Get(n)
		if !h {
			notFound = append(notFound, n)
			continue
		}
		dl[n] = d

	}
	if len(notFound) > 0 {
		panic(fmt.Sprintf("not found plugin:%v", notFound))
	}
	for n, d := range dl {
		p, err := d.Create()
		if err != nil {
			panic(fmt.Sprintf("create plugin [%s]", n))
			return nil
		}
		pl = append(pl, p)
	}

	return pl
}
