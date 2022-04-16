package path

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
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
	path              Path
	rootDocument      interface{}
	updateOperations  []PathRef
	documentEvalCache map[Path]interface{}
	suppressException bool
	resultIndex       int
}

func (e *EvaluationContextImpl) DocumentEvalCache() map[Path]interface{} {
	return e.documentEvalCache
}

func (e *EvaluationContextImpl) GetRoot() *RootPathToken {
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

func (e *EvaluationContextImpl) AddResult(pathString string, operation PathRef, model interface{}) {

}
