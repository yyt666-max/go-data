package auto

import (
	"reflect"
	"strings"
)

const (
	TagAutoLabel = "aolabel"
)

func createLabelHandler[T any](v T) map[string]labelList {

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		if !rv.IsValid() {
			return nil
		}
		if rv.IsNil() {
			return nil
		}
	}
	out := make(map[string]labelList)
	foreach("", rv, out)
	return out
}

var operaterStructType = reflect.TypeOf((*Label)(nil)).Elem()

func readLabelName(label string) string {
	if strings.Index(label, ",") == -1 {
		return label
	}
	return strings.Split(label, ",")[0]
}

func foreach(label string, rv reflect.Value, out map[string]labelList) {
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		if !rv.IsValid() {
			return
		}
		if rv.IsNil() {
			return
		}
		foreach(label, rv.Elem(), out)
		return
	}
	if rv.Type() == operaterStructType {
		if label != "" {
			out[label] = append(out[label], rv.Addr().Interface().(*Label))
		}
		return
	}

	switch rv.Kind() {
	case reflect.Struct:

		num := rv.NumField()
		rt := rv.Type()
		for i := 0; i < num; i++ {
			fieldValue := rv.Field(i)
			fieldType := rt.Field(i)
			if fieldType.Anonymous {
				foreach("", fieldValue, out)
				continue
			}

			tag, has := fieldType.Tag.Lookup(TagAutoLabel)
			if has {
				tag = readLabelName(tag)
			}
			foreach(tag, fieldValue, out)

		}

	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			foreach(label, rv.Index(i), out)
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			foreach(label, rv.MapIndex(key), out)
		}
	default:
		return
	}

}
