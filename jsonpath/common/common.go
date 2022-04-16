package common

type MapFunction interface {
	Map(currentValue interface{}, configuration *Configuration) interface{}
}
