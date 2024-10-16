package autowire

import (
	"fmt"
	"testing"
)

type TestInterfaceBean interface {
	Name() string
}

type TestBean struct {
	name string
}

func (t *TestBean) Name() string {
	return t.name
}

type TestAutowiredStruct struct {
	Test TestInterfaceBean `autowired:""`

	Test2 *TestBean `autowired:""`

	Name *string `autowired:"name"`
}

func TestName(t *testing.T) {

	var s = new(TestAutowiredStruct)
	Autowired(s)
	var in = new(TestBean)
	in.name = "test1"
	Autowired(in)
	a := "testName"
	Inject(&a, "name")
	var interfaceBean TestInterfaceBean
	Autowired[TestInterfaceBean](&TestBean{name: "test2"})

	Autowired(&interfaceBean)

	var emptyIntefacce TestInterfaceBean
	Autowired(&emptyIntefacce)
	err := CheckComplete()
	if err != nil {
		return
	}

	fmt.Println("empty:", emptyIntefacce.Name())
	fmt.Println("name:", s.Test.Name())
	fmt.Println("name2:", s.Test2.Name())
	fmt.Println("nameString:", *s.Name)

}
