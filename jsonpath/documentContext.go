package jsonpath

type EvaluationListener interface {
}

type ReadContext interface {
	Configuration() *Configuration
	Json() interface{}
	JsonString() string
	ReadWithFilters(path string, filters ...*Predicate)
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
