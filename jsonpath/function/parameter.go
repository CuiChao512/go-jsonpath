package function

import (
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
