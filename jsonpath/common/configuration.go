package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Option int

const (
	OPTION_DEFAULT_PATH_LEAF_TO_NULL Option = 0
	OPTION_ALWAYS_RETURN_LIST        Option = 1
	OPTION_AS_PATH_LIST              Option = 2
	OPTION_SUPPRESS_EXCEPTIONS       Option = 3
	OPTION_REQUIRE_PROPERTIES        Option = 4
)

type Configuration struct {
	jsonProvider        JsonProvider
	options             []Option
	mappingProvider     MappingProvider
	evaluationListeners []EvaluationListener
}

func (c *Configuration) JsonProvider() JsonProvider {
	return c.jsonProvider
}

func (c *Configuration) Options() []Option {
	return c.options
}

func (c *Configuration) MappingProvider() MappingProvider {
	return c.mappingProvider
}

func (c *Configuration) GetEvaluationListeners() []EvaluationListener {
	return c.evaluationListeners
}

var JsonProviderUndefined interface{}

type JsonProvider interface {
	IsArray(obj interface{}) bool
	IsMap(obj interface{}) bool
	GetArrayIndex(obj interface{}, idx int) interface{}
	GetMapValue(obj interface{}, key string) interface{}
	SetArrayIndex(array interface{}, idx int, newValue interface{}) error
	SetProperty(obj interface{}, key interface{}, value interface{}) error
	Parse(json string) (interface{}, error)
	ToJson(obj interface{}) (string, error)
	CreateArray() []interface{}
	CreateMap() map[string]interface{}
	Length(obj interface{}) (int, error)
	ToArray(obj interface{}) ([]interface{}, error)
	GetPropertyKeys(obj interface{}) ([]string, error)
	RemoveProperty(obj interface{}, key interface{}) error
	Unwrap(obj interface{}) interface{}
}

type MappingProvider interface {
	MapSlice(data interface{}, configuration *Configuration) interface{}
	MapMap(data interface{}, configuration *Configuration) interface{}
}

// defaultJsonProvider -----

type NativeJsonProvider struct {
}

func (*NativeJsonProvider) IsArray(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Slice
}

func (*NativeJsonProvider) IsMap(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Map
}

func (*NativeJsonProvider) GetArrayIndex(obj interface{}, idx int) interface{} {
	l, ok := obj.([]interface{})
	if !ok {
		return nil
	}
	return l[idx]
}

func (d *NativeJsonProvider) GetArrayIndexByUnwrap(obj interface{}, idx int, unwrap bool) interface{} {
	return d.GetArrayIndex(obj, idx)
}

func (d *NativeJsonProvider) SetArrayIndex(array interface{}, index int, newValue interface{}) error {
	if !d.IsArray(array) {
		return errors.New("unsupported operation")
	} else {
		l, _ := array.([]interface{})
		if index == len(l) {
			l = append(l, newValue)
		} else {
			l[index] = newValue
		}
		return nil
	}
}

func (d *NativeJsonProvider) GetMapValue(obj interface{}, key string) interface{} {
	m, ok := obj.(map[string]interface{})
	if !ok {
		return JsonProviderUndefined
	}
	value := m[key]
	if value == nil {
		return JsonProviderUndefined
	} else {
		return value
	}
}

func (d *NativeJsonProvider) GetPropertyKeys(obj interface{}) ([]string, error) {
	if d.IsArray(obj) {
		return nil, errors.New("slice dose not support getPropertyKeys operation")
	} else {
		m, ok := obj.(map[string]interface{})
		if !ok {
			return nil, errors.New("slice dose not support getPropertyKeys operation")
		}
		keys := make([]string, 0, len(m))
		for k, _ := range m {
			keys = append(keys, k)
		}
		return keys, nil
	}
}

func (d *NativeJsonProvider) SetProperty(obj interface{}, key interface{}, value interface{}) error {
	if d.IsMap(obj) {
		m, _ := obj.(map[string]interface{})
		m[UtilsToString(key)] = value
		return nil
	} else {
		return &JsonPathError{Message: "setProperty operation cannot be used with " + getTypeString(obj)}
	}
}

func (d *NativeJsonProvider) RemoveProperty(obj interface{}, key interface{}) error {
	if d.IsMap(obj) {
		m, _ := obj.(map[string]interface{})
		delete(m, UtilsToString(key))
	} else {
		s, _ := obj.([]interface{})

		var index int

		switch key.(type) {
		case int:
			k, _ := key.(int)
			index = k
		case string:
			keyString, _ := key.(string)
			var err error
			index, err = strconv.Atoi(keyString)
			if err != nil {
				return errors.New("%s can not as an index")
			}
		}
		s = append(s[:index], s[index+1:])
	}
	return nil
}

func (d *NativeJsonProvider) ToArray(obj interface{}) ([]interface{}, error) {
	if d.IsArray(obj) {
		s, _ := obj.([]interface{})
		return s, nil
	} else {
		return nil, &JsonPathError{Message: fmt.Sprintf("%s is not a slice", getTypeString(obj))}
	}
}

func (d *NativeJsonProvider) Length(obj interface{}) (int, error) {
	if d.IsArray(obj) || d.IsMap(obj) {
		return reflect.ValueOf(obj).Len(), nil
	} else {
		if reflect.TypeOf(obj).Kind() == reflect.String {
			return reflect.ValueOf(obj).Len(), nil
		}
	}

	return -1, &JsonPathError{Message: "length operation cannot be applied to " + getTypeString(obj)}
}

func getTypeString(obj interface{}) string {
	if obj == nil {
		return "null"
	} else {
		return reflect.TypeOf(obj).Kind().String()
	}
}

func (*NativeJsonProvider) Unwrap(obj interface{}) interface{} {
	return obj
}

func (*NativeJsonProvider) isPrimitiveNumber(n interface{}) bool {
	switch n.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}

func (d *NativeJsonProvider) unwrapNumber(number interface{}) interface{} {
	var unwrapNumber interface{}
	if d.isPrimitiveNumber(d) {
		unwrapNumber = number
	}
	return unwrapNumber
}

func (*NativeJsonProvider) Parse(jsonString string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*NativeJsonProvider) ToJson(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (*NativeJsonProvider) CreateArray() []interface{} {
	return []interface{}{}
}

func (*NativeJsonProvider) CreateMap() map[string]interface{} {
	return map[string]interface{}{}
}

type NativeMappingProvider struct{}

func (n *NativeMappingProvider) MapSlice(data interface{}, configuration *Configuration) interface{} {
	return data
}

func (*NativeMappingProvider) MapMap(data interface{}, configuration *Configuration) interface{} {
	return data
}

func DefaultConfiguration() *Configuration {
	return &Configuration{jsonProvider: &NativeJsonProvider{}, mappingProvider: &NativeMappingProvider{}}
}
