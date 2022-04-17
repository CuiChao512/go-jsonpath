package common

import (
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
	jsonProvider    JsonProvider
	options         []Option
	mappingProvider MappingProvider
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

var JsonProviderUndefined interface{}

type JsonProvider interface {
	Parse(json string) (interface{}, error)
	ToJson(obj interface{}) (string, error)
	CreateArray() interface{}
	CreateMap() interface{}
	IsArray(obj interface{}) bool
	Length(obj interface{}) int
	ToArray(obj interface{}) []interface{}
	GetPropertyKeys(obj interface{}) ([]string, error)
	GetArrayIndex(obj interface{}, idx int) interface{}
	SetArrayIndex(array interface{}, idx int, newValue interface{})
	GetMapValue(obj interface{}, key string) interface{}
	SetProperty(obj interface{}, key interface{}, value interface{})
	RemoveProperty(obj interface{}, key interface{})
	IsMap(obj interface{}) bool
	Unwrap(obj interface{}) interface{}
	ToIterable(obj interface{}) []interface{}
}

type MappingProvider interface {
	MapSlice(data interface{}, configuration *Configuration) interface{}
	MapMap(data interface{}, configuration *Configuration) interface{}
}

// defaultJsonProvider -----

type defaultJsonProvider struct {
}

func (*defaultJsonProvider) IsArray(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Slice
}

func (*defaultJsonProvider) IsMap(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Map
}

func (d *defaultJsonProvider) GetPropertyKeys(obj interface{}) ([]string, error) {
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

func (d *defaultJsonProvider) SetProperty(obj interface{}, key interface{}, value interface{}) error {
	if d.IsMap(obj) {
		m, _ := obj.(map[string]interface{})
		m[UtilsToString(key)] = value
		return nil
	} else {
		return &JsonPathError{Message: "setProperty operation cannot be used with " + getTypeString(obj)}
	}
}

func (d *defaultJsonProvider) RemoveProperty(obj interface{}, key interface{}) error {
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

func (d *defaultJsonProvider) ToIterable(obj interface{}) (interface{}, error) {
	if d.IsArray(obj) {
		return obj, nil
	} else {
		return nil, &JsonPathError{Message: fmt.Sprintf("%s is not a slice", getTypeString(obj))}
	}
}

func (d *defaultJsonProvider) Length(obj interface{}) (int, error) {
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

func (*defaultJsonProvider) Unwrap(obj interface{}) interface{} {
	return obj
}

type NativeJsonProvider struct {
	*defaultJsonProvider
}
