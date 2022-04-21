package jsonpath

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type EvaluationListener interface {
}

type ReadContext interface {
	Configuration() *common.Configuration
	Json() interface{}
	JsonString() string
	ReadWithFilters(path string, filters ...common.Predicate)
	Read(path string)
	Limit(maxResults int64) *ReadContext
	WithListeners(listeners ...*EvaluationListener)
}

type WriteContext interface {
}

type DocumentContext interface {
	ReadContext
	WriteContext
}

type JsonContext struct {
}
