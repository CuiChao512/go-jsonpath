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

func (ci *CharacterIndex) Length() int {
	return ci.endPosition + 1
}
func (ci *CharacterIndex) CharAt(idx int) rune {
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

func (ci *CharacterIndex) InBoundsByPosition(idx int) bool {
	return (idx >= 0) && (idx <= ci.endPosition)
}

func (ci *CharacterIndex) NextCharIs(c rune) bool {
	return ci.InBoundsByPosition(ci.position+1) && ([]rune(ci.charSequence)[ci.position+1] == c)
}

func (ci *CharacterIndex) IncrementPosition(charCount int) int {
	return ci.SetPosition(ci.position + charCount)
}

func (ci *CharacterIndex) DecrementEndPosition(charCount int) int {
	return ci.SetEndPosition(ci.endPosition - charCount)
}

func (ci *CharacterIndex) SetPosition(newPosition int) int {
	//position = min(newPosition, charSequence.length() - 1);
	ci.position = newPosition
	return ci.position
}

func (ci *CharacterIndex) SetEndPosition(newPosition int) int {
	ci.endPosition = newPosition
	return ci.endPosition
}

func (ci *CharacterIndex) Position() int {
	return ci.position
}

func (ci *CharacterIndex) IndexOfClosingSquareBracket(startPosition int) int {
	readPosition := startPosition
	for ci.InBoundsByPosition(readPosition) {
		if ci.CharAt(readPosition) == CLOSE_SQUARE_BRACKET {
			return readPosition
		}
		readPosition++
	}
	return -1
}

func (ci *CharacterIndex) IndexOfMatchingCloseChar(startPosition int, openChar rune,
	closeChar rune, skipStrings bool, skipRegex bool) (int, *InvalidPathError) {
	if ci.CharAt(startPosition) != openChar {
		return -1, &InvalidPathError{Message: "Expected " + string(openChar) + " but found " + string(ci.CharAt(startPosition))}
	}

	opened := 1
	readPosition := startPosition + 1

	for ci.InBoundsByPosition(readPosition) {
		if skipStrings {
			quoteChar := ci.CharAt(readPosition)
			if quoteChar == SINGLE_QUOTE || quoteChar == DOUBLE_QUOTE {
				readPosition = ci.NextIndexOfUnescapedFromStartPosition(readPosition, quoteChar)
				if readPosition == -1 {
					return -1, &InvalidPathError{Message: "Could not find matching close quote for " + string(quoteChar) + " when parsing : " + ci.charSequence}
				}
				readPosition++
			}
		}
		if skipRegex {
			if ci.CharAt(readPosition) == REGEX {
				readPosition = ci.NextIndexOfUnescapedFromStartPosition(readPosition, REGEX)
				if readPosition == -1 {
					return -1, &InvalidPathError{Message: "Could not find matching close quote for " + string(REGEX) + " when parsing regex in : " + ci.charSequence}
				}
				readPosition++
			}
		}
		if ci.CharAt(readPosition) == openChar {
			opened++
		}
		if ci.CharAt(readPosition) == closeChar {
			opened--
			if opened == 0 {
				return readPosition, nil
			}
		}
		readPosition++
	}
	return -1, nil
}

func (ci *CharacterIndex) IndexOfClosingBracket(startPosition int, skipStrings bool, skipRegex bool) (int, *InvalidPathError) {
	return ci.IndexOfMatchingCloseChar(startPosition, OPEN_PARENTHESIS, CLOSE_PARENTHESIS, skipStrings, skipRegex)
}

func (ci *CharacterIndex) IndexOfNextSignificantChar(c rune) int {
	return ci.IndexOfNextSignificantCharFromStartPosition(ci.position, c)
}

func (ci *CharacterIndex) IndexOfNextSignificantCharFromStartPosition(startPosition int, c rune) int {
	readPosition := startPosition + 1
	for !ci.IsOutOfBounds(readPosition) && ci.CharAt(readPosition) == SPACE {
		readPosition++
	}

	if ci.CharAt(readPosition) == c {
		return readPosition
	} else {
		return -1
	}
}

func (ci *CharacterIndex) NextIndexOf(c rune) int {
	return ci.NextIndexOfFromStartPosition(ci.position+1, c)
}

func (ci *CharacterIndex) NextIndexOfFromStartPosition(startPosition int, c rune) int {
	readPosition := startPosition
	for ci.IsOutOfBounds(readPosition) {
		if ci.CharAt(readPosition) == c {
			return readPosition
		}
		readPosition++
	}
	return -1
}

func (ci *CharacterIndex) NextIndexOfUnescaped(c rune) int {
	return ci.NextIndexOfUnescapedFromStartPosition(ci.position, c)
}

func (ci *CharacterIndex) NextIndexOfUnescapedFromStartPosition(startPosition int, c rune) int {
	readPosition := startPosition + 1
	inEscape := false

	for ci.IsOutOfBounds(readPosition) {
		if inEscape {
			inEscape = false
		} else if '\\' == ci.CharAt(readPosition) {
			inEscape = true
		} else if c == ci.CharAt(readPosition) {
			return readPosition
		}
		readPosition++
	}
	return -1
}

func (ci *CharacterIndex) CharAtOr(position int, defaultChar rune) rune {
	if !ci.InBoundsByPosition(position) {
		return defaultChar
	} else {
		return ci.CharAt(position)
	}
}

func (ci *CharacterIndex) NextSignificantCharIs(c rune) bool {
	return ci.NextSignificantCharIsFromStartPosition(ci.position, c)
}

func (ci *CharacterIndex) NextSignificantCharIsFromStartPosition(startPosition int, c rune) bool {
	readPosition := startPosition + 1
	for ci.IsOutOfBounds(readPosition) && ci.CharAt(readPosition) == SPACE {
		readPosition++
	}
	return !ci.IsOutOfBounds(readPosition) && ci.CharAt(readPosition) == c
}

func (ci *CharacterIndex) NextSignificantChar() rune {
	return ci.NextSignificantCharFromStartPosition(ci.position)
}
func (ci *CharacterIndex) NextSignificantCharFromStartPosition(startPosition int) rune {
	readPosition := startPosition + 1
	for !ci.IsOutOfBounds(readPosition) && ci.CharAt(startPosition) == SPACE {
		readPosition++
	}
	if !ci.IsOutOfBounds(readPosition) {
		return ci.CharAt(readPosition)
	} else {
		return ' '
	}
}

func (ci *CharacterIndex) ReadSignificantChar(c rune) {
	if ci.SkipBlanks().CurrentChar() != c {
		// TODO:throw new InvalidPathException(String.format("Expected character: %c", c));
	}
	ci.IncrementPosition(1)
}

func (ci *CharacterIndex) HasSignificantSubSequence(s string) bool {
	ci.SkipBlanks()

	if !ci.InBoundsByPosition(ci.position + len(s) - 1) {
		return false
	}
	if ci.SubSequence(ci.position, ci.position+len(s)) != s {
		return false
	}
	ci.IncrementPosition(len(s))
	return true
}

func (ci *CharacterIndex) IndexOfPreviousSignificantCharFromStartPosition(startPosition int) int {
	readPosition := startPosition - 1

	for !ci.IsOutOfBounds(readPosition) && ci.CharAt(readPosition) == SPACE {
		readPosition--
	}

	if !ci.IsOutOfBounds(readPosition) {
		return readPosition
	} else {
		return -1
	}
}

func (ci *CharacterIndex) IndexOfPreviousSignificantChar() int {
	return ci.IndexOfPreviousSignificantCharFromStartPosition(ci.position)
}

func (ci *CharacterIndex) PreviousSignificantChar() rune {
	return ci.PreviousSignificantCharFromStartPosition(ci.position)
}

func (ci *CharacterIndex) PreviousSignificantCharFromStartPosition(startPosition int) rune {
	previousSignificantCharIndex := ci.IndexOfPreviousSignificantCharFromStartPosition(startPosition)

	if previousSignificantCharIndex == -1 {
		return ' '
	} else {
		return ci.CharAt(previousSignificantCharIndex)
	}
}

func (ci *CharacterIndex) CurrentIsTail() bool {
	return ci.position >= ci.endPosition
}

func (ci *CharacterIndex) HasMoreCharacters() bool {
	return ci.InBoundsByPosition(ci.position + 1)
}

func (ci *CharacterIndex) InBounds() bool {
	return ci.InBoundsByPosition(ci.position)
}

func (ci *CharacterIndex) IsOutOfBounds(idx int) bool {
	return !ci.InBoundsByPosition(idx)
}

func (ci *CharacterIndex) SubSequence(start int, end int) string {
	return ci.charSequence[start:end]
}

func (ci *CharacterIndex) CharSequence() string {
	return ci.charSequence
}

func (ci *CharacterIndex) String() string {
	return ci.charSequence
}

func (ci *CharacterIndex) IsNumberCharacter(readPosition int) bool {
	c := ci.CharAt(readPosition)
	return UtilsCharIsDigit(c) || c == MINUS || c == PERIOD || c == SCI_E || c == SCI_e
}

func (ci *CharacterIndex) SkipBlanks() *CharacterIndex {
	for ci.InBounds() && ci.position < ci.endPosition && ci.CurrentChar() == SPACE {
		ci.IncrementPosition(1)
	}
	return ci
}

func (ci *CharacterIndex) SkipBlanksAtEnd() *CharacterIndex {
	for ci.InBounds() && ci.position < ci.endPosition && ci.LastCharIs(SPACE) {
		ci.DecrementEndPosition(1)
	}
	return ci
}

func (ci *CharacterIndex) Trim() *CharacterIndex {
	ci.SkipBlanks()
	ci.SkipBlanksAtEnd()
	return ci
}

func CreateCharacterIndex(pathString string) *CharacterIndex {
	return &CharacterIndex{
		charSequence: pathString,
		position:     0,
		endPosition:  len(pathString) - 1,
	}
}
