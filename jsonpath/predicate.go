package jsonpath

type Predicate interface {
	Apply(ctx *PredicateContext)
}

type PredicateContext interface {
	Item() *interface{}

	Root() *interface{}

	Configures() *Configuration
}
