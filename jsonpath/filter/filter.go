package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"strings"
)

type Evaluator interface {
	Evaluate(left ValueNode, right ValueNode, ctx jsonpath.PredicateContext) bool
}

type Filter interface {
	jsonpath.Predicate
	Or(other *jsonpath.Predicate) *OrFilter
	And(other *jsonpath.Predicate) *AndFilter
}

type FilterImpl struct {
}

func (filter *FilterImpl) String() string {
	return ""
}

func (filter *FilterImpl) Apply(ctx *jsonpath.PredicateContext) bool {
	return false
}

func (filter *FilterImpl) And(other *jsonpath.Predicate) *AndFilter {
	return nil
}

func (filter *FilterImpl) Or(other *jsonpath.Predicate) *OrFilter {
	return nil
}

type SingleFilter struct {
	predicate jsonpath.Predicate
}

func (filter *SingleFilter) Apply(ctx *jsonpath.PredicateContext) bool {
	return filter.predicate.Apply(ctx)
}

func (filter *SingleFilter) String() string {
	predicateString := filter.predicate.String()
	if strings.HasPrefix(predicateString, "(") {
		return "[?" + predicateString + "]"
	} else {
		return "[?(" + predicateString + ")]"
	}
}

func NewAndFilterByPredicates(predicates []*jsonpath.Predicate) *AndFilter {
	return &AndFilter{predicates: predicates}
}

func NewAndFilter(left *jsonpath.Predicate, right *jsonpath.Predicate) *AndFilter {
	predicates := []*jsonpath.Predicate{
		left, right,
	}
	return &AndFilter{predicates: predicates}
}

type AndFilter struct {
	FilterImpl
	predicates []*jsonpath.Predicate
}

func (filter *AndFilter) Apply(ctx *jsonpath.PredicateContext) bool {
	for _, predicate := range filter.predicates {
		if !(*predicate).Apply(ctx) {
			return false
		}
	}
	return true
}

func (filter *AndFilter) String() string {
	string_ := ""
	lenPredicates := len(filter.predicates)
	for i := 0; i < lenPredicates; i++ {
		p := filter.predicates[i]
		pString := (*p).String()
		if strings.HasPrefix(pString, "[?(") {
			pString = pString[3:]
		}
		string_ = string_ + pString
		if i < lenPredicates {
			string_ = string_ + "&&"
		}
	}
	return string_
}

type OrFilter struct {
}

// RelationalOperator
const (
	RelationalOperator_GTE      = ">="
	RelationalOperator_LTE      = "<="
	RelationalOperator_EQ       = "=="
	RelationalOperator_TSEQ     = "==="
	RelationalOperator_NE       = "!="
	RelationalOperator_TSNE     = "!=="
	RelationalOperator_LT       = "<"
	RelationalOperator_GT       = ">"
	RelationalOperator_REGEX    = "=~"
	RelationalOperator_MIN      = "MIN"
	RelationalOperator_IN       = "IN"
	RelationalOperator_CONTAINS = "CONTAINS"
	RelationalOperator_ALL      = "ALL"
	RelationalOperator_SIZE     = "SIZE"
	RelationalOperator_EXISTS   = "EXISTS"
	RelationalOperator_TYPE     = "TYPE"
	RelationalOperator_MATCHES  = "MATCHES"
	RelationalOperator_EMPTY    = "EMPTY"
	RelationalOperator_SUBSETOF = "SUBSETOF"
	RelationalOperator_ANYOF    = "ANYOF"
	RelationalOperator_NONEOF   = "NONEOF"
)

//LogicalOperator
const (
	LogicalOperator_AND = "&&"
	LogicalOperator_NOT = "!"
	LogicalOperator_OR  = "||"
)

type ExpressionNode interface {
	jsonpath.Predicate
	ExpressionNodeLabel()
}

//LogicalExpressionNode ----
type LogicalExpressionNode struct {
}

func (e *LogicalExpressionNode) ExpressionNodeLabel() {
	return
}
func (e *LogicalExpressionNode) Apply(ctx *jsonpath.PredicateContext) bool {
	return false
}
func (e *LogicalExpressionNode) String() string {
	return "nil"
}

//RelationExpressionNode -----

type RelationExpressionNode struct {
}

func (e *RelationExpressionNode) ExpressionNodeLabel() {
	return
}
func (e *RelationExpressionNode) Apply(ctx *jsonpath.PredicateContext) bool {
	return false
}
func (e *RelationExpressionNode) String() string {
	return "nil"
}
