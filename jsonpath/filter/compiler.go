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

func (c *Compiler) readLogicalOR() ExpressionNode {
	var ops []ExpressionNode
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
		return NewLogicalOrByList(ops)
	}
}

func (c *Compiler) readLogicalAND() ExpressionNode {
	var ops []ExpressionNode
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
		return NewLogicalAndByList(ops)
	}
}

func (c *Compiler) readLogicalANDOperand() ExpressionNode {
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
			return NewLogicalNot(c.readLogicalANDOperand())
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

func (c *Compiler) readValueNode() (ValueNode, error) {
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
		return c.readLiteral()
	}
}

func (c *Compiler) readLiteral() (ValueNode, error) {
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
		return c.readNumberLiteral(), nil
	case NULL:
		return c.readNullLiteral()
	case OPEN_OBJECT:
		return c.readJsonLiteral()
	case OPEN_ARRAY:
		return c.readJsonLiteral()
	case PATTERN:
		return c.readPattern()
	default:
		return c.readNumberLiteral(), nil
	}
}

func (c *Compiler) readExpression() *RelationExpressionNode {
	left, err0 := c.readValueNode()
	filter := c.filter
	savepoint := filter.Position()
	operator := c.readRelationalOperator()
	right, err1 := c.readValueNode()
	if err0 == nil && err1 == nil {
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
		return NewRelationExpressionNode(left, RelationalOperator_EXISTS, right)
	}
}

func (c *Compiler) readRelationalOperator() string {
	filter := c.filter
	begin := filter.SkipBlanks().Position()

	if c.isRelationalOperatorChar(filter.CurrentChar()) {
		for filter.InBounds() && c.isRelationalOperatorChar(filter.CurrentChar()) {
			filter.IncrementPosition(1)
		}
	} else {
		for filter.InBounds() && filter.CurrentChar() != SPACE {
			filter.IncrementPosition(1)
		}
	}
	return filter.SubSequence(begin, filter.Position())
}

func (c *Compiler) readNullLiteral() (*NullNode, error) {
	filter := c.filter

	begin := filter.Position()

	if filter.CurrentChar() == NULL && filter.InBoundsByPosition(filter.Position()+3) {
		nullValue := filter.SubSequence(filter.Position(), filter.Position()+4)
		if "null" == nullValue {
			log.Printf("NullLiteral from %d to %d -> [%s]", begin, filter.Position()+3, nullValue)
			filter.IncrementPosition(len(nullValue))
			return NewNullNode(), nil
		}
	}

	return nil, &jsonpath.InvalidPathError{Message: "Expected <null> value"}
}

func (c *Compiler) readJsonLiteral() (*JsonNode, error) {
	filter := c.filter

	begin := filter.Position()

	openChar := filter.CurrentChar()

	//TODO: assert openChar == OPEN_ARRAY || openChar == OPEN_OBJECT;

	closeChar := CLOSE_OBJECT
	if openChar == OPEN_ARRAY {
		closeChar = CLOSE_ARRAY
	}

	closingIndex, err := filter.IndexOfMatchingCloseChar(filter.Position(), openChar, closeChar, true, false)
	if err != nil {
		return nil, err
	} else if closingIndex == -1 {
		return nil, &jsonpath.InvalidPathError{
			Message: "String not closed. Expected " + string(SINGLE_QUOTE) + " in " + filter.String(),
		}
	} else {
		filter.SetPosition(closingIndex + 1)
	}

	json := filter.SubSequence(begin, filter.Position())
	return NewJsonNode(json), err
}

func parsePatternFlags(c [1]rune) int {
	//TODO: PatternFlag.parseFlags
	return 0
}

func (c *Compiler) endOfFlags(position int) int {
	endIndex := position
	var currentChar [1]rune
	for c.filter.InBoundsByPosition(endIndex) {
		currentChar[0] = c.filter.CharAt(endIndex)
		if parsePatternFlags(currentChar) > 0 {
			endIndex++
			continue
		}
		break
	}
	return endIndex
}

func (c *Compiler) readPattern() (*PatternNode, error) {
	filter := c.filter
	begin := filter.Position()
	closingIndex := filter.NextIndexOfUnescaped(PATTERN)

	if closingIndex == -1 {
		return nil, &jsonpath.InvalidPathError{Message: "Pattern not closed. Expected " + string(PATTERN) + " in " + filter.String()}
	} else {
		if filter.InBoundsByPosition(closingIndex + 1) {
			endFlagsIndex := c.endOfFlags(closingIndex + 1)
			if endFlagsIndex > closingIndex {
				flags := filter.SubSequence(closingIndex+1, endFlagsIndex)
				closingIndex += len(flags)
			}
		}
		filter.SetPosition(closingIndex + 1)
	}
	pattern := filter.SubSequence(begin, filter.Position())
	log.Printf("PatternNode from %d to %d -> [%s]", begin, filter.Position(), pattern)
	return NewPatternNode(pattern), nil
}

func (c *Compiler) readStringLiteral(endChar rune) (*StringNode, error) {
	filter := c.filter
	begin := filter.Position()

	closingSingleQuoteIndex := filter.NextIndexOfUnescaped(endChar)
	if closingSingleQuoteIndex == -1 {
		return nil, &jsonpath.InvalidPathError{Message: "String literal does not have matching quotes. Expected " + string(endChar) + " in " + filter.String()}
	} else {
		filter.SetPosition(closingSingleQuoteIndex + 1)
	}
	stringLiteral := filter.SubSequence(begin, filter.Position())
	log.Printf("StringLiteral from %d to %d -> [%s]", begin, filter.Position(), stringLiteral)
	return NewStringNode(stringLiteral, true), nil
}

func (c *Compiler) readNumberLiteral() *NumberNode {
	filter := c.filter
	begin := filter.Position()

	for filter.InBounds() && filter.IsNumberCharacter(filter.Position()) {
		filter.IncrementPosition(1)
	}
	numberLiteral := filter.SubSequence(begin, filter.Position())
	log.Printf("NumberLiteral from %d to %d -> [%s]", begin, filter.Position(), numberLiteral)
	return NewNumberNodeByString(numberLiteral)
}

func (c *Compiler) readBooleanLiteral() (*BooleanNode, error) {
	filter := c.filter
	begin := filter.Position()
	end := filter.Position() + 4
	if filter.CurrentChar() == TRUE {
		end = filter.Position() + 3
	}

	if !filter.InBoundsByPosition(end) {
		return nil, &jsonpath.InvalidPathError{Message: "Expected boolean literal"}
	}
	boolString := filter.SubSequence(begin, end+1)
	if boolString != "true" && boolString != "false" {
		return nil, &jsonpath.InvalidPathError{Message: "Expected boolean literal"}
	}
	filter.IncrementPosition(len(boolString))
	log.Printf("BooleanLiteral from %d to %d -> [%s]", begin, end, boolString)
	boolValue := false
	if boolString == "true" {
		boolValue = true
	}
	return NewBooleanNode(boolValue), nil
}

func (c *Compiler) readPath() (*PathNode, error) {
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
	if filter.CurrentChar() != CLOSE_PARENTHESIS {
		return false
	}

	idx := filter.IndexOfPreviousSignificantChar()
	if idx == -1 || filter.CharAt(idx) != OPEN_PARENTHESIS {
		return false
	}

	idx--

	for filter.InBoundsByPosition(idx) && idx > lowerBound {
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

func (c *Compiler) Compile() (jsonpath.Predicate, error) {
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
