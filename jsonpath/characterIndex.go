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
	position     int64
	endPosition  int64
}

func (ci *CharacterIndex) Length() int64 {
	return ci.endPosition + 1
}
func (ci *CharacterIndex) CharAt(idx int64) rune {
	return []rune(ci.charSequence)[idx]
}

func (ci *CharacterIndex) CurrentChar() rune {
	return []rune(ci.charSequence)[ci.position]
}

func (ci *CharacterIndex) CurrentCharIs(c rune) bool {
	return []rune(ci.charSequence)[ci.position] == c
}

func (ci *CharacterIndex) LastCharIs(c rune) bool {
	return []rune(ci.charSequence)[ci.endPosition] == c
}

func (ci *CharacterIndex) InBounds(idx int64) bool {
	return (idx >= 0) && (idx <= ci.endPosition)
}

func (ci *CharacterIndex) NextCharIs(c rune) bool {
	return ci.InBounds(ci.position+1) && ([]rune(ci.charSequence)[ci.position+1] == c)
}

func (ci *CharacterIndex) IncrementPosition(charCount int64) int64 {
	return ci.SetPosition(ci.position + charCount)
}

func (ci *CharacterIndex) DecrementEndPosition(charCount int64) int64 {
	return ci.SetEndPosition(ci.endPosition - charCount)
}

func (ci *CharacterIndex) SetPosition(newPosition int64) int64 {
	//position = min(newPosition, charSequence.length() - 1);
	ci.position = newPosition
	return ci.position
}

func (ci *CharacterIndex) SetEndPosition(newPosition int64) int64 {
	ci.endPosition = newPosition
	return ci.endPosition
}

func (ci *CharacterIndex) Position() int64 {
	return ci.position
}
