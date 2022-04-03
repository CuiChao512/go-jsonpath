package jsonpath

type Predicate interface {
	Apply(ctx *PredicateContext) bool
	String() string
}

type PredicateContext interface {
	Item() *interface{}

	Root() *interface{}

	Configures() *Configuration
}
