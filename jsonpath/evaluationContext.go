package jsonpath

import "cuichao.com/go-jsonpath/jsonpath/path"

type EvaluationContext interface {
	Configuration() *Configuration
	RootDocument() interface{}
	GetValue() interface{}
	GetValueUnwrap(unwrap bool) interface{}
}

type EvaluationContextImpl struct {
	configuration     *Configuration
	forUpdate         bool
	path              path.Path
	rootDocument      interface{}
	updateOperations  []path.Ref
	documentEvalCache map[path.Path]interface{}
	suppressException bool
	resultIndex       int
}

func (e *EvaluationContextImpl) DocumentEvalCache() map[path.Path]interface{} {
	return e.documentEvalCache
}

func (e *EvaluationContextImpl) GetRoot() *path.RootPathToken {
	//TODO:
	return nil
}

func (e *EvaluationContextImpl) Configuration() *Configuration {
	return nil
}

func (e *EvaluationContextImpl) JsonProvider() JsonProvider {
	return e.Configuration().jsonProvider
}

func (e *EvaluationContextImpl) Options() []Option {
	return e.Configuration().options
}

func (e *EvaluationContextImpl) RootDocument() interface{} {
	return nil
}

func (e *EvaluationContextImpl) GetValue() interface{} {
	return nil
}

func (e *EvaluationContextImpl) ToIterable(model interface{}) []interface{} {
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
