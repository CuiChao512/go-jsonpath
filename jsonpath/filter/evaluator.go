package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/predicate"
	"reflect"
	"strings"
)

type Evaluator interface {
	Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error)
}

type existsEvaluator struct {
}

func (*existsEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if !left.IsBooleanNode() && !right.IsBooleanNode() {
		return false, &jsonpath.JsonPathError{Message: "Failed to evaluate exists expression"}
	}
	leftNode, err := left.AsBooleanNode()
	if err != nil {
		return false, err
	}
	rightNode, err := right.AsBooleanNode()
	return leftNode.GetBoolean() == rightNode.GetBoolean(), nil
}

type notEqualsEvaluator struct {
}

func (*notEqualsEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	eqResult, err := evaluators[RelationalOperator_EQ].Evaluate(left, right, ctx)
	if err != nil {
		return false, err
	}
	return !eqResult, nil
}

type typeSafeNotEqualsEvaluator struct {
}

func (*typeSafeNotEqualsEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	tseqResult, err := evaluators[RelationalOperator_TSEQ].Evaluate(left, right, ctx)
	if err != nil {
		return false, err
	}
	return !tseqResult, nil
}

type equalsEvaluator struct{}

func (*equalsEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if left.IsJsonNode() && right.IsJsonNode() {
		leftNode, err := left.AsJsonNode()
		if err != nil {
			return false, err
		}

		rightNode, err := right.AsJsonNode()
		if err != nil {
			return false, err
		}
		return leftNode.EqualsByPredicateContext(rightNode, ctx), nil
	} else {
		return left.Equals(right), nil
	}
}

type typeSafeEqualsEvaluator struct{}

func (*typeSafeEqualsEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return false, nil
	}
	return evaluators[RelationalOperator_EQ].Evaluate(left, right, ctx)
}

type typeEvaluator struct{}

func (*typeEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	return reflect.ValueOf(right).Kind() == left.TypeOf(ctx), nil
}

type lessThanEvaluator struct{}

func (*lessThanEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if left.IsNumberNode() && right.IsNumberNode() {
		leftNode, err := left.AsNumberNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsNumberNode()
		if err != nil {
			return false, err
		}
		return leftNode.GetNumber().Cmp(*rightNode.GetNumber()) < 0, nil
	} else if left.IsStringNode() && right.IsStringNode() {
		leftNode, err := left.AsStringNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsStringNode()
		if err != nil {
			return false, err
		}
		return strings.Compare(leftNode.String(), rightNode.String()) < 0, nil
	} else if left.IsOffsetDateTimeNode() && right.IsOffsetDateTimeNode() { //workaround for issue: https://github.com/json-path/JsonPath/issues/613
		leftNode, err := left.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		return jsonpath.OffsetDateTimeCompare(leftNode.GetDate(), rightNode.GetDate()) < 0, nil
	}
	return false, nil
}

type lessThanEqualsEvaluator struct{}

func (*lessThanEqualsEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if left.IsNumberNode() && right.IsNumberNode() {
		leftNode, err := left.AsNumberNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsNumberNode()
		if err != nil {
			return false, err
		}
		return leftNode.GetNumber().Cmp(*rightNode.GetNumber()) <= 0, nil
	} else if left.IsStringNode() && right.IsStringNode() {
		leftNode, err := left.AsStringNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsStringNode()
		if err != nil {
			return false, err
		}
		return strings.Compare(leftNode.String(), rightNode.String()) <= 0, nil
	} else if left.IsOffsetDateTimeNode() && right.IsOffsetDateTimeNode() { //workaround for issue: https://github.com/json-path/JsonPath/issues/613
		leftNode, err := left.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		return jsonpath.OffsetDateTimeCompare(leftNode.GetDate(), rightNode.GetDate()) <= 0, nil
	}
	return false, nil
}

type greaterThanEvaluator struct{}

func (*greaterThanEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if left.IsNumberNode() && right.IsNumberNode() {
		leftNode, err := left.AsNumberNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsNumberNode()
		if err != nil {
			return false, err
		}
		return leftNode.GetNumber().Cmp(*rightNode.GetNumber()) > 0, nil
	} else if left.IsStringNode() && right.IsStringNode() {
		leftNode, err := left.AsStringNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsStringNode()
		if err != nil {
			return false, err
		}
		return strings.Compare(leftNode.String(), rightNode.String()) > 0, nil
	} else if left.IsOffsetDateTimeNode() && right.IsOffsetDateTimeNode() { //workaround for issue: https://github.com/json-path/JsonPath/issues/613
		leftNode, err := left.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		return jsonpath.OffsetDateTimeCompare(leftNode.GetDate(), rightNode.GetDate()) > 0, nil
	}
	return false, nil
}

type greaterThanEqualsEvaluator struct{}

func (*greaterThanEqualsEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if left.IsNumberNode() && right.IsNumberNode() {
		leftNode, err := left.AsNumberNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsNumberNode()
		if err != nil {
			return false, err
		}
		return leftNode.GetNumber().Cmp(*rightNode.GetNumber()) >= 0, nil
	} else if left.IsStringNode() && right.IsStringNode() {
		leftNode, err := left.AsStringNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsStringNode()
		if err != nil {
			return false, err
		}
		return strings.Compare(leftNode.String(), rightNode.String()) >= 0, nil
	} else if left.IsOffsetDateTimeNode() && right.IsOffsetDateTimeNode() { //workaround for issue: https://github.com/json-path/JsonPath/issues/613
		leftNode, err := left.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsOffsetDateTimeNode()
		if err != nil {
			return false, err
		}
		return jsonpath.OffsetDateTimeCompare(leftNode.GetDate(), rightNode.GetDate()) >= 0, nil
	}
	return false, nil
}

type sizeEvaluator struct{}

func (*sizeEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if !right.IsNumberNode() {
		return false, nil
	}
	rightNode, err := right.AsNumberNode()
	if err != nil {
		return false, nil
	}
	expectedSize := rightNode.GetNumber().IntPart()

	if left.IsStringNode() {
		leftNode, err := left.AsStringNode()
		if err != nil {
			return false, nil
		}
		return len(leftNode.String()) == int(expectedSize), nil
	} else if left.IsJsonNode() {

	}
	return false, nil
}

type emptyEvaluator struct {
}

func (*emptyEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	if left.IsStringNode() {
		leftNode, err := left.AsStringNode()
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsBooleanNode()
		if err != nil {
			return false, err
		}
		return leftNode.IsEmpty() == rightNode.GetBoolean(), nil
	} else if left.IsJsonNode() {
		leftNode, err := left.AsJsonNode()
		if err != nil {
			return false, err
		}

		leftIsEmpty, err := leftNode.IsEmpty(ctx)
		if err != nil {
			return false, err
		}
		rightNode, err := right.AsBooleanNode()
		if err != nil {
			return false, err
		}

		return leftIsEmpty == rightNode.GetBoolean(), nil
	}
	return false, nil
}

type inEvaluator struct{}

func (*inEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	var valueListNode *jsonpath.ValueListNode
	if right.IsJsonNode() {
		rightNode, err := right.AsJsonNode()
		if err != nil {
			return false, err
		}
		valueNode, err := rightNode.AsValueListNodeByPredicateContext(ctx)
		if err != nil {
			return false, err
		}
		if valueNode.IsUndefinedNode() {
			return false, nil
		} else {
			valueListNode, err = valueNode.AsValueListNode()
			if err != nil {
				return false, err
			}
		}
	} else {
		var err error = nil
		valueListNode, err = right.AsValueListNode()
		if err != nil {
			return false, err
		}
	}

	return valueListNode.Contains(left), nil
}

type notInEvaluator struct{}

func (*notInEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	isIn, err := evaluators[RelationalOperator_IN].Evaluate(left, right, ctx)
	if err != nil {
		return false, err
	}
	return !isIn, nil
}

type allEvaluator struct{}

func (*allEvaluator) Evaluate(left ValueNode, right ValueNode, ctx predicate.PredicateContext) (bool, error) {
	requiredValues, err := right.AsValueListNode()
	if err != nil {
		return false, err
	}
	if left.IsJsonNode() {
		jsonNode, err := left.AsJsonNode()
		if err != nil {
			return false, err
		}
		valueNode, err := jsonNode.AsValueListNodeByPredicateContext(ctx)
		if valueNode.IsValueListNode() {
			shouldContainAll, err := valueNode.AsValueListNode()
			if err != nil {
				return false, err
			}
			for _, required := range requiredValues.GetNodes() {
				if !shouldContainAll.Contains(required) {
					return false, err
				}
			}
		}
		return true, nil
	}
	return false, nil
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
	RelationalOperator_NIN      = "NIN"
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

var evaluators map[string]Evaluator

func init() {
	evaluators[RelationalOperator_EXISTS] = &existsEvaluator{}
	evaluators[RelationalOperator_NE] = &notEqualsEvaluator{}
	evaluators[RelationalOperator_TSNE] = &typeSafeNotEqualsEvaluator{}
	evaluators[RelationalOperator_EQ] = &equalsEvaluator{}
	evaluators[RelationalOperator_TSEQ] = &typeSafeEqualsEvaluator{}
	evaluators[RelationalOperator_LT] = &lessThanEvaluator{}
	evaluators[RelationalOperator_LTE] = &lessThanEqualsEvaluator{}
	evaluators[RelationalOperator_GT] = &greaterThanEvaluator{}
	evaluators[RelationalOperator_GTE] = &greaterThanEqualsEvaluator{}
	evaluators[RelationalOperator_REGEX] = &regexpEvaluator{}
	evaluators[RelationalOperator_SIZE] = &sizeEvaluator{}
	evaluators[RelationalOperator_EMPTY] = &emptyEvaluator{}
	evaluators[RelationalOperator_IN] = &inEvaluator{}
	evaluators[RelationalOperator_NIN] = &notInEvaluator{}
	evaluators[RelationalOperator_ALL] = &allEvaluator{}
	evaluators[RelationalOperator_CONTAINS] = &containsEvaluator{}
	evaluators[RelationalOperator_MATCHES] = &predicateMatchEvaluator{}
	evaluators[RelationalOperator_TYPE] = &typeEvaluator{}
	evaluators[RelationalOperator_SUBSETOF] = &subsetOfEvaluator{}
	evaluators[RelationalOperator_ANYOF] = &anyOfEvaluator{}
	evaluators[RelationalOperator_NONEOF] = &noneOfEvaluator{}
}

//LogicalOperator
const (
	LogicalOperator_AND = "&&"
	LogicalOperator_NOT = "!"
	LogicalOperator_OR  = "||"
)

type ExpressionNode interface {
	predicate.Predicate
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
	return CreateLogicalAnd(e, other)
}
func (e *LogicalExpressionNode) Or(other *LogicalExpressionNode) *LogicalExpressionNode {
	return CreateLogicalOr(e, other)
}
func (e *LogicalExpressionNode) Append(node ExpressionNode) *LogicalExpressionNode {
	e.chain = append(e.chain, node)
	return e
}

func (e *LogicalExpressionNode) GetOperator() string {
	return e.operator
}

func (e *LogicalExpressionNode) Apply(ctx PredicateContext) bool {
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
	return "(" + UtilsJoin(" "+e.operator+" ", "", chainString) + ")"
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

func CreateLogicalOr(left ExpressionNode, right ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNode(left, LogicalOperator_OR, right)
}

func CreateLogicalOrByList(operands []ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNodeByOperatorAndValues(LogicalOperator_OR, operands)
}

func CreateLogicalAnd(left ExpressionNode, right ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNode(left, LogicalOperator_AND, right)
}

func CreateLogicalAndByList(operands []ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNodeByOperatorAndValues(LogicalOperator_AND, operands)
}

func CreateLogicalNot(op ExpressionNode) *LogicalExpressionNode {
	return newLogicalExpressionNode(op, LogicalOperator_NOT, nil)
}

//RelationExpressionNode -----

type RelationExpressionNode struct {
}

func (e *RelationExpressionNode) ExpressionNodeLabel() {
	return
}
func (e *RelationExpressionNode) Apply(ctx PredicateContext) bool {
	return false
}
func (e *RelationExpressionNode) String() string {
	return "nil"
}

func CreateRelationExpressionNode(valueNode1 filter.ValueNode, operator string, valueNode2 filter.ValueNode) *RelationExpressionNode {
	return &RelationExpressionNode{}
}
