package function

import "cuichao.com/go-jsonpath/jsonpath"

type PathFunction interface {
	Invoke(currentPath string, parent PathRef, model interface{}, ctx jsonpath.EvaluationContext, parameters []*Parameter)
}
