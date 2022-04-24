package jsonpath

import (
	"errors"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type ReadContext interface {
	Configuration() *common.Configuration
	Json() interface{}
	JsonString() (string, error)
	ReadWithFilters(path string, filters ...common.Predicate) (interface{}, error)
	Read(path string) (interface{}, error)
	ReadJsonpath(path *Jsonpath) (interface{}, error)
	Limit(maxResults int) (ReadContext, error)
	WithListeners(listeners ...common.EvaluationListener) (ReadContext, error)
}

type WriteContext interface {
}

type DocumentContext interface {
	ReadContext
	WriteContext
}

type JsonContext struct {
	configuration *common.Configuration
	json          interface{}
}

func (jc *JsonContext) Configuration() *common.Configuration {
	return jc.configuration
}

func (jc *JsonContext) Json() interface{} {
	return jc.json
}

func (jc *JsonContext) JsonString() (string, error) {
	return jc.configuration.JsonProvider().ToJson(jc.json)
}

var jsonPathCache = make(map[string]*Jsonpath)

func (jc *JsonContext) pathFromCache(pathString string, filters []common.Predicate) (*Jsonpath, error) {
	var cacheKey string
	if filters == nil || len(filters) == 0 {
		cacheKey = pathString
	} else {
		cacheKey = common.UtilsConcat(pathString, common.UtilsToString(filters))
	}
	jp := jsonPathCache[cacheKey]
	if jp == nil {
		jsonpath, err := compileJsonpathByStringAndPredicateSlice(pathString, filters)
		if err != nil {
			return nil, err
		}
		jsonPathCache[cacheKey] = jsonpath
		return jsonpath, nil
	}
	return jp, nil
}

func (jc *JsonContext) Read(pathString string) (interface{}, error) {
	return jc.ReadWithFilters(pathString)
}

func (jc *JsonContext) ReadWithFilters(pathString string, filters ...common.Predicate) (interface{}, error) {
	if pathString == "" {
		return nil, errors.New("path can not be empty")
	}
	jp, err := jc.pathFromCache(pathString, filters)
	if err != nil {
		return nil, err
	}

	return jc.ReadJsonpath(jp)
}

func (jc *JsonContext) ReadJsonpath(path *Jsonpath) (interface{}, error) {
	if path == nil {
		return nil, errors.New("path can not be nil")
	}
	return path.readAnyByConfiguration(jc.json, jc.configuration)
}

func (jc *JsonContext) Limit(maxResults int) (ReadContext, error) {
	return jc.WithListeners(createLimitingEvaluationListener(maxResults))
}

func (jc *JsonContext) WithListeners(listeners ...common.EvaluationListener) (ReadContext, error) {
	newConfig := common.CreateConfigurationByJsonProviderOptionsMappingProviderEvaluationListeners(
		jc.configuration.JsonProvider(), jc.configuration.Options(), jc.configuration.MappingProvider(), listeners)
	return CreateJsonContextByAny(jc.json, newConfig)
}

type LimitingEvaluationListener struct {
	limit int
}

func (l *LimitingEvaluationListener) ResultFound(found common.FoundResult) common.EvaluationContinuation {
	if found.Index() == l.limit-1 {
		return common.ABORT
	} else {
		return common.CONTINUE
	}
}

func createLimitingEvaluationListener(maxResults int) *LimitingEvaluationListener {
	return &LimitingEvaluationListener{limit: maxResults}
}

func CreateJsonContextByAny(obj interface{}, configuration *common.Configuration) (*JsonContext, error) {
	if obj == nil {
		return nil, errors.New("json can not be nil")
	}
	if configuration == nil {
		return nil, errors.New("configuration can not be nil")
	}
	return &JsonContext{json: obj, configuration: configuration}, nil
}
