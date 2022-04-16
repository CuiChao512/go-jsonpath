package common

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func UtilsStringSliceContainsAll(strings1 []string, strings2 []string) bool {
	if strings1 == nil || strings2 == nil {
		return false
	}
	if len(strings1) != len(strings2) {
		return false
	}
	for _, str1 := range strings1 {
		if !UtilsSliceContains(strings2, str1) {
			return false
		}
	}
	return true
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

func UtilsConcat(s ...string) string {
	result := ""
	for _, str := range s {
		result += str
	}
	return result
}

func UtilsToString(obj ...interface{}) string {
	return fmt.Sprint(obj)
}

func UtilsCharIsDigit(char rune) bool {
	return char == '0' || char == '1' || char == '2' || char == '3' || char == '4' ||
		char == '5' || char == '6' || char == '7' || char == '8' || char == '9'
}

func UtilsGetPtrElem(ptr interface{}) interface{} {
	val := reflect.ValueOf(ptr)
	if val.Kind() == reflect.Ptr {
		return UtilsGetPtrElem(val.Elem().Interface())
	} else {
		return ptr
	}
}

func UtilsMaxInt(int1 int, int2 int) int {
	if int1 >= int2 {
		return int1
	}
	return int2
}

func UtilsMinInt(int1 int, int2 int) int {
	if int1 <= int2 {
		return int1
	}
	return int2
}

func UtilsGetTypeName(i interface{}) string {
	return reflect.TypeOf(UtilsGetPtrElem(i)).Name()
}

func UtilsIsFloat(v interface{}) bool {
	kind := reflect.ValueOf(v).Kind()

	return kind == reflect.Float32 || kind == reflect.Float64
}

func UtilsIsInt(v interface{}) bool {
	kind := reflect.ValueOf(v).Kind()

	return kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64

}

func UtilsIsNumber(v interface{}) bool {
	return UtilsIsFloat(v) || UtilsIsInt(v)
}

func UtilsNumberToFloat64(v interface{}) (float64, error) {
	if UtilsIsNumber(v) {
		switch v.(type) {
		case int:
			vv, _ := v.(int)
			return float64(vv), nil
		case int8:
			vv, _ := v.(int8)
			return float64(vv), nil
		case int16:
			vv, _ := v.(int16)
			return float64(vv), nil
		case int32:
			vv, _ := v.(int32)
			return float64(vv), nil
		case int64:
			vv, _ := v.(int64)
			return float64(vv), nil
		case float32:
			vv, _ := v.(float32)
			return float64(vv), nil
		case float64:
			vv, _ := v.(float64)
			return vv, nil
		}
	}
	return 0, errors.New("not a number")
}

func UtilsNumberToFloat64Force(v interface{}) float64 {
	f, _ := UtilsNumberToFloat64(v)
	return f
}

func UtilsStringUnescape(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	writer := new(strings.Builder)
	unicode := new(strings.Builder)
	hadSlash := false
	inUnicode := false
	runes := []rune(str)
	for i := 0; i < len(str); i++ {
		ch := runes[i]
		if inUnicode {
			unicode.WriteRune(ch)
			if unicode.Len() == 4 {
				value, err := strconv.ParseInt(unicode.String(), 16, 0)
				if err != nil {
					return "", &JsonPathError{Message: "Unable to parse unicode value: " + unicode.String()}
				}
				writer.WriteRune(rune(value))
				unicode.Reset()
				inUnicode = false
				hadSlash = false
			}
			continue
		}
		if hadSlash {
			hadSlash = false
			switch ch {
			case '\\':
				writer.WriteRune('\\')
			case '\'':
				writer.WriteRune('\'')
			//case '\"':
			case '"':
				writer.WriteRune('"')
			case 'r':
				writer.WriteRune('\r')
			case 'f':
				writer.WriteRune('\f')
			case 't':
				writer.WriteRune('\t')
			case 'n':
				writer.WriteRune('\n')
			case 'b':
				writer.WriteRune('\b')
			case 'u':
				{
					inUnicode = true
				}
			default:
				writer.WriteRune(ch)
			}
			continue
		} else if ch == '\\' {
			hadSlash = true
			continue
		}
		writer.WriteRune(ch)
	}
	if hadSlash {
		writer.WriteRune('\\')
	}
	return writer.String(), nil
}
