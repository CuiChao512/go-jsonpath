package jsonpath

import "strings"

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
	return strings.IndexRune(ci.charSequence).charAt(idx)
}

func (ci *CharacterIndex) CurrentChar() {
	return charSequence.charAt(position)
}

func (ci *CharacterIndex) CurrentCharIs(c rune) bool {
	return (charSequence.charAt(position) == c)
}

func (ci *CharacterIndex) LastCharIs(char c) bool {
	return charSequence.charAt(endPosition) == c
}

func (ci *CharacterIndex) NextCharIs(char c) bool {
	return inBounds(position+1) && (charSequence.charAt(position+1) == c)
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
