package function

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
)

type ILateBindingValue interface {
	Get() (interface{}, error)
}

type LateBindingValue struct {
	path          common.Path
	rootDocument  string
	configuration *common.Configuration
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

	if common.UtilsGetPtrElem(l) != common.UtilsGetPtrElem(o) {
		return false
	}

	that, _ := o.(*LateBindingValue)

	return l.path == that.path && l.rootDocument == that.rootDocument && l.configuration == that.configuration
}

func CreateLateBindingValue(path common.Path, rootDocument interface{}, configuration *common.Configuration) (*LateBindingValue, error) {
	l := &LateBindingValue{}
	l.path = path
	l.rootDocument = common.UtilsToString(rootDocument)
	l.configuration = configuration
	e, err := path.Evaluate(rootDocument, rootDocument, configuration)
	if err != nil {
		return nil, err
	}
	l.result, err = e.GetValue()
	if err != nil {
		return nil, err
	}

	return l, nil
}

type JsonLateBindingValue struct {
	jsonProvider  common.JsonProvider
	jsonParameter *Parameter
}

func (j *JsonLateBindingValue) Get() (interface{}, error) {
	return j.jsonProvider.Parse(j.jsonParameter.GetJson())
}

func CreateJsonLateBindingValue(jsonProvider common.JsonProvider, jsonParameter *Parameter) *JsonLateBindingValue {
	return &JsonLateBindingValue{
		jsonParameter: jsonParameter,
		jsonProvider:  jsonProvider,
	}
}
