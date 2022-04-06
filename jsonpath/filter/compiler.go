package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"fmt"
	"log"
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

func (c *Compiler) readLogicalOR() jsonpath.Predicate {
	var ops []jsonpath.Predicate
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

func (c *Compiler) readLogicalAND() jsonpath.Predicate {
	var ops []jsonpath.Predicate
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

func (c *Compiler) readLogicalANDOperand() jsonpath.Predicate {
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

func (c *Compiler) readValueNode() (ValueNode, *jsonpath.InvalidPathError) {
	filter := c.filter
	switch filter.SkipBlanks().CurrentChar() {
	case DOC_CONTEXT:
		return c.readPath()
	case EVAL_CONTEXT:
		return c.readPath()
	case NOT:
		filter.IncrementPosition(1)
		switch filter.SkipBlanks().CurrentChar() {
		case DOC_CONTEXT:
			return c.readPath()
		case EVAL_CONTEXT:
			return c.readPath()
		default:
			return nil, &jsonpath.InvalidPathError{Message: fmt.Sprintf("Unexpected character: %c", NOT)}
		}
	default:
		return c.readLiteral(),nil
	}
}

func (c *Compiler) readLiteral() (ValueNode,*jsonpath.InvalidPathError){
	switch c.filter.SkipBlanks().CurrentChar() {
	case SINGLE_QUOTE:
		return c.readStringLiteral(SINGLE_QUOTE)
	case DOUBLE_QUOTE:
		return c.readStringLiteral(DOUBLE_QUOTE)
	case TRUE:
		return c.readBooleanLiteral()
	case FALSE:
		return c.readBooleanLiteral()
	case MINUS:
		return c.readNumberLiteral()
	case NULL:
		return c.readNullLiteral()
	case OPEN_OBJECT:
		return readJsonLiteral()
	case OPEN_ARRAY:
		return readJsonLiteral()
	case PATTERN:
		return readPattern()
	default:
		return readNumberLiteral()
	}
}

func (c *Compiler) readExpression() jsonpath.Predicate {
	left, err0 := c.readValueNode()
	filter := c.filter
	savepoint := filter.Position()
	operator := c.readRelationalOperator()
	right, err1 := c.readValueNode()
	if err0 == nil && err1 == nil{
		return NewRelationExpressionNode(left, operator, right)
	} else {
		filter.SetPosition(savepoint)
		pathNode, _ := left.AsPathNode()
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

func (c *Compiler) readRelationalOperator() string{
	filter := c.filter
	begin := filter.SkipBlanks().Position()

	if c.isRelationalOperatorChar(filter.CurrentChar()){
		for ;filter.InBounds() && c.isRelationalOperatorChar(filter.CurrentChar());{
			filter.IncrementPosition(1)
		}
	} else {
		for ;filter.InBounds() && filter.CurrentChar() != SPACE;{
			filter.IncrementPosition(1)
		}
	}
	return filter.SubSequence(begin, filter.Position())
}

func (c *Compiler) readNullLiteral() (*NullNode,*jsonpath.InvalidPathError){
	filter := c.filter

	begin := filter.Position()

	if filter.CurrentChar() == NULL && filter.InBoundsByPosition(filter.Position() + 3){
		nullValue := filter.SubSequence(filter.Position(), filter.Position() + 4)
		if "null" == nullValue{
			log.Printf("NullLiteral from %d to %d -> [%s]", begin, filter.Position()+3, nullValue)
			filter.IncrementPosition(len(nullValue))
			return NewNullNode(),nil
		}
	}

	return nil,&jsonpath.InvalidPathError{Message: "Expected <null> value"}
}

func (c *Compiler) readJsonLiteral() (*JsonNode,*jsonpath.InvalidPathError){
	filter := c.filter

	begin := filter.Position()

	openChar := filter.CurrentChar()

	//TODO: assert openChar == OPEN_ARRAY || openChar == OPEN_OBJECT;

	closeChar := CLOSE_OBJECT
	if openChar == OPEN_ARRAY{
		closeChar = CLOSE_ARRAY
	}

	closingIndex,err := filter.IndexOfMatchingCloseChar(filter.Position(),openChar,closeChar,true,false)
	if err != nil{
		return nil, err
	} else if closingIndex == - 1{
		return nil, &jsonpath.InvalidPathError{
			Message: "String not closed. Expected " + string(SINGLE_QUOTE) + " in " + filter.String(),
		}
	} else {
		filter.SetPosition(closingIndex + 1)
	}

	json := filter.SubSequence(begin, filter.Position())
	return NewJsonNode(json),err
}

func (c *Compiler) readPath() (*PathNode, *jsonpath.InvalidPathError) {
	filter := c.filter
	previousSignificantChar := filter.PreviousSignificantChar()
	begin := filter.Position()
	filter.IncrementPosition(1)

	for filter.InBounds() {
		if filter.CurrentChar() == OPEN_SQUARE_BRACKET {
			closingSquareBracketIndex, err := filter.IndexOfMatchingCloseChar(filter.Position(), OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, true, false)
			if err != nil {
				return nil, err
			} else if closingSquareBracketIndex == -1 {
				return nil, &jsonpath.InvalidPathError{Message: "Square brackets does not match in filter " + filter.String()}
			} else {
				filter.SetPosition(closingSquareBracketIndex + 1)
			}

			closingFunctionBracket := filter.CurrentChar() == CLOSE_PARENTHESIS && c.currentCharIsClosingFunctionBracket(begin)
			closingLogicalBracket := filter.CurrentChar() == CLOSE_PARENTHESIS && !closingFunctionBracket

			if !filter.InBounds() || c.isRelationalOperatorChar(filter.CurrentChar()) || filter.CurrentChar() == SPACE || closingLogicalBracket {
				break
			} else {
				filter.IncrementPosition(1)
			}
		}
	}

	shouldExists := !(previousSignificantChar == NOT)
	path := filter.SubSequence(begin, filter.Position())
	return NewPathNodeWithString(path, false, shouldExists), nil
}

func (c *Compiler) currentCharIsClosingFunctionBracket(lowerBound int) bool {
	filter := c.filter
	if filter.CurrentChar() != CLOSE_PARENTHESIS{
		return false
	}

	idx := filter.IndexOfPreviousSignificantChar()
	if idx == -1 || filter.CharAt(idx) != OPEN_PARENTHESIS{
		return false
	}

	idx--

	for ;filter.InBoundsByPosition(idx) && idx > lowerBound{
		if filter.CharAt(idx) == PERIOD {
			return true
		}
		idx--
	}
	return false
}

func (*Compiler) isLogicalOperatorChar(c rune) bool {
	return c == AND || c == OR
}

func (*Compiler) isRelationalOperatorChar(c rune) bool {
	return c == LT || c == GT || c == EQ || c == TILDE || c == NOT
}

func (c *Compiler) Compile() (jsonpath.Predicate, *jsonpath.InvalidPathError) {
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
