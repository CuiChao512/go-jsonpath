package common

type Path interface {
	Evaluate(document interface{}, rootDocument interface{}, configuration *Configuration) (EvaluationContext, error)
	EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *Configuration, forUpdate bool) EvaluationContext
	String() string
	IsDefinite() bool
	IsFunctionPath() bool
	IsRootPath() bool
}

type PathRef interface {
	GetAccessor() interface{}
	Set(newVal interface{}, configuration *Configuration) error
	Convert(mapFunction MapFunction, configuration *Configuration) error
	Delete(configuration *Configuration) error
	Add(newVal interface{}, configuration *Configuration) error
	Put(key string, newVal interface{}, configuration *Configuration) error
	RenameKey(oldKeyName string, newKeyName string, configuration *Configuration) error
	CompareTo(o PathRef) int
}

type EvaluationContext interface {
	Configuration() *Configuration
	RootDocument() interface{}
	GetValue() interface{}
	GetValueUnwrap(unwrap bool) interface{}
}

type Predicate interface {
	Apply(ctx PredicateContext) (bool, error)
	String() string
}

type PredicateContext interface {
	Item() interface{}

	Root() interface{}

	Configuration() *Configuration
}
