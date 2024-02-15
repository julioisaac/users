package utils

import "reflect"

func IsNilOrEmpty(i interface{}) bool {
	if i == nil {
		return true
	}

	value := reflect.ValueOf(i)
	kind := value.Kind()

	switch kind {
	case reflect.Slice, reflect.Map, reflect.Chan:
		return value.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	case reflect.Int64:
		return value.IsZero()
	case reflect.String:
		return value.String() == ""
	default:
		return false
	}
}
