package path

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
)

type PredicateContextImpl struct {
	contextDocument   interface{}
	rootDocument      interface{}
	configuration     *common.Configuration
	documentPathCache map[common.Path]interface{}
}

func (pc *PredicateContextImpl) Item() interface{} {
	return pc.contextDocument
}

func (pc *PredicateContextImpl) Root() interface{} {
	return pc.rootDocument
}

func (pc *PredicateContextImpl) Configuration() *common.Configuration {
	return pc.configuration
}

func (pc *PredicateContextImpl) Evaluate(path2 common.Path) interface{} {
	return nil
}

func CreatePredicateContextImpl(contextDocument interface{}, rootDocument interface{}, configuration *common.Configuration, documentPathCache map[common.Path]interface{}) common.PredicateContext {
	return &PredicateContextImpl{
		contextDocument:   contextDocument,
		rootDocument:      rootDocument,
		configuration:     configuration,
		documentPathCache: documentPathCache,
	}
}
