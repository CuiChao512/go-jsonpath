package function

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type ILateBindingValue interface {
	Get() (interface{}, error)
}

type LateBindingValue struct {
	path          path.Path
	rootDocument  string
	configuration *jsonpath.Configuration
	result        interface{}
}

func (l *LateBindingValue) Get() (interface{}, error) {
	return l.result, nil
}

func (l *LateBindingValue) Equals(o interface{}) bool {
	if l == o {
		return true
	}

	if o == nil {
		return false
	}

	if jsonpath.UtilsGetPtrElem(l) != jsonpath.UtilsGetPtrElem(o) {
		return false
	}

	that, _ := o.(*LateBindingValue)

	return l.path == that.path && l.rootDocument == that.rootDocument && l.configuration == that.configuration
}

func CreateLateBindingValue(path path.Path, rootDocument interface{}, configuration *jsonpath.Configuration) (*LateBindingValue, error) {
	l := &LateBindingValue{}
	l.path = path
	l.rootDocument = jsonpath.UtilsToString(rootDocument)
	l.configuration = configuration
	e, err := path.Evaluate(rootDocument, rootDocument, configuration)
	if err != nil {
		return nil, err
	}
	l.result = e.GetValue()
	return l, nil
}

type JsonLateBindingValue struct {
	jsonProvider  jsonpath.JsonProvider
	jsonParameter *Parameter
}

func (j *JsonLateBindingValue) Get() (interface{}, error) {
	return j.jsonProvider.Parse(j.jsonParameter.GetJson())
}

func CreateJsonLateBindingValue(jsonProvider jsonpath.JsonProvider, jsonParameter *Parameter) *JsonLateBindingValue {
	return &JsonLateBindingValue{
		jsonParameter: jsonParameter,
		jsonProvider:  jsonProvider,
	}
}
