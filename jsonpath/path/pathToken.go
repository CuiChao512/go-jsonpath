package path

import "cuichao.com/go-jsonpath/jsonpath"

type Token interface {
	GetTokenCount() int
	IsPathDefinite() bool
	String() string
	Invoke(pathFunction PathFunction, currentPath string, parent PathRef, model interface{}, ctx *jsonpath.EvaluationContextImpl)
}
