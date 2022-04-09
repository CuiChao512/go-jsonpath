package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/filter"
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

type CompiledPath struct {
	filterStack []filter.PredicateNode
	path        *jsonpath.CharacterIndex
}

func (cp *CompiledPath) readWhitespace() {
	for cp.path.InBounds() {
		c := cp.path.CurrentChar()
		if cp.isWhitespace(c) {
			break
		}
		cp.path.IncrementPosition(1)
	}
}

func (*CompiledPath) isWhitespace(c rune) bool {
	return c == SPACE || c == TAB || c == LF || c == CR
}

func (cp *CompiledPath) Evaluate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration) (jsonpath.EvaluationContext, error) {
	return nil, nil
}

func (cp *CompiledPath) EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration, forUpdate bool) jsonpath.EvaluationContext {
	return nil
}

func (cp *CompiledPath) String() string {
	return ""
}

func (cp *CompiledPath) IsDefinite() bool {
	return false
}

func (cp *CompiledPath) IsFunctionPath() bool {
	return false
}

func (cp *CompiledPath) IsRootPath() bool {
	return false
}

func Compile(pathString string) Path {
	return &CompiledPath{}
}
