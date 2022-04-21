package path

import (
	"errors"
	"fmt"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

var documentEvalCache = map[common.Path]interface{}{}

type EvaluationContextImpl struct {
	configuration     *common.Configuration
	forUpdate         bool
	path              common.Path
	rootDocument      interface{}
	updateOperations  []common.PathRef
	valueResult       []interface{}
	pathResult        []interface{}
	suppressException bool
	resultIndex       int
}

func (*EvaluationContextImpl) DocumentEvalCache() map[common.Path]interface{} {
	return documentEvalCache
}

func (e *EvaluationContextImpl) GetRoot() (*RootPathToken, error) {
	compiledPath, ok := e.path.(*CompiledPath)
	if !ok {
		return nil, errors.New("path can not cast to *CompiledPath")
	}
	return compiledPath.root, nil
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

func (e *EvaluationContextImpl) GetValue() (interface{}, error) {
	return e.GetValueUnwrap(true)
}

func (e *EvaluationContextImpl) GetValueUnwrap(unwrap bool) (interface{}, error) {
	if e.path.IsDefinite() {
		if e.resultIndex == 0 {
			if e.suppressException {
				return nil, nil
			}
			return nil, fmt.Errorf("no result:%s", &common.PathNotFoundError{Message: "No results for path: " + e.path.String()})
		}
		if length, err := e.JsonProvider().Length(e.valueResult); err != nil {
			return nil, err
		} else {
			var value interface{}
			if length > 0 {
				value = e.JsonProvider().GetArrayIndex(e.valueResult, length-1)
			}

			if value != nil && unwrap {
				value = e.JsonProvider().Unwrap(value)
			}
			return value, nil
		}
	}
	return e.valueResult, nil
}

func (e *EvaluationContextImpl) ForUpdate() bool {
	return e.forUpdate
}

func (e *EvaluationContextImpl) AddResult(pathString string, operation common.PathRef, model interface{}) error {
	if e.forUpdate {
		e.updateOperations = append(e.updateOperations, operation)
	}

	if err := e.configuration.JsonProvider().SetArrayIndex(&e.valueResult, e.resultIndex, model); err != nil {
		return err
	}
	if err := e.configuration.JsonProvider().SetArrayIndex(&e.pathResult, e.resultIndex, pathString); err != nil {
		return err
	}

	e.resultIndex++

	if len(e.configuration.GetEvaluationListeners()) == 0 {
		idx := e.resultIndex - 1
		for _, listener := range e.configuration.GetEvaluationListeners() {
			continuation := listener.ResultFound(createFoundResultImpl(idx, pathString, model))
			if continuation == common.ABORT {
				return &common.EvaluationAbortError{}
			}
		}
	}
	return nil
}

func CreateEvaluationContextImpl(path common.Path, rootDocument interface{}, configuration *common.Configuration, forUpdate bool) *EvaluationContextImpl {
	e := &EvaluationContextImpl{}
	e.forUpdate = forUpdate
	e.path = path
	e.rootDocument = rootDocument
	e.configuration = configuration
	e.valueResult = configuration.JsonProvider().CreateArray()
	e.pathResult = configuration.JsonProvider().CreateArray()
	e.updateOperations = []common.PathRef{}
	e.suppressException = common.UtilsSliceContains(configuration.Options(), common.OPTION_SUPPRESS_EXCEPTIONS)
	return e
}

type FoundResultImpl struct {
	index  int
	path   string
	result interface{}
}

func (f *FoundResultImpl) Index() int {
	return f.index
}

func (f *FoundResultImpl) Path() string {
	return f.path
}

func (f *FoundResultImpl) Result() interface{} {
	return f.result
}

func createFoundResultImpl(idx int, path string, model interface{}) common.FoundResult {
	return &FoundResultImpl{index: idx, path: path, result: model}
}
