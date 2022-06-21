package path

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/function"
)

type KeySetFunction struct {
}

func (*KeySetFunction) Next(value interface{}) {}
func (*KeySetFunction) GetValue() interface{}  { return nil }
func (*KeySetFunction) Invoke(nextAndGet PathFunctionNextAndGet, currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters []*function.Parameter) (interface{}, error) {
	if ctx.Configuration().JsonProvider().IsMap(model) {
		return ctx.Configuration().JsonProvider().GetPropertyKeys(model)
	}
	return nil, nil
}

type Append struct {
}

func (*Append) Next(value interface{}) {}
func (*Append) GetValue() interface{}  { return nil }

func (*Append) Invoke(nextAndGet PathFunctionNextAndGet, currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters []*function.Parameter) (interface{}, error) {
	jsonProvider := ctx.Configuration().JsonProvider()
	var array []interface{}
	if jsonProvider.IsArray(model) {
		ok := false
		array, ok = model.([]interface{})
		if ok && parameters != nil && len(parameters) > 0 {
			for _, param := range parameters {
				l := len(array)
				val, err := param.GetValue()
				if err != nil {
					return nil, err
				}
				err = jsonProvider.SetArrayIndex(&array, l, val)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return array, nil
}
