package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"fmt"
)

const (
	DOC_CONTEXT          = '$'
	EVAL_CONTEXT         = '@'
	OPEN_SQUARE_BRACKET  = '['
	CLOSE_SQUARE_BRACKET = ']'
	OPEN_PARENTHESIS     = '('
	CLOSE_PARENTHESIS    = ')'
	OPEN_OBJECT          = '{'
	CLOSE_OBJECT         = '}'
	OPEN_ARRAY           = '['
	CLOSE_ARRAY          = ']'

	SINGLE_QUOTE = '\''
	DOUBLE_QUOTE = '"'

	SPACE  = ' '
	PERIOD = '.'

	AND = '&'
	OR  = '|'

	MINUS       = '-'
	LT          = '<'
	GT          = '>'
	EQ          = '='
	TILDE       = '~'
	TRUE        = 't'
	FALSE       = 'f'
	NULL        = 'n'
	NOT         = '!'
	PATTERN     = '/'
	IGNORE_CASE = 'i'
)

type Compiler struct {
	filter *jsonpath.CharacterIndex
}

func (c *Compiler) readLogicalOR() *jsonpath.Predicate {
	var ops []*jsonpath.Predicate
	ops = append(ops, c.readLogicalAND())
	filter := c.filter
	for {
		savepoint := filter.Position()
		if filter.HasSignificantSubSequence(LogicalOperator_OR) {
			ops = append(ops, c.readLogicalAND())
		} else {
			filter.SetPosition(savepoint)
			break
		}
	}

	if len(ops) == 1 {
		return ops[0]
	} else {
		return CreateLogicalExpressionOr(ops)
	}
}

func (c *Compiler) readLogicalAND() *jsonpath.Predicate {
	var ops []*jsonpath.Predicate
	ops = append(ops, c.readLogicalANDOperand())
	filter := *c.filter
	for {
		savepoint := filter.Position()
		if filter.HasSignificantSubSequence(LogicalOperator_AND) {
			ops = append(ops, c.readLogicalANDOperand())
		} else {
			filter.SetPosition(savepoint)
			break
		}
	}

	if len(ops) == 1 {
		return ops[0]
	} else {
		return CreateLogicalExpressionAnd(ops)
	}
}

func (c *Compiler) readLogicalANDOperand() *jsonpath.Predicate {
	filter := c.filter
	savepoint := filter.SkipBlanks().Position()
	if filter.SkipBlanks().CurrentCharIs(NOT) {
		filter.ReadSignificantChar(NOT)
		switch filter.SkipBlanks().CurrentChar() {
		case DOC_CONTEXT:
			fallthrough
		case EVAL_CONTEXT:
			filter.SetPosition(savepoint)
			break
		default:
			return CreateLogicalExpressionNot(c.readLogicalANDOperand())
		}
	}

	if filter.SkipBlanks().CurrentCharIs(OPEN_PARENTHESIS) {
		filter.ReadSignificantChar(OPEN_PARENTHESIS)
		op := c.readLogicalOR()
		filter.ReadSignificantChar(CLOSE_PARENTHESIS)
		return op
	}

	return c.readExpression()
}

func (c *Compiler) readValueNode() (*ValueNode, *jsonpath.InvalidPathError) {
	filter := c.filter
	switch filter.SkipBlanks().CurrentChar() {
	case DOC_CONTEXT:
		return c.readPath(), nil
	case EVAL_CONTEXT:
		return c.readPath(), nil
	case NOT:
		filter.IncrementPosition(1)
		switch filter.SkipBlanks().CurrentChar() {
		case DOC_CONTEXT:
			return c.readPath(), nil
		case EVAL_CONTEXT:
			return c.readPath(), nil
		default:
			return nil, &jsonpath.InvalidPathError{Message: fmt.Sprintf("Unexpected character: %c", NOT)}
		}
	}
}

func (c *Compiler) readExpression() *jsonpath.Predicate {
	left, err := c.readValueNode()
	filter := c.filter
	savepoint := filter.Position()
	operator := c.readRelationalOperator()
	right := readValueNode()
	if err == nil {
		return NewRelationExpressionNode(left, operator, right)
	} else {
		pathNode, _ := (*left).AsPathNode()
		pathNode = pathNode.AsExistsCheck(pathNode.ShouldExists())
		var right *BooleanNode
		if pathNode.ShouldExists() {
			right = TRUE_NODE
		} else {
			right = FALSE_NODE
		}
		return NewRelationalExpressionNode(left, RelationalOperator_EXISTS, right)
	}
}

func (c *Compiler) Compile() (*jsonpath.Predicate, *jsonpath.InvalidPathError) {
	result := c.readLogicalOR()
	filter := c.filter
	filter.SkipBlanks()
	if filter.InBounds() {
		return nil, &jsonpath.InvalidPathError{
			Message: fmt.Sprintf("Expected end of filter expression instead of: %s",
				c.filter.SubSequence(filter.Position(), filter.Length())),
		}
	}
	return result, nil
}

func FilterCompile(filterString string) CompiledFilter {
	return CompiledFilter{}
}

type CompiledFilter struct {
}
