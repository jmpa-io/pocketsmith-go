package pocketsmith

import (
	"reflect"
)

// A helper function to check if an interface has a value or not.
// https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1.
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array,
		reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
