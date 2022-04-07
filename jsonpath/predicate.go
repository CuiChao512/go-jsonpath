package jsonpath

import "cuichao.com/go-jsonpath/jsonpath/path"

type Predicate interface {
	Apply(ctx *PredicateContext) bool
	String() string
}

type PredicateContext interface {
	Item() interface{}

	Root() interface{}

	Configuration() *Configuration
}

type PredicateContextImpl struct {
	contextDocument   interface{}
	rootDocument      interface{}
	configuration     *Configuration
	documentPathCache map[path.Path]interface{}
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

func (pc *PredicateContextImpl) Evaluate(path2 path.Path) interface{} {
	return nil
}
