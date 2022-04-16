package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/predicate"
	"reflect"
)

type ValueNode interface {
	TypeOf(ctx predicate.PredicateContext) reflect.Kind
	IsPatternNode() bool
	AsPatternNode() (*jsonpath.PatternNode, error)
	IsPathNode() bool
	AsPathNode() (*jsonpath.PathNode, error)
	IsNumberNode() bool
	AsNumberNode() (*jsonpath.NumberNode, error)
	IsStringNode() bool
	AsStringNode() (*jsonpath.StringNode, error)
	IsBooleanNode() bool
	AsBooleanNode() (*jsonpath.BooleanNode, error)
	IsPredicateNode() bool
	AsPredicateNode() (*jsonpath.PredicateNode, error)
	IsValueListNode() bool
	AsValueListNode() (*jsonpath.ValueListNode, error)
	IsNullNode() bool
	AsNullNode() (*jsonpath.NullNode, error)
	IsUndefinedNode() bool
	AsUndefinedNode() (*jsonpath.UndefinedNode, error)
	IsClassNode() bool
	AsClassNode() (*jsonpath.ClassNode, error)
	IsOffsetDateTimeNode() bool
	AsOffsetDateTimeNode() (*jsonpath.OffsetDateTimeNode, error)
	IsJsonNode() bool
	AsJsonNode() (*jsonpath.JsonNode, error)
	String() string
	Equals(o interface{}) bool
}

type ValueNodeDefault struct {
}

func (n *ValueNodeDefault) TypeOf(ctx predicate.PredicateContext) reflect.Kind {
	return reflect.Invalid
}

func (n *ValueNodeDefault) IsPatternNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsPatternNode() (*jsonpath.PatternNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected regexp node"}
}

func (n *ValueNodeDefault) IsPathNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsPathNode() (*jsonpath.PathNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected path node"}
}

func (n *ValueNodeDefault) IsNumberNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsNumberNode() (*jsonpath.NumberNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected number node"}
}

func (n *ValueNodeDefault) IsStringNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsStringNode() (*jsonpath.StringNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected string node"}
}

func (n *ValueNodeDefault) IsBooleanNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsBooleanNode() (*jsonpath.BooleanNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected boolean node"}
}

func (n *ValueNodeDefault) IsPredicateNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsPredicateNode() (*jsonpath.PredicateNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected predicate node"}
}

func (n *ValueNodeDefault) IsValueListNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsValueListNode() (*jsonpath.ValueListNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected value list node"}
}

func (n *ValueNodeDefault) IsNullNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsNullNode() (*jsonpath.NullNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected null node"}
}

func (n *ValueNodeDefault) IsUndefinedNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsUndefinedNode() (*jsonpath.UndefinedNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected undefined node"}
}

func (n *ValueNodeDefault) IsJsonNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsJsonNode() (*jsonpath.JsonNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected json node"}
}

func (n *ValueNodeDefault) IsClassNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsClassNode() (*jsonpath.ClassNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected class node"}
}

func (n *ValueNodeDefault) IsOffsetDateTimeNode() bool {
	return false
}

func (_ *ValueNodeDefault) AsOffsetDateTimeNode() (*jsonpath.OffsetDateTimeNode, error) {
	return nil, &jsonpath.InvalidPathError{Message: "Expected offset date time node"}
}

func (n *ValueNodeDefault) String() string {
	return ""
}

func (n *ValueNodeDefault) Equals(o interface{}) bool {
	return false
}
