package path

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type CompiledPath struct {
	root       *RootPathToken
	isRootPath bool
}

func (cp *CompiledPath) Evaluate(document interface{}, rootDocument interface{}, configuration *common.Configuration) (common.EvaluationContext, error) {
	return cp.EvaluateForUpdate(document, rootDocument, configuration, false)
}

func (cp *CompiledPath) EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *common.Configuration, forUpdate bool) (common.EvaluationContext, error) {
	ctx := CreateEvaluationContextImpl(cp, rootDocument, configuration, forUpdate)
	var op common.PathRef
	if ctx.ForUpdate() {
		op = CreateRootPathRef(rootDocument)
	} else {
		op = PathRefNoOp
	}
	if err := cp.root.Evaluate("", op, document, ctx); err != nil {
		return nil, err
	}
	return ctx, nil
}

func (cp *CompiledPath) String() string {
	return cp.root.String()
}

func (cp *CompiledPath) IsDefinite() bool {
	return cp.root.IsPathDefinite()
}

func (cp *CompiledPath) IsFunctionPath() bool {
	return cp.root.IsFunctionPath()
}

func (cp *CompiledPath) IsRootPath() bool {
	return cp.isRootPath
}

func (cp *CompiledPath) GetRoot() *RootPathToken {
	return cp.root
}

func CreateCompiledPath(rootPathToken *RootPathToken, isRootPath bool) *CompiledPath {
	return &CompiledPath{root: rootPathToken, isRootPath: isRootPath}
}
