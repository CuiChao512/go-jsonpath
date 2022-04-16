package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/configuration"
	"cuichao.com/go-jsonpath/jsonpath/predicate"
)

type EvaluationListener interface {
}

type ReadContext interface {
	Configuration() *configuration.Configuration
	Json() interface{}
	JsonString() string
	ReadWithFilters(path string, filters ...*predicate.Predicate)
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
