package common

type PredicateContextImpl struct {
	contextDocument   interface{}
	rootDocument      interface{}
	configuration     *Configuration
	documentPathCache map[Path]interface{}
}

func (pc *PredicateContextImpl) Item() interface{} {
	return pc.contextDocument
}

func (pc *PredicateContextImpl) Root() interface{} {
	return pc.rootDocument
}

func (pc *PredicateContextImpl) Configuration() *Configuration {
	return pc.configuration
}

func (pc *PredicateContextImpl) Evaluate(path2 Path) (interface{}, error) {
	var result interface{}
	if path2.IsRootPath() {
		if pc.documentPathCache[path2] != nil {
			result = pc.documentPathCache[path2]
		} else {
			r, err := path2.Evaluate(pc.rootDocument, pc.rootDocument, pc.configuration)
			if err != nil {
				return nil, err
			}
			result, err = r.GetValue()
			if err != nil {
				return nil, err
			}
		}
	} else {
		r, err := path2.Evaluate(pc.rootDocument, pc.rootDocument, pc.configuration)
		if err != nil {
			return nil, err
		}
		result, err = r.GetValue()
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func CreatePredicateContextImpl(contextDocument interface{}, rootDocument interface{}, configuration *Configuration, documentPathCache map[Path]interface{}) PredicateContext {
	return &PredicateContextImpl{
		contextDocument:   contextDocument,
		rootDocument:      rootDocument,
		configuration:     configuration,
		documentPathCache: documentPathCache,
	}
}
