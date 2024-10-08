package types

import "reflect"

func isArrayOrSlice(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Array || kind == reflect.Slice
}
