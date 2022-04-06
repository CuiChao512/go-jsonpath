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
	chain    []ExpressionNode
	operator string
}

func (e *LogicalExpressionNode) ExpressionNodeLabel() {
	return
}
func (e *LogicalExpressionNode) And(other *LogicalExpressionNode) *LogicalExpressionNode {
	return NewLogicalAnd(e, other)
}
func (e *LogicalExpressionNode) Or(other *LogicalExpressionNode) *LogicalExpressionNode {
	return NewLogicalOr(e, other)
}
func (e *LogicalExpressionNode) Append(node ExpressionNode) *LogicalExpressionNode {
	e.chain = append(e.chain, node)
	return e
}

func (e *LogicalExpressionNode) GetOperator() string {
	return e.operator
}

func (e *LogicalExpressionNode) Apply(ctx *jsonpath.PredicateContext) bool {
	if e.operator == LogicalOperator_OR {
		for _, expression := range e.chain {
			if expression.Apply(ctx) {
				return true
			}
		}
		return false
	} else if e.operator == LogicalOperator_AND {
		for _, expression := range e.chain {
			if !expression.Apply(ctx) {
				return false
			}
		}
		return true
	} else {
		expression := e.chain[0]
		return !expression.Apply(ctx)
	}
}
func (e *LogicalExpressionNode) String() string {
	var chainString []string
	for _, e := range e.chain {
		chainString = append(chainString, e.String())
	}
	return "(" + jsonpath.UtilsJoin(" "+e.operator+" ", "", chainString) + ")"
}

func newLogicalExpressionNode(left ExpressionNode, operator string, right ExpressionNode) *LogicalExpressionNode {
	var chain []ExpressionNode
	chain[0] = left
	chain[1] = right
	return &LogicalExpressionNode{
		chain:    chain,
		operator: operator,
	}
}

func newLogicalExpressionNodeByOperatorAndValues(operator string, values []ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		chain:    values,
		operator: operator,
	}
}

func NewLogicalOr(left ExpressionNode, right ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNode(left, LogicalOperator_OR, right)
}

func NewLogicalOrByList(operands []ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNodeByOperatorAndValues(LogicalOperator_OR, operands)
}

func NewLogicalAnd(left ExpressionNode, right ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNode(left, LogicalOperator_AND, right)
}

func NewLogicalAndByList(operands []ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNodeByOperatorAndValues(LogicalOperator_AND, operands)
}

func NewLogicalNot(op ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNode(op, LogicalOperator_NOT, nil)
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

func NewRelationExpressionNode(valueNode1 ValueNode, operator string, valueNode2 ValueNode) *RelationExpressionNode {
	return &RelationExpressionNode{}
}
