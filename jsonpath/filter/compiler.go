package filter

const (
	DOC_CONTEXT          = "$"
	EVAL_CONTEXT         = "@"
	OPEN_SQUARE_BRACKET  = "["
	CLOSE_SQUARE_BRACKET = "]"
	OPEN_PARENTHESIS     = "("
	CLOSE_PARENTHESIS    = ")"
	OPEN_OBJECT          = "{"
	CLOSE_OBJECT         = "}"
	OPEN_ARRAY           = "["
	CLOSE_ARRAY          = "]"

	SINGLE_QUOTE = "'"
	DOUBLE_QUOTE = "\""

	SPACE  = " "
	PERIOD = "."

	AND = "&"
	OR  = "|"

	MINUS       = "-"
	LT          = "<"
	GT          = ">"
	EQ          = "="
	TILDE       = "~"
	TRUE        = "t"
	FALSE       = "f"
	NULL        = "n"
	NOT         = "!"
	PATTERN     = "/"
	IGNORE_CASE = "i"
)

type Compiler struct {
}
