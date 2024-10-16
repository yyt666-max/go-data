package register

import (
	"testing"
)

type typeInterface interface {
	sayHello()
}
type testType struct {
}

func (t *testType) sayHello() {
	//fmt.Println("hello")
}

type args[T any] struct {
	v T
}
type testCase[T any] struct {
	name     string
	args     args[T]
	handlers []func(v T)
}

func TestRunHandler(t *testing.T) {

	tests := []testCase[typeInterface]{
		// TODO: Add test cases.
		{
			name: "test",
			args: args[typeInterface]{
				v: typeInterface(new(testType)),
			},
			handlers: []func(v typeInterface){
				func(v typeInterface) {
					v.sayHello()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.args.v
			for _, h := range tt.handlers {
				Handle(h)
			}

			Call(v)
		})
	}
}

func TestStructRunHandler(t *testing.T) {

	tests := []testCase[*testType]{
		// TODO: Add test cases.
		{
			name: "test",
			args: args[*testType]{
				v: new(testType),
			},
			handlers: []func(v *testType){
				func(v *testType) {
					v.sayHello()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.args.v
			for _, h := range tt.handlers {
				Handle(h)
			}

			Call(v)
		})
	}
}
func TestStructInterfaceRunHandler(t *testing.T) {

	cases := new(testType)
	isCall := 0
	Handle(func(v typeInterface) {
		isCall = 1
	})
	Call(cases)

	if isCall == 1 {
		t.Error("want not call")
	} else {
		t.Log("ok")
	}
	isCall = 0
	Handle(func(v typeInterface) {
		isCall = 2
	})
	Call[typeInterface](cases)

	if isCall == 2 {
		t.Log("ok")
	} else {
		t.Error("want to call")
	}

}
