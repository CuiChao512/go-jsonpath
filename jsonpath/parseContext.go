package jsonpath

import (
	"errors"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type ParseContext interface {
	ParseString(json string) (DocumentContext, error)
	ParseAny(json interface{}) (DocumentContext, error)
}

type ParseContextImpl struct {
	configuration *common.Configuration
}

func (pCtx *ParseContextImpl) ParseString(json string) (DocumentContext, error) {
	if json == "" {
		return nil, errors.New("json string can not be empty")
	}
	obj, err := pCtx.configuration.JsonProvider().Parse(json)
	if err != nil {
		return nil, err
	}
	return CreateJsonContextByAny(obj, pCtx.configuration)
}

func (pCtx *ParseContextImpl) ParseAny(json interface{}) (DocumentContext, error) {
	if json == nil {
		return nil, errors.New("json object can not be nil")
	}
	return CreateJsonContextByAny(json, pCtx.configuration)
}

func createParseContextImpl() *ParseContextImpl {
	return CreateParseContextImplByConfiguration(common.DefaultConfiguration())
}

func CreateParseContextImplByConfiguration(configuration *common.Configuration) *ParseContextImpl {
	return &ParseContextImpl{
		configuration: configuration,
	}
}
