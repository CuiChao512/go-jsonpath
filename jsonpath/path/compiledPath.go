package path

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/function"
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

func invertScannerFunctionRelationship(path *RootPathToken) (*RootPathToken, error) {
	if path.IsFunctionPath() {
		next, err := path.nextToken()
		if err != nil {
			return nil, err
		}
		switch next.(type) {
		case *ScanPathToken:
			var token Token = path
			var prior Token = nil

			for true {
				token, err = token.nextToken()
				if token != nil {
					switch token.(type) {
					case *FunctionPathToken:
						continue
					}
				}
				break
			}
			// Invert the relationship $..path.function() to $.function($..path)
			switch token.(type) {
			case *FunctionPathToken:
				prior.SetNext(nil)
				path.SetTail(prior)

				// Now generate a new parameter from our path
				parameter := &function.Parameter{}
				compiledPath, err := CreateCompiledPath(path, true)
				if err != nil {
					return nil, err
				}
				parameter.SetPath(compiledPath)
				parameter.SetType(function.PATH)
				functionToken, _ := token.(*FunctionPathToken)
				functionToken.SetParameters([]*function.Parameter{parameter})
				functionRoot := CreateRootPathToken('$')
				functionRoot.SetTail(functionToken)
				functionRoot.SetNext(functionToken)

				// Define the function as the root
				return functionRoot, nil

			}
		}
	}
	return path, nil
}

func CreateCompiledPath(rootPathToken *RootPathToken, isRootPath bool) (*CompiledPath, error) {
	newRoot, err := invertScannerFunctionRelationship(rootPathToken)
	if err != nil {
		return nil, err
	}
	return &CompiledPath{root: newRoot, isRootPath: isRootPath}, nil
}
