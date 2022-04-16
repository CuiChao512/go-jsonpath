package path

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/evaluationContext"
)

type CompiledPath struct {
	root       *RootPathToken
	isRootPath bool
}

func (cp *CompiledPath) Evaluate(document interface{}, rootDocument interface{}, configuration *common.Configuration) (evaluationContext.EvaluationContext, error) {
	return nil, nil
}

func (cp *CompiledPath) EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *common.Configuration, forUpdate bool) evaluationContext.EvaluationContext {
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
