package autowire

import (
	"fmt"
	"reflect"
	"testing"
)

type TestInterface interface {
	isComplete() string
	OnComplete()
}
type TestStruct struct {
	complete bool
}

func (t *TestStruct) OnComplete() {
	t.complete = true
}

func (t *TestStruct) isComplete() string {
	if t.complete {
		return "complete"
	} else {
		return "not complete"
	}
}

type TestNameStruct struct {
	TestInterface `autowired:""`
	Name          string        `autowired:"name"`
	test          TestInterface `autowired:""`
}

func (t *TestNameStruct) TestInterfaceInject() {
	fmt.Println("autowired TestInterfaceInject ok:", t.Name)
}

type TestInterfaceInject interface {
	TestInterfaceInject()
}

func TestAutowired2(t *testing.T) {
	type InjectCase func()

	injects := []InjectCase{func() {
		Inject("vv", "rt")
	}, func() {
		var name = "name inject"
		Inject(&name, "name")
		st := new(TestStruct)
		//Autowired(st)
		Inject[TestInterface](st)
	}}
	for _, h := range injects {
		h()
	}

	var v2 *TestNameStruct
	Autowired(&v2)

	var testInterfaceInject TestInterfaceInject

	Autowired(&testInterfaceInject)
	var v = new(TestNameStruct)
	Autowired[TestInterfaceInject](v)
	//Inject(v)
	err := CheckComplete()
	if err != nil {
		t.Error(err)
		return
	}
	if testInterfaceInject == nil {
		t.Error("not autowired interface")
		return
	}

	testInterfaceInject.TestInterfaceInject()
	if testInterfaceInject == nil {
		t.Error("not autowired target")
		return
	}
	t.Log("isComplete:", v2.isComplete())
	if v2.Name == "" {
		t.Error("not autowired member name:", v2.Name)
	}
	if v2.test == nil {
		t.Error("not autowired member:test")
	}
	if v2.TestInterface == nil {
		t.Error("not autowired member:TestInterface")
	}

	t.Log(".test=>", v2.test.isComplete())
	t.Log(".TestInterface=>", v.isComplete())
	t.Log("name=", v2.Name)
	t.Log("ok")

}

type TTT int

func (T TTT) Test() {
	//TODO implement me
	panic("implement me")
}

type TTI interface {
	Test()
}

func TestCheckr(t *testing.T) {
	//var i = new(int)
	var i TTI = new(TTT)
	v := reflect.ValueOf(&i)
	et := elemType(v.Type())
	fmt.Println(et.String(), ":", et.Kind())
	//v := reflect.ValueOf(i)
	//nv := reflect.New(v.Type().Elem())
	//fmt.Println(nv.Type().String())
	//fmt.Println(v.Type().String())
	//v.Set(nv)

}

//func TestGet(t *testing.T) {
//	name := "name-test"
//	Inject(name, "name")
//	var r string
//	Get(&r, "name")
//	if r == name {
//		t.Log("ok")
//		return
//	}
//	t.Error("want:", name, " bug get:", r)
//}
