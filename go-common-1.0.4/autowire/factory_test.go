package autowire

import (
	"fmt"
	"testing"
)
import "reflect"

type TestFactoryInterface interface {
	TestInterface() string
}
type TestFactoryStruct struct {
	name   string
	action string
}

func (t *TestFactoryStruct) OnComplete() {
	t.action = "complete"
}

type TestFactoryAuto struct {
	Target TestFactoryInterface `autowired:""`
}

func (t *TestFactoryStruct) TestInterface() string {
	return fmt.Sprintf("%s:%s", t.name, t.action)
}

func TestAuto(t *testing.T) {
	Auto[TestFactoryInterface](func() reflect.Value {
		return reflect.ValueOf(&TestFactoryStruct{name: "Auto"})

	})
	v := new(TestFactoryAuto)
	Autowired(v)
	err := CheckComplete()
	if err != nil {
		t.Error(err)
		return
	}
	if v.Target == nil {
		t.Error("autowired nil")
		return
	}
	t.Logf("test fatory :%s", v.Target.TestInterface())
}
