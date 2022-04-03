package path

import "cuichao.com/go-jsonpath/jsonpath"

type Compiler struct {
}

type CompiledPath struct {
}

func (cp *CompiledPath) Evaluate(document *interface{}, rootDocument *interface{}, configuration *jsonpath.Configuration) *jsonpath.EvaluationContext {
	return nil
}

func (cp *CompiledPath) EvaluateForUpdate(document *interface{}, rootDocument *interface{}, configuration *jsonpath.Configuration, forUpdate bool) *jsonpath.EvaluationContext {
	return nil
}

func (cp *CompiledPath) ToString() string {
	return ""
}

func Compile(pathString string) Path {
	return &CompiledPath{}
}
