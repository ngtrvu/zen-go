package utils

import (
	"reflect"
)

func GetBaseStructType(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)

	// Check if i is a struct or a pointer to a struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	if t.Kind() != reflect.Struct {
		return t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// If the field is an embedded struct, return its type
		if field.Anonymous {
			return field.Type
		}
	}

	return nil
}

func CreateInstanceFromObject(obj interface{}) interface{} {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}

func CreateArrayFromObject(obj interface{}) interface{} {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return reflect.New(reflect.SliceOf(t)).Elem().Interface()
}
