package function

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type ParamType int32

const (
	JSON ParamType = 0
	PATH ParamType = 1
)

type Parameter struct {
	paramType   ParamType
	path        path.Path
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

func (p *Parameter) GetPath() path.Path {
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

func (p *Parameter) SetPath(path path.Path) {
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

func CreatePathParameter(p path.Path) *Parameter {
	return &Parameter{path: p, paramType: PATH}
}

func ParametersToList(typeName common.Type, ctx path.EvaluationContext, parameters []*Parameter) ([]interface{}, error) {
	var values *[]interface{}
	for _, param := range parameters {
		value, err := param.GetValue()
		if err != nil {
			return nil, err
		}
		parameterConsume(typeName, ctx, values, value)
	}
	return *values, nil
}

func parameterConsume(expectedType common.Type, ctx path.EvaluationContext, collection *[]interface{}, value interface{}) {
	if ctx.Configuration().JsonProvider().IsArray(value) {
		for _, o := range ctx.Configuration().JsonProvider().ToIterable(value) {
			if expectedType == common.TYPE_NUMBER {
				*collection = append(*collection, o)
			} else {
				*collection = append(*collection, common.UtilsToString(o))
			}
		}
	} else {
		if expectedType == common.TYPE_NUMBER {
			*collection = append(*collection, value)
		} else {
			*collection = append(*collection, common.UtilsToString(value))
		}
	}
}
