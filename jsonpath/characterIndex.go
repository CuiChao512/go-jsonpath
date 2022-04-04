package jsonpath

const (
	OPEN_PARENTHESIS     rune = '('
	CLOSE_PARENTHESIS    rune = ')'
	CLOSE_SQUARE_BRACKET rune = ']'
	SPACE                rune = ' '
	ESCAPE               rune = '\\'
	SINGLE_QUOTE         rune = '\''
	DOUBLE_QUOTE         rune = '"'
	MINUS                rune = '-'
	PERIOD               rune = '.'
	REGEX                rune = '/'
	SCI_E                rune = 'E'
	SCI_e                rune = 'e'
)

type CharacterIndex struct {
	charSequence string
	position     int
	endPosition  int
}
