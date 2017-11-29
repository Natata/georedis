package geochat

import (
	"fmt"
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

func rawToNeighbors(r interface{}, options ...Option) ([]*NeighborData, error) {
	v := reflect.ValueOf(r)

	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("wrong type: %v", v.Kind())
	}

	results := make([]*NeighborData, v.Len())
	var err error
	for i := 0; i < v.Len(); i++ {
		results[i], err = NewNeighborData(unpackValue(v.Index(i)), options...)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
