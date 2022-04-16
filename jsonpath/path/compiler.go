package path

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/filter"
	"cuichao.com/go-jsonpath/jsonpath/function"
	"cuichao.com/go-jsonpath/jsonpath/predicate"
	"strconv"
	"strings"
)

const (
	DOC_CONTEXT  = '$'
	EVAL_CONTEXT = '@'

	OPEN_SQUARE_BRACKET  = '['
	CLOSE_SQUARE_BRACKET = ']'
	OPEN_PARENTHESIS     = '('
	CLOSE_PARENTHESIS    = ')'
	OPEN_BRACE           = '{'
	CLOSE_BRACE          = '}'

	WILDCARD     = '*'
	PERIOD       = '.'
	SPACE        = ' '
	TAB          = '\t'
	CR           = '\r'
	LF           = '\n'
	BEGIN_FILTER = '?'
	COMMA        = ','
	SPLIT        = ':'
	MINUS        = '-'
	SINGLE_QUOTE = '\''
	DOUBLE_QUOTE = '"'
)

type Compiler struct {
	filterStack []predicate.Predicate
	path        *common.CharacterIndex
}

func (c *Compiler) readWhitespace() {
	for c.path.InBounds() {
		char := c.path.CurrentChar()
		if c.isWhitespace(char) {
			break
		}
		c.path.IncrementPosition(1)
	}
}

func (*Compiler) isWhitespace(c rune) bool {
	return c == SPACE || c == TAB || c == LF || c == CR
}

func (*Compiler) isPathContext(c rune) bool {
	return c == DOC_CONTEXT || c == EVAL_CONTEXT
}

func (c *Compiler) readContextToken() (*RootPathToken, error) {
	c.readWhitespace()

	if !c.isPathContext(c.path.CurrentChar()) {
		return nil, fail("Path must start with '$' or '@'")
	}

	pathToken := CreateRootPathToken(c.path.CurrentChar())

	if c.path.CurrentIsTail() {
		return pathToken, nil
	}

	c.path.IncrementPosition(1)

	if c.path.CurrentChar() != PERIOD && c.path.CurrentChar() != OPEN_SQUARE_BRACKET {
		return nil, fail("Illegal character at position " + strconv.FormatInt(int64(c.path.Position()), 10) + " expected '.' or '['")
	}

	appender := pathToken.GetPathTokenAppender()

	_, err := c.readNextToken(appender)
	if err != nil {
		return nil, err
	}
	return pathToken, nil
}

func (c *Compiler) readNextToken(appender TokenAppender) (bool, error) {
	switch c.path.CurrentChar() {
	case OPEN_SQUARE_BRACKET:
		readResult, err := c.readBracketPropertyToken(appender)
		errMsg := "Could not parse token starting at position " + strconv.Itoa(c.path.Position()) + ". Expected ?, ', 0-9, * "
		if err != nil || !readResult {
			return false, fail(errMsg)
		}

		readResult, err = c.readArrayToken(appender)
		if err != nil || !readResult {
			return false, fail(errMsg)
		}
		readResult, err = c.readWildCardToken(appender)
		if err != nil || !readResult {
			return false, fail(errMsg)
		}
		readResult, err = c.readFilterToken(appender)
		if err != nil || !readResult {
			return false, fail(errMsg)
		}
		readResult, err = c.readPlaceholderToken(appender)
		if err != nil || !readResult {
			return false, fail(errMsg)
		}
		return true, nil
	case PERIOD:
		readResult, err := c.readDotToken(appender)
		if err != nil {
			return false, err
		}
		if !readResult {
			return false, fail("Could not parse token starting at position " + strconv.Itoa(c.path.Position()))
		}
	case WILDCARD:
		readResult, err := c.readWildCardToken(appender)
		if err != nil {
			return false, err
		}
		if !readResult {
			return false, fail("Could not parse token starting at position " + strconv.Itoa(c.path.Position()))
		}
	default:
		readResult, err := c.readPropertyOrFunctionToken(appender)
		if err != nil {
			return false, err
		}
		if !readResult {
			return false, fail("Could not parse token starting at position " + strconv.Itoa(c.path.Position()))
		}
	}
	return true, nil
}

func (c *Compiler) readDotToken(appender TokenAppender) (bool, error) {
	if c.path.CurrentCharIs(PERIOD) && c.path.NextCharIs(PERIOD) {
		appender.AppendPathToken(CreateScanPathToken())
		c.path.IncrementPosition(2)
	} else if !c.path.HasMoreCharacters() {
		return false, &common.InvalidPathError{Message: "Path must not end with a '."}
	} else {
		c.path.IncrementPosition(1)
	}

	if c.path.CurrentCharIs(PERIOD) {
		return false, &common.InvalidPathError{Message: "Character '.' on position " + strconv.Itoa(c.path.Position()) + " is not valid."}
	}

	return c.readNextToken(appender)
}

func (c *Compiler) readPropertyOrFunctionToken(appender TokenAppender) (bool, error) {
	path := c.path
	if path.CurrentCharIs(OPEN_SQUARE_BRACKET) || path.CurrentCharIs(WILDCARD) || path.CurrentCharIs(PERIOD) || path.CurrentCharIs(SPACE) {
		return false, nil
	}
	startPosition := path.Position()
	readPosition := startPosition
	endPosition := 0

	isFunction := false

	for path.InBoundsByPosition(readPosition) {
		char := path.CharAt(readPosition)
		if char == SPACE {
			return false, &common.InvalidPathError{Message: "Use bracket notion ['my prop'] if your property contains blank characters. position: " + strconv.Itoa(path.Position())}
		} else if char == PERIOD || char == OPEN_SQUARE_BRACKET {
			endPosition = readPosition
			break
		} else if char == OPEN_PARENTHESIS {
			isFunction = true
			endPosition = readPosition
			break
		}
		readPosition++
	}
	if endPosition == 0 {
		endPosition = path.Length()
	}

	var functionParameters []*function.Parameter
	if isFunction {
		parenthesisCount := 1
		for i := readPosition + 1; i < path.Length(); i++ {
			if path.CharAt(i) == CLOSE_PARENTHESIS {
				parenthesisCount--
			} else if path.CharAt(i) == OPEN_PARENTHESIS {
				parenthesisCount++
			}
			if parenthesisCount == 0 {
				break
			}
		}

		if parenthesisCount != 0 {
			functionName := path.SubSequence(startPosition, endPosition)
			return false, &common.InvalidPathError{Message: "Arguments to function: '" + functionName + "' are not closed properly."}
		}

		if path.InBoundsByPosition(readPosition + 1) {
			// read the next token to determine if we have a simple no-args function call
			char := path.CharAt(readPosition + 1)
			if char != CLOSE_PARENTHESIS {
				path.SetPosition(endPosition + 1)
				// parse the arguments of the function - arguments that are inner queries or JSON document(s)
				functionName := path.SubSequence(startPosition, endPosition)
				var err error = nil
				functionParameters, err = c.parseFunctionParameters(functionName)
				if err != nil {
					return false, err
				}
			} else {
				path.SetPosition(readPosition + 1)
			}
		} else {
			path.SetPosition(readPosition)
		}
	} else {
		path.SetPosition(endPosition)
	}

	property := path.SubSequence(startPosition, endPosition)
	if isFunction {
		appender.AppendPathToken(CreateFunctionPathToken(property, functionParameters))
	} else {
		appender.AppendPathToken(CreatePropertyPathToken([]string{property}, string(SINGLE_QUOTE)))
	}
	readResult, err := c.readNextToken(appender)
	if err != nil {
		return false, err
	}
	return path.CurrentIsTail() || readResult, nil
}

func (c *Compiler) parseFunctionParameters(funcName string) ([]*function.Parameter, error) {
	paramType := function.JSON

	paramTypeUpdated := false

	// Parenthesis starts at 1 since we're marking the start of a function call, the close paren will denote the
	// last parameter boundary

	groupParen, groupBracket, groupBrace, groupQuote := 1, 0, 0, 0

	path := c.path
	endOfStream := false
	priorChar := rune(0)
	var parameters []*function.Parameter
	parameter := ""
	for path.InBounds() && !endOfStream {
		char := path.CurrentChar()
		path.IncrementPosition(1)

		// we're at the start of the stream, and don't know what type of parameter we have
		if !paramTypeUpdated {
			if c.isWhitespace(char) {
				continue
			}

			if char == OPEN_BRACE || common.UtilsCharIsDigit(char) || DOUBLE_QUOTE == char {
				paramType = function.JSON
			} else if c.isPathContext(char) {
				paramType = function.PATH // read until we reach a terminating comma and we've reset grouping to zero
			}
			paramTypeUpdated = true
		}

		switch char {
		case DOUBLE_QUOTE:
			if priorChar != '\\' && groupQuote > 0 {
				groupQuote--
			} else {
				groupQuote++
			}
		case OPEN_PARENTHESIS:
			groupParen++
		case OPEN_BRACE:
			groupBrace++
		case OPEN_SQUARE_BRACKET:
			groupBracket++
		case CLOSE_BRACE:
			if 0 == groupBrace {
				return nil, &common.InvalidPathError{Message: "Unexpected close brace '}' at character position: " + strconv.Itoa(path.Position())}
			}
			groupBrace--
		case CLOSE_SQUARE_BRACKET:
			if 0 == groupBracket {
				return nil, &common.InvalidPathError{Message: "Unexpected close bracket ']' at character position: " + strconv.Itoa(path.Position())}
			}
			groupBracket--

		// In either the close paren case where we have zero paren groups left, capture the parameter, or where
		// we've encountered a COMMA do the same
		case CLOSE_PARENTHESIS:
			groupParen--
			//CS304 Issue link: https://github.com/json-path/JsonPath/issues/620
			if 0 > groupParen || priorChar == '(' {
				parameter += string(char)
			}
		case COMMA:
			// In this state we've reach the end of a function parameter and we can pass along the parameter string
			// to the parser
			if 0 == groupQuote && 0 == groupBrace && 0 == groupBracket && ((0 == groupParen && CLOSE_PARENTHESIS == char) || 1 == groupParen) {
				endOfStream = 0 == groupParen

				if paramTypeUpdated {
					var param *function.Parameter = nil
					switch paramType {
					case function.JSON:
						// parse the json and set the value
						param = function.CreateJsonParameter(parameter)
					case function.PATH:
						var predicates []predicate.Predicate
						compiler := createPathCompiler(common.CreateCharacterIndex(parameter), &predicates)
						compiledPath, err := compiler.compile()
						if err != nil {
							return nil, err
						}
						param = function.CreatePathParameter(compiledPath)
					}
					if paramTypeUpdated {
						parameters = append(parameters, param)
					}
					parameter = ""
					paramTypeUpdated = false
				}
			}
		}

		if !paramTypeUpdated && !(char == COMMA && 0 == groupBrace && 0 == groupBracket && 1 == groupParen) {
			parameter += string(char)
		}
		priorChar = char
	}
	if 0 != groupBrace || 0 != groupParen || 0 != groupBracket {
		return nil, &common.InvalidPathError{Message: "Arguments to function: '" + funcName + "' are not closed properly."}
	}
	return parameters, nil
}

func (c *Compiler) compile() (Path, error) {
	root, err := c.readContextToken()
	if err != nil {
		return nil, err
	}

	return &CompiledPath{root: root, isRootPath: root.GetPathFragment() == "$"}, nil
}

func (c *Compiler) readPlaceholderToken(appender TokenAppender) (bool, error) {
	path := c.path
	if !path.CurrentCharIs(OPEN_SQUARE_BRACKET) {
		return false, nil
	}
	questionMarkIndex := path.IndexOfNextSignificantChar(BEGIN_FILTER)
	if questionMarkIndex == -1 {
		return false, nil
	}
	nextSignificantChar := path.NextSignificantCharFromStartPosition(questionMarkIndex)
	if nextSignificantChar != CLOSE_SQUARE_BRACKET && nextSignificantChar != COMMA {
		return false, nil
	}

	expressionBeginIndex := path.Position() + 1
	expressionEndIndex := path.NextIndexOfFromStartPosition(expressionBeginIndex, CLOSE_SQUARE_BRACKET)

	if expressionEndIndex == -1 {
		return false, nil
	}

	expression := path.SubSequence(expressionBeginIndex, expressionEndIndex)
	tokens := strings.Split(expression, ",")

	if len(c.filterStack) < len(tokens) {
		return false, &common.InvalidPathError{Message: "Not enough predicates supplied for filter [" + expression + "] at position " + strconv.Itoa(path.Position())}
	}

	var predicates []predicate.Predicate
	for _, token := range tokens {
		if token != "" {
			token = strings.TrimSpace(token)
		}
		if "?" != token {
			return false, &common.InvalidPathError{Message: "Expected '?' but found " + token}
		}
		predicates = append(predicates, c.filterStack[len(c.filterStack)-1])
		c.filterStack = c.filterStack[0 : len(c.filterStack)-1]
	}

	appender.AppendPathToken(CreatePredicatePathToken(predicates))

	path.SetPosition(expressionEndIndex + 1)

	readResult, err := c.readNextToken(appender)
	if err != nil {
		return false, err
	}
	return path.CurrentIsTail() || readResult, nil
}

func (c *Compiler) readFilterToken(appender TokenAppender) (bool, error) {
	path := c.path
	if !path.CurrentCharIs(OPEN_SQUARE_BRACKET) && !path.NextSignificantCharIs(BEGIN_FILTER) {
		return false, nil
	}

	openStatementBracketIndex := path.Position()
	questionMarkIndex := path.IndexOfNextSignificantChar(BEGIN_FILTER)
	if questionMarkIndex == -1 {
		return false, nil
	}
	openBracketIndex := path.IndexOfNextSignificantCharFromStartPosition(questionMarkIndex, OPEN_PARENTHESIS)
	if openBracketIndex == -1 {
		return false, nil
	}
	closeBracketIndex, err := path.IndexOfClosingBracket(openBracketIndex, true, true)
	if err != nil {
		return false, err
	}
	if closeBracketIndex == -1 {
		return false, nil
	}
	if !path.NextSignificantCharIsFromStartPosition(closeBracketIndex, CLOSE_SQUARE_BRACKET) {
		return false, nil
	}
	closeStatementBracketIndex := path.IndexOfNextSignificantCharFromStartPosition(closeBracketIndex, CLOSE_SQUARE_BRACKET)

	criteria := path.SubSequence(openStatementBracketIndex, closeStatementBracketIndex+1)

	predicate0, e := filter.Compile(criteria)
	if e != nil {
		return false, nil
	}
	appender.AppendPathToken(CreatePredicatePathToken([]predicate.Predicate{predicate0}))

	path.SetPosition(closeStatementBracketIndex + 1)
	readResult, e := c.readNextToken(appender)
	if e != nil {
		return false, e
	}
	return path.CurrentIsTail() || readResult, nil
}

func (c *Compiler) readWildCardToken(appender TokenAppender) (bool, error) {
	path := c.path
	inBracket := path.CurrentCharIs(OPEN_SQUARE_BRACKET)

	if inBracket && !path.NextSignificantCharIs(WILDCARD) {
		return false, nil
	}
	if !path.CurrentCharIs(WILDCARD) && path.IsOutOfBounds(path.Position()+1) {
		return false, nil
	}
	if inBracket {
		wildCardIndex := path.IndexOfNextSignificantChar(WILDCARD)
		if !path.NextSignificantCharIsFromStartPosition(wildCardIndex, CLOSE_SQUARE_BRACKET) {
			offset := wildCardIndex + 1
			return false, &common.InvalidPathError{Message: "Expected wildcard token to end with ']' on position " + strconv.Itoa(offset)}
		}
		bracketCloseIndex := path.IndexOfNextSignificantCharFromStartPosition(wildCardIndex, CLOSE_SQUARE_BRACKET)
		path.SetPosition(bracketCloseIndex + 1)
	} else {
		path.IncrementPosition(1)
	}

	appender.AppendPathToken(CreateWildcardPathToken())
	readResult, e := c.readNextToken(appender)
	if e != nil {
		return false, e
	}
	return path.CurrentIsTail() || readResult, nil
}

func (c *Compiler) readArrayToken(appender TokenAppender) (bool, error) {
	path := c.path
	if !path.CurrentCharIs(OPEN_SQUARE_BRACKET) {
		return false, nil
	}
	nextSignificantChar := path.NextSignificantChar()
	if !common.UtilsCharIsDigit(nextSignificantChar) && nextSignificantChar != MINUS && nextSignificantChar != SPLIT {
		return false, nil
	}

	expressionBeginIndex := path.Position() + 1
	expressionEndIndex := path.NextIndexOfFromStartPosition(expressionBeginIndex, CLOSE_SQUARE_BRACKET)

	if expressionEndIndex == -1 {
		return false, nil
	}

	expression := strings.TrimSpace(path.SubSequence(expressionBeginIndex, expressionEndIndex))

	if "*" == expression {
		return false, nil
	}

	//check valid chars
	for i := 0; i < len(expression); i++ {
		char := []rune(expression)[i]
		if !common.UtilsCharIsDigit(char) && char != COMMA && char != MINUS && char != SPLIT && char != SPACE {
			return false, nil
		}
	}

	isSliceOperation := strings.Contains(expression, ":")

	if isSliceOperation {
		arraySliceOperation, err := ParseArraySliceOperation(expression)
		if err != nil {
			return false, nil
		}
		appender.AppendPathToken(CreateArraySlicePathToken(arraySliceOperation))
	} else {
		arrayIndexOperation, err := ParseArrayIndexOperation(expression)
		if err != nil {
			return false, nil
		}
		appender.AppendPathToken(CreateArrayIndexPathToken(arrayIndexOperation))
	}

	path.SetPosition(expressionEndIndex + 1)
	readResult, e := c.readNextToken(appender)
	if e != nil {
		return false, e
	}
	return path.CurrentIsTail() || readResult, nil

}

func (c *Compiler) readBracketPropertyToken(appender TokenAppender) (bool, error) {
	path := c.path
	if !path.CurrentCharIs(OPEN_SQUARE_BRACKET) {
		return false, nil
	}
	potentialStringDelimiter := path.NextSignificantChar()
	if potentialStringDelimiter != SINGLE_QUOTE && potentialStringDelimiter != DOUBLE_QUOTE {
		return false, nil
	}

	var properties = make([]string, 0)

	startPosition := path.Position() + 1
	readPosition := startPosition
	endPosition := 0
	inProperty := false
	inEscape := false
	lastSignificantWasComma := false

	for path.InBoundsByPosition(readPosition) {
		char := path.CharAt(readPosition)

		if inEscape {
			inEscape = false
		} else if '\\' == char {
			inEscape = true
		} else if char == CLOSE_SQUARE_BRACKET && !inProperty {
			if lastSignificantWasComma {
				return false, fail("Found empty property at index " + strconv.Itoa(readPosition))
			}
			break
		} else if char == potentialStringDelimiter {
			if inProperty {
				nextSignificantChar := path.NextSignificantCharFromStartPosition(readPosition)
				if nextSignificantChar != CLOSE_SQUARE_BRACKET && nextSignificantChar != COMMA {
					return false, fail("Property must be separated by comma or Property must be terminated close square bracket at index " + strconv.Itoa(readPosition))
				}
				endPosition = readPosition
				prop := path.SubSequence(startPosition, endPosition)
				property, err := common.UtilsStringUnescape(prop)
				if err != nil {
					return false, err
				}
				properties = append(properties, property)
				inProperty = false
			} else {
				startPosition = readPosition + 1
				inProperty = true
				lastSignificantWasComma = false
			}
		} else if char == COMMA && !inProperty {
			if lastSignificantWasComma {
				return false, fail("Found empty property at index " + strconv.Itoa(readPosition))
			}
			lastSignificantWasComma = true
		}
		readPosition++
	}

	if inProperty {
		return false, fail("Property has not been closed - missing closing " + string(potentialStringDelimiter))
	}

	endBracketIndex := path.IndexOfNextSignificantCharFromStartPosition(endPosition, CLOSE_SQUARE_BRACKET) + 1

	path.SetPosition(endBracketIndex)

	appender.AppendPathToken(CreatePropertyPathToken(properties, string(potentialStringDelimiter)))

	readResult, e := c.readNextToken(appender)
	if e != nil {
		return false, e
	}
	return path.CurrentIsTail() || readResult, nil

}

func fail(message string) *common.InvalidPathError {
	return &common.InvalidPathError{Message: message}
}

func createPathCompiler(path *common.CharacterIndex, filterStack *[]predicate.Predicate) *Compiler {
	return &Compiler{path: path, filterStack: *filterStack}
}

func Compile(pathString string, filters ...predicate.Predicate) (Path, error) {
	ci := common.CreateCharacterIndex(pathString)

	if ci.CharAt(0) != DOC_CONTEXT && ci.CharAt(0) != EVAL_CONTEXT {
		ci = common.CreateCharacterIndex("$." + pathString)
	}
	ci.Trim()

	if ci.LastCharIs('.') {
		return nil, fail("Path must not end with a '.' or '..'")
	}

	filterStack := filters[:]

	return createPathCompiler(ci, &filterStack).compile()
}
