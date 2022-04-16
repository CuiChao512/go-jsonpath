package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/configuration"
)

type CompiledPath struct {
	root       *RootPathToken
	isRootPath bool
}

func (cp *CompiledPath) Evaluate(document interface{}, rootDocument interface{}, configuration *configuration.Configuration) (jsonpath.EvaluationContext, error) {
	return nil, nil
}

func (cp *CompiledPath) EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *configuration.Configuration, forUpdate bool) jsonpath.EvaluationContext {
	return nil
}

func (cp *CompiledPath) String() string {
	return ""
}

func (cp *CompiledPath) IsDefinite() bool {
	return false
}

func (cp *CompiledPath) IsFunctionPath() bool {
	return false
}

func (cp *CompiledPath) IsRootPath() bool {
	return false
}

func (cp *CompiledPath) GetRoot() *RootPathToken {
	return cp.root
}
