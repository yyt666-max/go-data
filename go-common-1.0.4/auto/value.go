package auto

import (
	"reflect"
	"strings"
)

const (
	TagAutoValue = "aovalue"
)

func Auto(tagValue string, value string, target any) {
	autoSetValue(strings.ToUpper(tagValue), value, reflect.ValueOf(target))
}
func autoSetValue(tagValue string, value string, rv reflect.Value) {
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		if !rv.IsValid() {
			return
		}
		if rv.IsNil() {
			return
		}
		autoSetValue(tagValue, value, rv.Elem())
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
				autoSetValue(tagValue, value, fieldValue)
				continue
			}
			autoValue := fieldType.Tag.Get(TagAutoValue)
			if strings.HasPrefix(strings.ToUpper(autoValue), tagValue) {
				for fieldValue.Kind() == reflect.Ptr {
					fieldValue = fieldValue.Elem()
				}

				if fieldValue.Kind() == reflect.String {
					fieldValue.SetString(value)
				}
				continue
			}

			switch fieldValue.Kind() {

			case reflect.Struct, reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Array, reflect.Map:
				autoSetValue(tagValue, value, fieldValue)
			default:
				continue
			}

		}
	default:
		return
	}
}
