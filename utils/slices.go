package utils

import (
	"fmt"
	"reflect"
	"time"
)

// Unique represents to returns unique items in slice
func Unique[S ~[]E, E any](s S) S {
	keys := make(map[any]bool)
	list := make(S, 0)
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Remove represents to remove the item from slice
func Remove[S ~[]E, E any](s S, item E) S {
	for i, v := range s {
		if reflect.DeepEqual(v, item) {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// Filter represents to filter the items from slice by field name and value
func Filter[S ~[]E, E any](s S, fieldName string, value any) S {
	var result S

	for _, entry := range s {
		val := reflect.ValueOf(entry)
		if val.Kind() == reflect.Struct {
			for i := 0; i < val.NumField(); i++ {
				if val.Type().Field(i).Name == fieldName && val.Field(i).Interface() == value {
					result = append(result, entry)
				}
			}
		}
	}

	return result
}

// Filter represents to filter the items from slice by field name and value
func GetValues[S ~[]E, E any, T any](s S, fieldName string, value T) []T {
	result := make([]T, 0)

	for _, entry := range s {
		val := reflect.ValueOf(entry)
		if val.Kind() == reflect.Struct {
			for i := 0; i < val.NumField(); i++ {
				if val.Type().Field(i).Name == fieldName {
					result = append(result, val.Field(i).Interface().(T))
				}
			}
		}
	}

	return result
}

func GenerateRandomCode() string {
	return fmt.Sprintf("%06d", time.Now().Nanosecond()%1000000)
}

func Contains[S ~[]E, E any](s S, item E) bool {
	for _, entry := range s {
		if reflect.DeepEqual(entry, item) {
			return true
		}
	}
	return false
}
