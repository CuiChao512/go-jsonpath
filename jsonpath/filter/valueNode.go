package filter

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
)

type iPatternNode interface {
	IsPatternNode() bool
	AsPatternNode() (*PatternNode, error)
}

type iPathNode interface {
	IsPathNode() bool
	AsPathNode() (*PathNode, error)
}

type iNumberNode interface {
	IsNumberNode() bool
	AsNumberNode() (*NumberNode, error)
}

type iStringNode interface {
	IsStringNode() bool
	AsStringNode() (*StringNode, error)
}

type iBooleanNode interface {
	IsBooleanNode() bool
	AsBooleanNode() (*BooleanNode, error)
}

type iPredicateNode interface {
	IsPredicateNode() bool
	AsPredicateNode() (*PredicateNode, error)
}

type iValueListNode interface {
	IsValueListNode() bool
	AsValueListNode() (*ValueListNode, error)
}

type iNullNode interface {
	IsNullNode() bool
	AsNullNode() (*NullNode, error)
}

type iUndefinedNode interface {
	IsUndefinedNode() bool
	AsUndefinedNode() (*UndefinedNode, error)
}

type iClassNode interface {
	IsClassNode() bool
	AsClassNode() (*ClassNode, error)
}

type iOffsetDateTimeNode interface {
	IsOffsetDateTimeNode() bool
	AsOffsetDateTimeNode() (*OffsetDateTimeNode, error)
}
type iJsonNode interface {
	IsJsonNode() bool
	AsJsonNode() (*JsonNode, error)
}

type ValueNode interface {
	TypeOf(ctx common.PredicateContext) reflect.Kind
	iPatternNode
	iPathNode
	iNumberNode
	iStringNode
	iBooleanNode
	iPredicateNode
	iValueListNode
	iNullNode
	iUndefinedNode
	iClassNode
	iOffsetDateTimeNode
	iJsonNode
	String() string
	Equals(o interface{}) bool
}

type defaultPatternNode struct {
}

func (n *defaultPatternNode) IsPatternNode() bool {
	return false
}

func (_ *defaultPatternNode) AsPatternNode() (*PatternNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected regexp node"}
}

type defaultPathNode struct {
}

func (n *defaultPathNode) IsPathNode() bool {
	return false
}

func (_ *defaultPathNode) AsPathNode() (*PathNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected path node"}
}

type defaultNumberNode struct {
}

func (n *defaultNumberNode) IsNumberNode() bool {
	return false
}

func (_ *defaultNumberNode) AsNumberNode() (*NumberNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected number node"}
}

type defaultStringNode struct {
}

func (n *defaultStringNode) IsStringNode() bool {
	return false
}

func (_ *defaultStringNode) AsStringNode() (*StringNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected string node"}
}

type defaultBooleanNode struct {
}

func (n *defaultBooleanNode) IsBooleanNode() bool {
	return false
}

func (_ *defaultBooleanNode) AsBooleanNode() (*BooleanNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected boolean node"}
}

type defaultPredicateNode struct {
}

func (n *defaultPredicateNode) IsPredicateNode() bool {
	return false
}

func (_ *defaultPredicateNode) AsPredicateNode() (*PredicateNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected predicate node"}
}

type defaultValueListNode struct {
}

func (n *defaultValueListNode) IsValueListNode() bool {
	return false
}

func (_ *defaultValueListNode) AsValueListNode() (*ValueListNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected value list node"}
}

type defaultNullNode struct {
}

func (n *defaultNullNode) IsNullNode() bool {
	return false
}

func (_ *defaultNullNode) AsNullNode() (*NullNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected null node"}
}

type defaultUndefinedNode struct {
}

func (n *defaultUndefinedNode) IsUndefinedNode() bool {
	return false
}

func (_ *defaultUndefinedNode) AsUndefinedNode() (*UndefinedNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected undefined node"}
}

type defaultJsonNode struct {
}

func (n *defaultJsonNode) IsJsonNode() bool {
	return false
}

func (_ *defaultJsonNode) AsJsonNode() (*JsonNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected json node"}
}

type defaultClassNode struct {
}

func (n *defaultClassNode) IsClassNode() bool {
	return false
}

func (_ *defaultClassNode) AsClassNode() (*ClassNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected class node"}
}

type defaultOffsetDateTimeNode struct {
}

func (n *defaultOffsetDateTimeNode) IsOffsetDateTimeNode() bool {
	return false
}

func (_ *defaultOffsetDateTimeNode) AsOffsetDateTimeNode() (*OffsetDateTimeNode, error) {
	return nil, &common.InvalidPathError{Message: "Expected offset date time node"}
}
