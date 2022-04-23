package jsonpath

import (
	"errors"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/filter"
)

type Jsonpath struct {
	path common.Path
}

func (j *Jsonpath) GetPath() string {
	return j.path.String()
}

func (j *Jsonpath) find(path string, data interface{}, configs *common.Configuration) (interface{}, error) {
	return nil, nil
}

func (j *Jsonpath) readAnyByConfiguration(jsonObject interface{}, config *common.Configuration) (interface{}, error) {
	optAsPathList := common.UtilsSliceContains(config.Options(), common.OPTION_AS_PATH_LIST)
	optAlwaysReturnList := common.UtilsSliceContains(config.Options(), common.OPTION_ALWAYS_RETURN_LIST)
	optSuppressException := common.UtilsSliceContains(config.Options(), common.OPTION_SUPPRESS_EXCEPTIONS)

	if j.path.IsFunctionPath() {
		if optAsPathList || optAlwaysReturnList {
			if optSuppressException {
				if j.path.IsDefinite() {
					return nil, nil
				} else {
					return config.JsonProvider().CreateArray(), nil
				}
			}
		}
		evaluationContext, err := j.path.Evaluate(jsonObject, jsonObject, config)
		if err != nil {
			return nil, err
		}
		pathList, err := evaluationContext.GetPathList()
		if err != nil {
			return nil, err
		}
		if optSuppressException && len(pathList) == 0 {
			return config.JsonProvider().CreateArray(), nil
		}
		return evaluationContext.GetValueUnwrap(true)
	} else if optAsPathList {
		evaluationContext, err := j.path.Evaluate(jsonObject, jsonObject, config)
		if err != nil {
			return nil, err
		}
		pathList, err := evaluationContext.GetPathList()
		if err != nil {
			return nil, err
		}
		if optSuppressException && len(pathList) == 0 {
			return config.JsonProvider().CreateArray(), nil
		}
		return evaluationContext.GetPath()
	} else {
		evaluationContext, err := j.path.Evaluate(jsonObject, jsonObject, config)
		if err != nil {
			return nil, err
		}
		pathList, err := evaluationContext.GetPathList()
		if err != nil {
			return nil, err
		}
		if optSuppressException && len(pathList) == 0 {
			if optAlwaysReturnList {
				return config.JsonProvider().CreateArray(), nil
			} else {
				if j.path.IsDefinite() {
					return nil, nil
				} else {
					return config.JsonProvider().CreateArray(), nil
				}
			}
		}
		res, err := evaluationContext.GetValue()
		if err != nil {
			return nil, err
		}
		if optAlwaysReturnList && j.path.IsDefinite() {
			array := config.JsonProvider().CreateArray()
			err = config.JsonProvider().SetArrayIndex(array, 0, res)
			if err != nil {
				return nil, err
			}
			return array, nil
		} else {
			return res, nil
		}
	}
}

func CreateJsonpathByStringAndPredicates(jsonpath string, filters []common.Predicate) (*Jsonpath, error) {
	if jsonpath == "" {
		return nil, errors.New("json can not be null or empty")
	}
	p, err := filter.PathCompileByStringAndPredicateSlice(jsonpath, filters)
	if err != nil {
		return nil, err
	}
	return &Jsonpath{path: p}, nil
}

func compileJsonpathByStringAndPredicateSlice(jsonpath string, filters []common.Predicate) (*Jsonpath, error) {
	if jsonpath == "" {
		return nil, errors.New("json can not be null or empty")
	}
	return CreateJsonpathByStringAndPredicates(jsonpath, filters)
}

func compileJsonpath(jsonpath string, filters ...common.Predicate) (*Jsonpath, error) {
	if jsonpath == "" {
		return nil, errors.New("json can not be null or empty")
	}
	return CreateJsonpathByStringAndPredicates(jsonpath, filters)
}

func JsonpathParseString(json string) (DocumentContext, error) {
	pc := createParseContextImpl()
	return pc.parseString(json)
}

func JsonpathParseObject(json interface{}) (DocumentContext, error) {
	pc := createParseContextImpl()
	return pc.parseAny(json)
}
