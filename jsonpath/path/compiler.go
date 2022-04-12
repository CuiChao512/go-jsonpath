package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/filter"
	"strconv"
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
	filterStack []filter.PredicateNode
	path        *jsonpath.CharacterIndex
}

func (cp *Compiler) readWhitespace() {
	for cp.path.InBounds() {
		c := cp.path.CurrentChar()
		if cp.isWhitespace(c) {
			break
		}
		cp.path.IncrementPosition(1)
	}
}

func (*Compiler) isWhitespace(c rune) bool {
	return c == SPACE || c == TAB || c == LF || c == CR
}

func (*Compiler) isPathContext(c rune) bool {
	return c == DOC_CONTEXT || c == EVAL_CONTEXT
}

func (cp *Compiler) readContextToken() (*RootPathToken, error) {
	cp.readWhitespace()

	if !cp.isPathContext(cp.path.CurrentChar()) {
		return nil, fail("Path must start with '$' or '@'")
	}

	pathToken := CreateRootPathToken(cp.path.CurrentChar())

	if cp.path.CurrentIsTail() {
		return pathToken, nil
	}

	cp.path.IncrementPosition(1)

	if cp.path.CurrentChar() != PERIOD && cp.path.CurrentChar() != OPEN_SQUARE_BRACKET {
		return nil, fail("Illegal character at position " + strconv.FormatInt(int64(cp.path.Position()), 10) + " expected '.' or '['")
	}

	appender := pathToken.GetPathTokenAppender()

	cp.readNextToken(appender)
	return pathToken, nil
}

func (cp *Compiler) readNextToken(appender TokenAppender) bool {
	switch cp.path.CurrentChar() {
	case OPEN_SQUARE_BRACKET:
	case PERIOD:
	case WILDCARD:
	default:
	}
	return true
}

func fail(message string) *jsonpath.InvalidPathError {
	return &jsonpath.InvalidPathError{Message: message}
}

func Compile(pathString string, filters ...jsonpath.Predicate) (Path, error) {
	ci := jsonpath.NewCharacterIndex(pathString)

	if ci.CharAt(0) != DOC_CONTEXT && ci.CharAt(0) != EVAL_CONTEXT {
		ci = jsonpath.NewCharacterIndex("$." + pathString)
	}
	ci.Trim()

	if ci.LastCharIs('.') {
		return nil, fail("Path must not end with a '.' or '..'")
	}

	//filterStack := filters[:]

	return &CompiledPath{}, nil
}
