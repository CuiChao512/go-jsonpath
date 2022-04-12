package jsonpath

type Option int

const (
	OPTION_DEFAULT_PATH_LEAF_TO_NULL Option = 0
	OPTION_ALWAYS_RETURN_LIST        Option = 1
	OPTION_AS_PATH_LIST              Option = 2
	OPTION_SUPPRESS_EXCEPTIONS       Option = 3
	OPTION_REQUIRE_PROPERTIES        Option = 4
)

type Configuration struct {
	jsonProvider JsonProvider
	options      []Option
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
