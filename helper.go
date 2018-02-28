package georedis

import (
	"reflect"
)

func unpackValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		if !v.IsNil() {
			v = v.Elem()
		}
	}
	return v
}
