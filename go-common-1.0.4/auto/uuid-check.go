package auto

import (
	"reflect"
)

const (
	TagCheck = "aocheck"
)

type IDCheckMap map[string]*IDCheck

type IDCheck struct {
	uuids []string
	name  string
}

func (i *IDCheck) UUIDS() []string {
	return i.uuids
}

func (i *IDCheck) Name() string {
	return i.name
}

func searchIDCheck[T any](v T) map[string]IDCheckMap {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		if !rv.IsValid() {
			return nil
		}
		if rv.IsNil() {
			return nil
		}
	}

	out := map[string]IDCheckMap{}
	recursion("", "", rv, out)
	return out
}

func recursion(label string, jsonLabel string, rv reflect.Value, out map[string]IDCheckMap) {
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		if !rv.IsValid() {
			return
		}
		if rv.IsNil() {
			return
		}
		recursion(label, jsonLabel, rv.Elem(), out)
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
				recursion(label, jsonLabel, fieldValue, out)
				continue
			}
			tag, has := fieldType.Tag.Lookup(TagCheck)
			if has {
				tag = readLabelName(tag)
			}
			jsonTag, has := fieldType.Tag.Lookup("json")
			if has {
				jsonTag = readLabelName(jsonTag)
			}

			recursion(tag, jsonTag, fieldValue, out)
		}
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			recursion(label, jsonLabel, rv.Index(i), out)
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			recursion(label, jsonLabel, rv.MapIndex(key), out)
		}
	case reflect.String:
		if label != "" {
			if _, ok := out[label]; !ok {
				out[label] = make(IDCheckMap)
			}
			if _, ok := out[label][jsonLabel]; !ok {
				out[label][jsonLabel] = &IDCheck{
					uuids: []string{rv.String()},
					name:  jsonLabel,
				}
			} else {
				out[label][jsonLabel].uuids = append(out[label][jsonLabel].uuids, rv.String())
			}
		}
	default:
		return
	}
}
