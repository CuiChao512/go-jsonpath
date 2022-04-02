package filter

import "cuichao.com/go-jsonpath/jsonpath"

type Evaluator interface {
	Evaluate(left ValueNode, right ValueNode, ctx jsonpath.PredicateContext) bool
}
