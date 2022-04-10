package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
)

type Path interface {
	Evaluate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration) (jsonpath.EvaluationContext, error)
	EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration, forUpdate bool) jsonpath.EvaluationContext
	String() string
	IsDefinite() bool
	IsFunctionPath() bool
	IsRootPath() bool
}

type Ref interface {
	GetAccessor() interface{}
	Set(newVal interface{}, configuration *jsonpath.Configuration)
	Convert(mapFunction jsonpath.MapFunction, configuration *jsonpath.Configuration)
	Delete(configuration *jsonpath.Configuration)
	Add(newVal interface{}, configuration *jsonpath.Configuration)
	Put(key string, newVal interface{}, configuration *jsonpath.Configuration)
	RenameKey(oldKeyName string, newKeyName string, configuration *jsonpath.Configuration)
}
