package autowire

import (
	"fmt"
	"strings"
)

type OnCreate interface {
	OnCreate()
}
type Inited interface {
	OnInit()
}
type Complete interface {
	OnComplete()
}
type CompleteList []Complete

func (cl CompleteList) OnComplete() {
	for _, c := range cl {
		c.OnComplete()
	}
}

func CheckComplete() error {
	doFactory()
	requiresBean := make([]string, 0)

	for name, ls := range empty {
		if len(ls) > 0 {
			requiresBean = append(requiresBean, name)
		}
	}
	if len(requiresBean) > 0 {
		return fmt.Errorf("bean not init:[%s]", strings.Join(requiresBean, ","))
	}

	for name, ls := range requires {

		if len(ls) > 0 {
			ds := make([]string, 0, len(ls))
			for _, v := range ls {
				ds = append(ds, v.name)
			}
			return fmt.Errorf("bean [%s] not init dependent on:[%s]", name, strings.Join(ds, ","))

		}
	}
	for _, b := range objects {
		if v, ok := b.Interface().(Inited); ok {
			v.OnInit()
		}
	}
	return nil
}
