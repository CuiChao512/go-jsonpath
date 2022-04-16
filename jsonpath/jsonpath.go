package jsonpath

import "cuichao.com/go-jsonpath/jsonpath/configuration"

type Jsonpath struct {
}

func (j *Jsonpath) find(path string, data interface{}, configs *configuration.Configuration) (interface{}, error) {
	return nil, nil
}

type MapFunction interface {
	Map(currentValue interface{}, configuration *configuration.Configuration) interface{}
}
