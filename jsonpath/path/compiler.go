package path

import "cuichao.com/go-jsonpath/jsonpath"

type CompiledPath struct {
}

func (cp *CompiledPath) Evaluate(document *interface{}, rootDocument *interface{}, configuration *jsonpath.Configuration) jsonpath.EvaluationContext {
	return nil
}

func (cp *CompiledPath) EvaluateForUpdate(document *interface{}, rootDocument *interface{}, configuration *jsonpath.Configuration, forUpdate bool) jsonpath.EvaluationContext {
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

func Compile(pathString string) Path {
	return &CompiledPath{}
}
