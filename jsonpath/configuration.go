package jsonpath

type Configuration struct {
	jsonProvider JsonProvider
}

func (c *Configuration) JsonProvider() JsonProvider {
	return c.jsonProvider
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
	GetPropertyKeys(obj interface{}) []string
	GetArrayIndex(obj interface{}, idx int) interface{}
	SetArrayIndex(array interface{}, idx int, newValue interface{})
	GetMapValue(obj interface{}, key string) interface{}
	SetProperty(obj interface{}, key interface{}, value interface{})
	RemoveProperty(obj interface{}, key interface{})
	IsMap(obj interface{}) bool
	Unwrap(obj interface{}) interface{}
}
