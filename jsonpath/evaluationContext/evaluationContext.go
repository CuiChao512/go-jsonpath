package evaluationContext

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type EvaluationContext interface {
	Configuration() *common.Configuration
	RootDocument() interface{}
	GetValue() interface{}
	GetValueUnwrap(unwrap bool) interface{}
}

type EvaluationContextImpl struct {
	configuration     *common.Configuration
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

func (e *EvaluationContextImpl) Configuration() *common.Configuration {
	return e.configuration
}

func (e *EvaluationContextImpl) JsonProvider() common.JsonProvider {
	return e.Configuration().JsonProvider()
}

func (e *EvaluationContextImpl) Options() []common.Option {
	return e.Configuration().Options()
}

func (e *EvaluationContextImpl) RootDocument() interface{} {
	return e.rootDocument
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
