package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
)

type Path interface {
	Evaluate(document *interface{}, rootDocument *interface{}, configuration *jsonpath.Configuration) *jsonpath.EvaluationContext
	EvaluateForUpdate(document *interface{}, rootDocument *interface{}, configuration *jsonpath.Configuration, forUpdate bool) *jsonpath.EvaluationContext
}
