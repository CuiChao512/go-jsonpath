package filter

type ValueNode interface {
	IsPatternNode() bool
	AsPatternNode() PatternNode
	IsPathNode() bool
	AsPathNode() PathNode
	IsNumberNode() bool
	AsNumberNode() NumberNode
	IsStringNode() bool
	AsStringNode() StringNode
	IsBooleanNode() bool
	AsBooleanNode() BooleanNode
	IsPredicateNode() bool
	AsPredicateNode() PredicateNode
	IsValueListNode() bool
	AsValueListNode() ValueListNode
	IsNullNode() bool
	AsNullNode() NullNode
	IsUndefinedNode() bool
	AsUndefinedNode() UndefinedNode
	IsClassNode() bool
	AsClassNode() ClassNode
	IsOffsetDateTimeNode() bool
	AsOffsetDateTimeNode() OffsetDateTimeNode
}
