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

func (pc *PredicateContextImpl) Evaluate(path2 common.Path) (interface{}, error) {
	var result interface{}
	if path2.IsRootPath() {
		if pc.documentPathCache[path2] != nil {
			result = pc.documentPathCache[path2]
		} else {
			r, err := path2.Evaluate(pc.rootDocument, pc.rootDocument, pc.configuration)
			if err != nil {
				return nil, err
			}
			result = r.GetValue()
		}
	} else {
		r, err := path2.Evaluate(pc.rootDocument, pc.rootDocument, pc.configuration)
		if err != nil {
			return nil, err
		}
		result = r.GetValue()
	}
	return result, nil
}

func CreatePredicateContextImpl(contextDocument interface{}, rootDocument interface{}, configuration *common.Configuration, documentPathCache map[common.Path]interface{}) common.PredicateContext {
	return &PredicateContextImpl{
		contextDocument:   contextDocument,
		rootDocument:      rootDocument,
		configuration:     configuration,
		documentPathCache: documentPathCache,
	}
}
