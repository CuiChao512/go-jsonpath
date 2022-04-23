package jsonpath

import (
	"errors"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type parseContext interface {
	parseString(json string) (DocumentContext, error)
	parseAny(json interface{}) (DocumentContext, error)
}

type parseContextImpl struct {
	configuration *common.Configuration
}

func (pCtx *parseContextImpl) parseString(json string) (DocumentContext, error) {
	if json == "" {
		return nil, errors.New("json string can not be empty")
	}
	obj, err := pCtx.configuration.JsonProvider().Parse(json)
	if err != nil {
		return nil, err
	}
	return CreateJsonContextByAny(obj, pCtx.configuration)
}

func (pCtx *parseContextImpl) parseAny(json interface{}) (DocumentContext, error) {
	if json == nil {
		return nil, errors.New("json object can not be nil")
	}
	return CreateJsonContextByAny(json, pCtx.configuration)
}

func createParseContextImpl() *parseContextImpl {
	return createParseContextImplByConfiguration(common.DefaultConfiguration())
}

func createParseContextImplByConfiguration(configuration *common.Configuration) *parseContextImpl {
	return &parseContextImpl{
		configuration: configuration,
	}
}
