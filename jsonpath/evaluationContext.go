package jsonpath

import "cuichao.com/go-jsonpath/jsonpath/path"

type EvaluationContext interface {
	Configuration() *Configuration
	RootDocument() interface{}
	GetValue() interface{}
	GetValueUnwrap(unwrap bool) interface{}
}

type EvaluationContextImpl struct {
}

func (e *EvaluationContextImpl) Configuration() *Configuration {
	return nil
}
func (e *EvaluationContextImpl) RootDocument() interface{} {
	return nil
}

func (e *EvaluationContextImpl) GetValue() interface{} {
	return nil
}

func (e *EvaluationContextImpl) GetValueUnwrap(unwrap bool) interface{} {
	return nil
}

func (e *EvaluationContextImpl) ForUpdate() bool {
	return false
}

func (e *EvaluationContextImpl) AddResult(pathString string, operation path.Ref, model interface{}) {

}
