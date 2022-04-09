package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"reflect"
)

type ValueNode interface {
	TypeOf(ctx jsonpath.PredicateContext) reflect.Kind
	IsPatternNode() bool
	AsPatternNode() (*PatternNode, error)
	IsPathNode() bool
	AsPathNode() (*PathNode, error)
	IsNumberNode() bool
	AsNumberNode() (*NumberNode, error)
	IsStringNode() bool
	AsStringNode() (*StringNode, error)
	IsBooleanNode() bool
	AsBooleanNode() (*BooleanNode, error)
	IsPredicateNode() bool
	AsPredicateNode() (*PredicateNode, error)
	IsValueListNode() bool
	AsValueListNode() (*ValueListNode, error)
	IsNullNode() bool
	AsNullNode() (*NullNode, error)
	IsUndefinedNode() bool
	AsUndefinedNode() (*UndefinedNode, error)
	IsClassNode() bool
	AsClassNode() (*ClassNode, error)
	IsOffsetDateTimeNode() bool
	AsOffsetDateTimeNode() (*OffsetDateTimeNode, error)
	String() string
	Equals(o interface{}) bool
}

type valueNodeDefault struct {
}

func (n *valueNodeDefault) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.Invalid
}

func (n *valueNodeDefault) IsPatternNode() bool {
	return false
}

func (_ *valueNodeDefault) AsPatternNode() (*PatternNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *valueNodeDefault) IsPathNode() bool {
	return false
}

func (_ *valueNodeDefault) AsPathNode() (*PathNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected path node"}
}

func (n *valueNodeDefault) IsNumberNode() bool {
	return false
}

func (_ *valueNodeDefault) AsNumberNode() (*NumberNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected number node"}
}

func (n *valueNodeDefault) IsStringNode() bool {
	return false
}

func (_ *valueNodeDefault) AsStringNode() (*StringNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected string node"}
}

func (n *valueNodeDefault) IsBooleanNode() bool {
	return false
}

func (_ *valueNodeDefault) AsBooleanNode() (*BooleanNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected boolean node"}
}

func (n *valueNodeDefault) IsPredicateNode() bool {
	return false
}

func (_ *valueNodeDefault) AsPredicateNode() (*PredicateNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected predicate node"}
}

func (n *valueNodeDefault) IsValueListNode() bool {
	return false
}

func (_ *valueNodeDefault) AsValueListNode() (*ValueListNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected value list node"}
}

func (n *valueNodeDefault) IsNullNode() bool {
	return false
}

func (_ *valueNodeDefault) AsNullNode() (*NullNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected null node"}
}

func (n *valueNodeDefault) IsUndefinedNode() bool {
	return false
}

func (_ *valueNodeDefault) AsUndefinedNode() (*UndefinedNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected undefined node"}
}

func (n *valueNodeDefault) IsJsonNode() bool {
	return false
}

func (_ *valueNodeDefault) AsJsonNode() (*JsonNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected json node"}
}

func (n *valueNodeDefault) IsClassNode() bool {
	return false
}

func (_ *valueNodeDefault) AsClassNode() (*ClassNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected class node"}
}

func (n *valueNodeDefault) IsOffsetDateTimeNode() bool {
	return false
}

func (_ *valueNodeDefault) AsOffsetDateTimeNode() (*OffsetDateTimeNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected offset date time node"}
}

func (n *valueNodeDefault) String() string {
	return ""
}

func (n *valueNodeDefault) Equals(o interface{}) bool {
	return false
}
