package function

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type ParamType int32

const (
	NULL ParamType = -1
	JSON ParamType = 0
	PATH ParamType = 1
)

type Parameter struct {
	paramType   ParamType
	path        common.Path
	lateBinding ILateBindingValue
	evaluated   bool
	json        string
}

func (p *Parameter) GetValue() (interface{}, error) {
	return p.lateBinding.Get()
}

func (p *Parameter) SetLateBinding(lateBinding ILateBindingValue) {
	p.lateBinding = lateBinding
}

func (p *Parameter) GetPath() common.Path {
	return p.path
}

func (p *Parameter) SetEvaluated(evaluated bool) {
	p.evaluated = evaluated
}

func (p *Parameter) HasEvaluated() bool {
	return p.evaluated
}

func (p *Parameter) GetType() ParamType {
	return p.paramType
}

func (p *Parameter) SetType(paramType ParamType) {
	p.paramType = paramType
}

func (p *Parameter) SetPath(path common.Path) {
	p.path = path
}

func (p *Parameter) GetJson() string {
	return p.json
}

func (p *Parameter) GetILateBindingValue() ILateBindingValue {
	return p.lateBinding
}

func CreateJsonParameter(json string) *Parameter {
	return &Parameter{json: json, paramType: JSON}
}

func CreatePathParameter(p common.Path) *Parameter {
	return &Parameter{path: p, paramType: PATH}
}

func ParametersToList(typeName common.Type, ctx common.EvaluationContext, parameters []*Parameter) ([]interface{}, error) {
	var values []interface{}
	for _, param := range parameters {
		value, err := param.GetValue()
		if err != nil {
			return nil, err
		}
		values, err = parameterConsume(typeName, ctx, values, value)
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

func parameterConsume(expectedType common.Type, ctx common.EvaluationContext, collection []interface{}, value interface{}) ([]interface{}, error) {
	if ctx.Configuration().JsonProvider().IsArray(value) {
		array, err := ctx.Configuration().JsonProvider().ToArray(value)
		if err != nil {
			return nil, err
		}
		for _, o := range array {
			if expectedType == common.TYPE_NUMBER {
				collection = append(collection, o)
			} else {
				collection = append(collection, common.UtilsToString(o))
			}
		}
	} else {
		if expectedType == common.TYPE_NUMBER {
			collection = append(collection, value)
		} else {
			collection = append(collection, common.UtilsToString(value))
		}
	}
	return collection, nil
}
