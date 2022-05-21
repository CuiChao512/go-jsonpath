package filter

import (
	"errors"
	"fmt"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"log"
	"strconv"
	"strings"
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
	filter *common.CharacterIndex
}

func (c *Compiler) readLogicalOR() (ExpressionNode, error) {
	var ops []ExpressionNode
	op, err := c.readLogicalAND()
	if err != nil {
		return nil, err
	}
	ops = append(ops, op)
	filter := c.filter
	for {
		savepoint := filter.Position()
		if filter.HasSignificantSubSequence(LogicalOperator_OR) {
			op, err = c.readLogicalAND()
			if err != nil {
				return nil, err
			}
			ops = append(ops, op)
		} else {
			filter.SetPosition(savepoint)
			break
		}
	}

	if len(ops) == 1 {
		return ops[0], nil
	} else {
		return CreateLogicalOrByList(ops), nil
	}
}

func (c *Compiler) readLogicalAND() (ExpressionNode, error) {
	var ops []ExpressionNode
	op, err := c.readLogicalANDOperand()
	if err != nil {
		return nil, err
	}
	ops = append(ops, op)
	filter := c.filter
	for {
		savepoint := filter.Position()
		if filter.HasSignificantSubSequence(LogicalOperator_AND) {
			op, err := c.readLogicalANDOperand()
			if err != nil {
				return nil, err
			}
			ops = append(ops, op)
		} else {
			filter.SetPosition(savepoint)
			break
		}
	}

	if len(ops) == 1 {
		return ops[0], nil
	} else {
		return CreateLogicalAndByList(ops), nil
	}
}

func (c *Compiler) readLogicalANDOperand() (ExpressionNode, error) {
	filter := c.filter
	savepoint := filter.SkipBlanks().Position()
	if filter.SkipBlanks().CurrentCharIs(NOT) {
		err := filter.ReadSignificantChar(NOT)
		if err != nil {
			return nil, err
		}
		switch filter.SkipBlanks().CurrentChar() {
		case DOC_CONTEXT:
			fallthrough
		case EVAL_CONTEXT:
			filter.SetPosition(savepoint)
			break
		default:
			expressionNode, err := c.readLogicalANDOperand()
			if err != nil {
				return nil, err
			}
			return CreateLogicalNot(expressionNode), nil
		}
	}

	if filter.SkipBlanks().CurrentCharIs(OPEN_PARENTHESIS) {
		err := filter.ReadSignificantChar(OPEN_PARENTHESIS)
		if err != nil {
			return nil, err
		}
		op, err := c.readLogicalOR()
		if err != nil {
			return nil, err
		}
		err = filter.ReadSignificantChar(CLOSE_PARENTHESIS)
		if err != nil {
			return nil, err
		}
		return op, nil
	}

	return c.readExpression()
}

func (c *Compiler) readValueNode() (ValueNode, error) {
	filter := c.filter
	currentChar := filter.SkipBlanks().CurrentChar()
	println("readValueNode currentChar:", string(currentChar))
	switch currentChar {
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
			return nil, &common.InvalidPathError{Message: fmt.Sprintf("Unexpected character: %c", NOT)}
		}
	default:
		return c.readLiteral()
	}
}

func (c *Compiler) readLiteral() (ValueNode, error) {
	currentChar := c.filter.SkipBlanks().CurrentChar()
	println("readLiteral: currentChar=", string(currentChar))
	switch currentChar {
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
		return c.readJsonLiteral()
	case OPEN_ARRAY:
		return c.readJsonLiteral()
	case PATTERN:
		return c.readPattern()
	default:
		println("readNumberLiteral")
		return c.readNumberLiteral()
	}
}

func (c *Compiler) readExpression() (*RelationExpressionNode, error) {
	left, err0 := c.readValueNode()
	if err0 != nil {
		switch err0.(type) {
		case *common.InvalidPathError:
		default:
			return nil, err0
		}
	}
	filter := c.filter
	savepoint := filter.Position()
	operator, err1 := c.readRelationalOperator()

	if err0 == nil && err1 == nil {
		right, err2 := c.readValueNode()
		if err2 == nil {
			return CreateRelationExpressionNode(left, operator, right), nil
		}
	}
	filter.SetPosition(savepoint)
	pathNode, err3 := left.AsPathNode()
	if err3 != nil {
		return nil, err3
	}
	pathNode = pathNode.AsExistsCheck(pathNode.ShouldExists())
	var right *BooleanNode
	if pathNode.ShouldExists() {
		right = TRUE_NODE
	} else {
		right = FALSE_NODE
	}
	return CreateRelationExpressionNode(pathNode, RelationalOperator_EXISTS, right), nil
}

func (c *Compiler) readRelationalOperator() (string, error) {
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
	operator := filter.SubSequence(begin, filter.Position())
	if evaluators[operator] == nil {
		return "", &common.InvalidPathError{Message: "Filter operator " + operator + " is not supported!"}
	}
	return operator, nil
}

func (c *Compiler) readNullLiteral() (*NullNode, error) {
	filter := c.filter

	begin := filter.Position()

	if filter.CurrentChar() == NULL && filter.InBoundsByPosition(filter.Position()+3) {
		nullValue := filter.SubSequence(filter.Position(), filter.Position()+4)
		if "null" == nullValue {
			log.Printf("NullLiteral from %d to %d -> [%s]", begin, filter.Position()+3, nullValue)
			filter.IncrementPosition(len(nullValue))
			return CreateNullNode(), nil
		}
	}

	return nil, &common.InvalidPathError{Message: "Expected <null> value"}
}

func (c *Compiler) readJsonLiteral() (*JsonNode, error) {
	filter := c.filter

	begin := filter.Position()

	openChar := filter.CurrentChar()

	if openChar != OPEN_ARRAY && openChar != OPEN_OBJECT {
		return nil, errors.New("not a json array or object")
	}

	closeChar := CLOSE_OBJECT
	if openChar == OPEN_ARRAY {
		closeChar = CLOSE_ARRAY
	}

	closingIndex, err := filter.IndexOfMatchingCloseChar(filter.Position(), openChar, closeChar, true, false)
	if err != nil {
		return nil, err
	} else if closingIndex == -1 {
		return nil, &common.InvalidPathError{
			Message: "String not closed. Expected " + string(SINGLE_QUOTE) + " in " + filter.String(),
		}
	} else {
		filter.SetPosition(closingIndex + 1)
	}

	json := filter.SubSequence(begin, filter.Position())
	return CreateJsonNodeByString(json), err
}

func parsePatternFlags(flag rune) int {
	if flag == 'i' || flag == 'm' || flag == 's' || flag == 'U' {
		return 1
	}
	return 0
}

func (c *Compiler) endOfFlags(position int) int {
	endIndex := position
	var currentChar rune
	for c.filter.InBoundsByPosition(endIndex) {
		currentChar = c.filter.CharAt(endIndex)
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
		return nil, &common.InvalidPathError{Message: "Pattern not closed. Expected " + string(PATTERN) + " in " + filter.String()}
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
	return CreatePatternNodeByString(pattern)
}

func (c *Compiler) readStringLiteral(endChar rune) (*StringNode, error) {
	filter := c.filter
	begin := filter.Position()

	closingSingleQuoteIndex := filter.NextIndexOfUnescaped(endChar)
	if closingSingleQuoteIndex == -1 {
		return nil, &common.InvalidPathError{Message: "String literal does not have matching quotes. Expected " + string(endChar) + " in " + filter.String()}
	} else {
		filter.SetPosition(closingSingleQuoteIndex + 1)
	}
	stringLiteral := filter.SubSequence(begin, filter.Position())
	log.Printf("StringLiteral from %d to %d -> [%s]", begin, filter.Position(), stringLiteral)
	return CreateStringNode(stringLiteral, true)
}

func (c *Compiler) readNumberLiteral() (*NumberNode, error) {
	filter := c.filter
	begin := filter.Position()

	for filter.InBounds() && filter.IsNumberCharacter(filter.Position()) {
		filter.IncrementPosition(1)
	}
	numberLiteral := filter.SubSequence(begin, filter.Position())
	log.Printf("NumberLiteral from %d to %d -> [%s]", begin, filter.Position(), numberLiteral)
	return CreateNumberNodeByString(numberLiteral)
}

func (c *Compiler) readBooleanLiteral() (*BooleanNode, error) {
	filter := c.filter
	begin := filter.Position()
	end := filter.Position() + 4
	if filter.CurrentChar() == TRUE {
		end = filter.Position() + 3
	}

	if !filter.InBoundsByPosition(end) {
		return nil, &common.InvalidPathError{Message: "Expected boolean literal"}
	}
	boolString := filter.SubSequence(begin, end+1)
	if boolString != "true" && boolString != "false" {
		return nil, &common.InvalidPathError{Message: "Expected boolean literal"}
	}
	filter.IncrementPosition(len(boolString))
	log.Printf("BooleanLiteral from %d to %d -> [%s]", begin, end, boolString)
	boolValue := false
	if boolString == "true" {
		boolValue = true
	}
	return CreateBooleanNode(boolValue), nil
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
				return nil, &common.InvalidPathError{Message: "Square brackets does not match in filter " + filter.String()}
			} else {
				filter.SetPosition(closingSquareBracketIndex + 1)
			}
		}
		closingFunctionBracket := filter.CurrentChar() == CLOSE_PARENTHESIS && c.currentCharIsClosingFunctionBracket(begin)
		closingLogicalBracket := filter.CurrentChar() == CLOSE_PARENTHESIS && !closingFunctionBracket

		if !filter.InBounds() || c.isRelationalOperatorChar(filter.CurrentChar()) || filter.CurrentChar() == SPACE || closingLogicalBracket {
			break
		} else {
			filter.IncrementPosition(1)
		}
	}

	shouldExists := !(previousSignificantChar == NOT)
	path := filter.SubSequence(begin, filter.Position())
	return CreatePathNodeWithString(path, false, shouldExists)
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

func (c *Compiler) Compile() (common.Predicate, error) {
	result, err := c.readLogicalOR()
	if err != nil {
		switch err.(type) {
		case *common.InvalidPathError:
			return nil, err
		default:
			return nil, &common.InvalidPathError{Message: "Failed to parse filter: " + c.filter.String() +
				", error on position: " + strconv.Itoa(c.filter.Position()) + ", char: " + string(c.filter.CurrentChar())}
		}
	}
	filter := c.filter
	filter.SkipBlanks()
	if filter.InBounds() {
		return nil, &common.InvalidPathError{
			Message: fmt.Sprintf("Expected end of filter expression instead of: %s",
				c.filter.SubSequence(filter.Position(), filter.Length())),
		}
	}
	return result, nil
}

func CreateFilterCompiler(filterString string) (*Compiler, error) {
	compiler := &Compiler{}
	compiler.filter = common.CreateCharacterIndex(filterString)

	f := compiler.filter
	f.Trim()

	if !f.CurrentCharIs('[') || !f.LastCharIs(']') {
		return nil, &common.InvalidPathError{Message: "Filter must start with '[' and end with ']'. " + filterString}
	}

	f.IncrementPosition(1)
	f.DecrementEndPosition(1)
	f.Trim()

	if !f.CurrentCharIs('?') {
		return nil, &common.InvalidPathError{Message: "Filter must start with '[?' and end with ']'. " + filterString}
	}

	f.IncrementPosition(1)
	f.Trim()
	if !f.CurrentCharIs('(') || !f.LastCharIs(')') {
		return nil, &common.InvalidPathError{Message: "Filter must start with '[?(' and end with ')]'. " + filterString}
	}

	return compiler, nil
}

func Compile(filterString string) (*CompiledFilter, error) {
	compiler, err := CreateFilterCompiler(filterString)
	if err != nil {
		return nil, err
	}
	compiledFilter, err := compiler.Compile()
	if err != nil {
		return nil, err
	}
	return &CompiledFilter{predicate: compiledFilter}, nil
}

type CompiledFilter struct {
	predicate common.Predicate
}

func (cf *CompiledFilter) Apply(ctx common.PredicateContext) (bool, error) {
	return cf.predicate.Apply(ctx)
}

func (cf *CompiledFilter) String() string {
	predicateString := cf.predicate.String()
	if strings.HasPrefix(predicateString, "(") {
		return "[?" + predicateString + "]"
	} else {
		return "[?(" + predicateString + ")]"
	}
}
