package filter

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/path"
	"reflect"
)

type ValueNode interface {
	TypeOf(ctx path.PredicateContext) reflect.Kind
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
	IsJsonNode() bool
	AsJsonNode() (*JsonNode, error)
	String() string
	Equals(o interface{}) bool
}

type ValueNodeDefault struct {
}

func (n *ValueNodeDefault) TypeOf(ctx path.PredicateContext) reflect.Kind {
	return reflect.Invalid
}

func (n *ValueNodeDefault) IsPatternNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsPatternNode() (*PatternNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected regexp node"}
}

func (n *ValueNodeDefault) IsPathNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsPathNode() (*PathNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected path node"}
}

func (n *ValueNodeDefault) IsNumberNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsNumberNode() (*NumberNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected number node"}
}

func (n *ValueNodeDefault) IsStringNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsStringNode() (*StringNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected string node"}
}

func (n *ValueNodeDefault) IsBooleanNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsBooleanNode() (*BooleanNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected boolean node"}
}

func (n *ValueNodeDefault) IsPredicateNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsPredicateNode() (*PredicateNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected predicate node"}
}

func (n *ValueNodeDefault) IsValueListNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsValueListNode() (*ValueListNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected value list node"}
}

func (n *ValueNodeDefault) IsNullNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsNullNode() (*NullNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected null node"}
}

func (n *ValueNodeDefault) IsUndefinedNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsUndefinedNode() (*UndefinedNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected undefined node"}
}

func (n *ValueNodeDefault) IsJsonNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsJsonNode() (*JsonNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected json node"}
}

func (n *ValueNodeDefault) IsClassNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsClassNode() (*ClassNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected class node"}
}

func (n *ValueNodeDefault) IsOffsetDateTimeNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsOffsetDateTimeNode() (*OffsetDateTimeNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected offset date time node"}
}

func (n *ValueNodeDefault) String() string {
	return ""
}

func (n *ValueNodeDefault) Equals(o interface{}) bool {
	return false
}
