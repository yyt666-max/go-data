package autowire

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

var (
	beans    = make(map[string]reflect.Value)
	requires = make(map[string][]*requireValue) // map[requireType]map[name][]*value

	factory = make(map[string]*factoryItem)

	objects      = make([]reflect.Value, 0)
	requireCount = make(map[int]int)
	empty        = make(map[string][]reflect.Value)
	locker       sync.Mutex
)

type requireValue struct {
	value      reflect.Value
	isExported bool
	id         int
	name       string
}

// Autowired 自动注入依赖
// v 注入目标, 必须是指针
// 目标为struct, 会遍历其成员,并根据注解 autowired 进行注入,并且按类型名将对象注入给其依赖性
// 如果输入为interface类型,则会将目标按类型名注入给其他依赖项
// 如果输入是基础类型, 则必须有name,并按name进行处理
//
// names 为可选的对象别名,
//
//	如果不传,则用识别到的类型名进行存储与注入, 如果类型名的bean已存在,则报冲突panic
//	如果names为空,则会忽略类型名的冲突,并使用所有的名称进行注入, 如names存在冲突,一样会报冲突panic
func Autowired[T any](v T, names ...string) {

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr && rv.Kind() != reflect.Interface {
		log.Fatal("only autowired ptr but get:", rv.Kind().String())
	}
	if rv.IsNil() {
		log.Fatal("not allow autowired nil ptr for:", rv.String())
	}

	ert := elemType(reflect.TypeOf(new(T)))

	if ert.PkgPath() == "" {
		log.Fatal("not support autowired type:", ert.String())
	}
	autowired(ert, rv, names...)

}

//
//func Get[T any](v *T, name ...string) bool {
//	locker.Lock()
//	defer locker.Unlock()
//
//	rv := reflect.ValueOf(v)
//	tp := rv.Type()
//	for _, n := range name {
//		value, has := beans[n]
//		if !has {
//			continue
//		}
//
//		if rv.IsNil() {
//			rv.Set(reflect.New(rv.Type()))
//		}
//		rv.Elem().Set(value)
//		return true
//	}
//
//	value, has := beans[fmt.Sprintf("%s.%s", tp.PkgPath(), tp.User())]
//	if !has {
//		return false
//	}
//	if rv.IsNil() {
//		rv.Set(reflect.New(rv.Type()))
//	}
//	rv.Elem().Set(value)
//	return true
//
//}

func autowired(et reflect.Type, v reflect.Value, names ...string) {
	locker.Lock()
	complete, needInject := doAutowired(et, v)
	locker.Unlock()
	if complete {
		v.Interface().(Complete).OnComplete()
	}
	if needInject {
		inject(et, v, names...)
	}

}
func doAutowired(et reflect.Type, v reflect.Value) (bool, bool) {

	if v.Kind() == reflect.Interface {
		return false, true
	}

	pkgName, name := et.PkgPath(), et.Name()
	beanName := fmt.Sprintf("%s.%s", pkgName, name)

	v = setEmpty(v)
	if v.IsNil() {
		// 如果传入值为nil, 则整体注入
		if pkgName == "" { // 不支持匿名注入
			panic("not support type:" + name)
		}

		bv, has := beans[beanName]
		if has {
			// 存在则整体注入
			v.Set(bv)

		} else {
			// 不存在则缓存起来,到check时再按内部字段注入
			empty[beanName] = append(empty[beanName], v)
		}

		return false, false
	}

	if v.Kind() == reflect.Interface {
		// interface 直接执行下一步 inject
		return false, true
	}

	return autowiredRoot(v), true
}
func autoWiredSetField(id int, root string, path []string, v reflect.Value, field reflect.StructField) {

	path = append(path, field.Name)
	beanName, has := field.Tag.Lookup("autowired")

	if !has {

		// 只对非指针类型的子字段执行内部注入
		if v.Kind() == reflect.Struct {
			t := elemType(v.Type())

			for i := 0; i < t.NumField(); i++ {
				autoWiredSetField(id, root, path, v.Field(i), t.Field(i))
			}
		}
		return
	}
	if beanName == "" {
		t := elemType(field.Type)
		if t.PkgPath() == "" {
			panic(fmt.Sprintf("anonymous autowired not support field [%s:%s] type [%s] ", root, strings.Join(path, "."), t.Name()))
		}
		beanName = fmt.Sprint(t.PkgPath(), ".", t.Name())
	}
	value, has := beans[beanName]
	if has {
		setValue(v, value, field.IsExported())
		return
	}

	requires[beanName] = append(requires[beanName], &requireValue{
		value:      v,
		id:         id,
		name:       root,
		isExported: field.IsExported(),
	})
	requireCount[id] += 1

}
func setFieldValue(field, value reflect.Value, exported bool) {
	if exported {

		field.Set(value)

	} else {
		//nolint:gosec
		reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
			Elem().
			Set(value)
	}

}
func setValue(target, source reflect.Value, exported bool) {
	target = setEmpty(target)
	tt := target.Type()
	ett := elemType(tt)
	switch tt.Kind() {
	case reflect.Ptr:
		if source.Kind() == reflect.Ptr {
			setFieldValue(target, source, exported)
			return
		}
		setFieldValue(target, reflect.New(ett), exported)
		target.Elem().Set(source)
		return
	case reflect.Interface:
		st := source.Type()
		//if st.Kind() == reflect.Ptr {
		/// bean 必须是指针
		for !st.Implements(tt) && st.Kind() == reflect.Ptr {
			st = st.Elem()
			source = source.Elem()
		}
		setFieldValue(target, source, exported)
		return
		//}
	default:

	}

	if tt.Kind() == reflect.Interface {
		st := source.Type()
		for !st.Implements(tt) {
			st = st.Elem()
			source = source.Elem()
		}
		setFieldValue(target, source, exported)
		return
	}
	for source.Kind() == reflect.Ptr {
		source = source.Elem()
	}
	setFieldValue(target, source, exported)

}

func autowiredRoot(v reflect.Value) bool {

	id := len(objects)
	objects = append(objects, v)
	v = elemValue(v)
	t := elemType(v.Type())
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("just support inject struct but get :%s of %s.%s", t.Kind().String(), t.PkgPath(), t.Name()))
	}
	root := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	num := t.NumField()
	vStruct := v.Elem()
	for i := 0; i < num; i++ {
		f := t.Field(i)
		autoWiredSetField(id, root, nil, vStruct.Field(i), f)
	}

	if requireCount[id] == 0 {
		iv := v.Interface()
		if _, ok := iv.(Complete); ok {
			return true

		}
	}
	return false
}

func Inject[T any](v T, name ...string) {

	et := elemType(reflect.TypeOf(new(T)))
	inject(et, reflect.ValueOf(v), name...)

}
func inject(et reflect.Type, value reflect.Value, name ...string) {

	locker.Lock()
	completeList, err := doInject(et, value, name...)

	locker.Unlock()
	if err != nil {
		panic(err)
	}
	completeList.OnComplete()

}
func doInject(et reflect.Type, value reflect.Value, name ...string) (completeList Complete, err error) {
	if value.Kind() != reflect.Ptr && value.Kind() != reflect.Interface {
		vp := reflect.New(value.Type())
		vp.Elem().Set(value)
		value = vp
	}
	if value.IsNil() {
		return nil, nil
	}

	rv := elemValue(value)
	pn, _ := reflectName(et)
	if pn != "" && len(name) == 0 {
		return addBean(reflectTypeName(et), rv)
	}

	return addBean(name[0], rv)

}
func addBean(name string, value reflect.Value) (completerList CompleteList, err error) {
	if _, has := beans[name]; has {

		return nil, fmt.Errorf(" [%s] conflicts with existing", name)

	}
	beans[name] = value

	for _, vs := range empty[name] {
		setValue(vs, value, true)

	}
	delete(empty, name)

	for _, rv := range requires[name] {
		setValue(rv.value, value, rv.isExported)
		requireCount[rv.id]--

		if requireCount[rv.id] == 0 {
			v := objects[rv.id]
			if h, ok := v.Interface().(Complete); ok {

				completerList = append(completerList, h)
			}
			requireCount[rv.id] = -1
		}
	}
	delete(requires, name)
	return
}

func setEmpty(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}
	t := v.Type()
	t = t.Elem()

	if t.Kind() != reflect.Ptr && t.Kind() != reflect.Interface {
		return v
	}
	if v.IsNil() {
		if v.Kind() == reflect.Ptr {
			v.Set(reflect.New(v.Type()))
		}
	}
	return setEmpty(v.Elem())
}
