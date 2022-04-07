package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
)

type Path interface {
	Evaluate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration) (jsonpath.EvaluationContext, error)
	EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration, forUpdate bool) jsonpath.EvaluationContext
	String() string
	IsDefinite() bool
	IsFunctionPath() bool
	IsRootPath() bool
}
