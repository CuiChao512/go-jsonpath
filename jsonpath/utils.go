package jsonpath

import "reflect"

func UtilsJoin(delimiter string, warp string, stringValues []string) string {
	if len(stringValues) == 0 {
		return ""
	}
	r := ""
	for _, stringValue := range stringValues {
		r += delimiter + warp + stringValue
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
