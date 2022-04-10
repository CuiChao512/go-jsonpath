package function

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type PathFunction interface {
	Invoke(currentPath string, parent path.Ref, model interface{}, ctx jsonpath.EvaluationContext, parameters []*Parameter)
}
