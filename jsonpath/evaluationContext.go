package jsonpath

type EvaluationContext interface {
	Configuration() *Configuration
	RootDocument() *interface{}
	GetValue() *interface{}
	GetValueUnwrap(unwrap bool) *interface{}
}
