package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type EvaluationListener interface {
}

type ReadContext interface {
	Configuration() *common.Configuration
	Json() interface{}
	JsonString() string
	ReadWithFilters(path string, filters ...*path.Predicate)
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
