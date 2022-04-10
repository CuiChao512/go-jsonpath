package jsonpath

type Jsonpath struct {
}

func (j *Jsonpath) find(path string, data interface{}, configs *Configuration) (interface{}, error) {
	return nil, nil
}

type MapFunction interface {
	Map(currentValue interface{}, configuration *Configuration) interface{}
}
