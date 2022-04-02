package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
)

type ValueNode interface {
	TypeOf(ctx jsonpath.PredicateContext) string
	IsPatternNode() bool
	AsPatternNode() (*PatternNode, *jsonpath.InvalidPathError)
	IsPathNode() bool
	AsPathNode() (*PathNode, *jsonpath.InvalidPathError)
	IsNumberNode() bool
	AsNumberNode() (*NumberNode, *jsonpath.InvalidPathError)
	IsStringNode() bool
	AsStringNode() (*StringNode, *jsonpath.InvalidPathError)
	IsBooleanNode() bool
	AsBooleanNode() (*BooleanNode, *jsonpath.InvalidPathError)
	IsPredicateNode() bool
	AsPredicateNode() (*PredicateNode, *jsonpath.InvalidPathError)
	IsValueListNode() bool
	AsValueListNode() (*ValueListNode, *jsonpath.InvalidPathError)
	IsNullNode() bool
	AsNullNode() (*NullNode, *jsonpath.InvalidPathError)
	IsUndefinedNode() bool
	AsUndefinedNode() (*UndefinedNode, *jsonpath.InvalidPathError)
	IsClassNode() bool
	AsClassNode() (*ClassNode, *jsonpath.InvalidPathError)
	IsOffsetDateTimeNode() bool
	AsOffsetDateTimeNode() (*OffsetDateTimeNode, *jsonpath.InvalidPathError)
}

type valueNodeDefault struct {
}

func (n *valueNodeDefault) TypeOf(ctx jsonpath.PredicateContext) string {
	return "void"
}

func (n *valueNodeDefault) IsPatternNode() bool {
	return false
}

func (_ *valueNodeDefault) AsPatternNode() (*PatternNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *valueNodeDefault) IsPathNode() bool {
	return false
}

func (_ *valueNodeDefault) AsPathNode() (*PathNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected path node"}
}

func (n *valueNodeDefault) IsNumberNode() bool {
	return false
}

func (_ *valueNodeDefault) AsNumberNode() (*NumberNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected number node"}
}

func (n *valueNodeDefault) IsStringNode() bool {
	return false
}

func (_ *valueNodeDefault) AsStringNode() (*StringNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected string node"}
}

func (n *valueNodeDefault) IsBooleanNode() bool {
	return false
}

func (_ *valueNodeDefault) AsBooleanNode() (*BooleanNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected boolean node"}
}

func (n *valueNodeDefault) IsPredicateNode() bool {
	return false
}

func (_ *valueNodeDefault) AsPredicateNode() (*PredicateNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *valueNodeDefault) IsValueListNode() bool {
	return false
}

func (_ *valueNodeDefault) AsValueListNode() (*ValueListNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *valueNodeDefault) IsNullNode() bool {
	return false
}

func (_ *valueNodeDefault) AsNullNode() (*NullNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *valueNodeDefault) IsUndefinedNode() bool {
	return false
}

func (_ *valueNodeDefault) AsUndefinedNode() (*UndefinedNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *valueNodeDefault) IsClassNode() bool {
	return false
}

func (_ *valueNodeDefault) AsClassNode() (*ClassNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *valueNodeDefault) IsOffsetDateTimeNode() bool {
	return false
}

func (_ *valueNodeDefault) AsOffsetDateTimeNode() (*OffsetDateTimeNode, *jsonpath.InvalidPathError) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}
