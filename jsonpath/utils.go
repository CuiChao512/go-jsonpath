package jsonpath

import (
	"fmt"
	"reflect"
)

func UtilsJoin(delimiter string, warp string, values interface{}) string {
	r := ""
	s := reflect.ValueOf(values)
	if s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			r += delimiter + warp + UtilsToString(s.Index(i).Interface())
		}
	} else {
		return delimiter + warp + UtilsToString(values)
	}
	return r
}

func UtilsSliceContains(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			if item == s.Index(i).Interface() {
				return true
			}
		}
	}
	return false
}

func UtilsIsSlice(slice interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() == reflect.Slice {
		return true
	}
	return false
}

func UtilsIsMap(mapObj interface{}) bool {
	s := reflect.ValueOf(mapObj)
	if s.Kind() == reflect.Map {
		return true
	}
	return false
}

func UtilsConcat(string ...string) string {
	//TODO: Utils.concat
	return ""
}

func UtilsToString(obj ...interface{}) string {
	return fmt.Sprint(obj)
}

func UtilsCharIsDigit(char rune) bool {
	return char == '0' || char == '1' || char == '2' || char == '3' || char == '4' ||
		char == '5' || char == '6' || char == '7' || char == '8' || char == '9'
}
