package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"regexp"
	"strings"
)

// PatternNode -------patternNode------
type PatternNode struct {
	valueNodeDefault
	pattern         string
	compiledPattern *regexp.Regexp
}

func NewPatternNode(pattern string) *PatternNode {

	begin := strings.Index(pattern, "/")
	end := strings.LastIndex(pattern, "/")
	purePattern := pattern[begin:end]
	compiledPattern, _ := regexp.Compile(purePattern)
	return &PatternNode{pattern: purePattern, compiledPattern: compiledPattern}
}

func (pn *PatternNode) IsPatternNode() bool {
	return true
}

func (pn *PatternNode) AsPatternNode() (*PatternNode, *jsonpath.InvalidPathError) {
	return pn, nil
}

func (pn *PatternNode) String() string {
	if !strings.HasPrefix(pn.pattern, "/") {
		return "/" + pn.pattern + "/"
	} else {
		return pn.pattern
	}
}

// PathNode ------PathNode-----
type PathNode struct {
	valueNodeDefault
}

// NumberNode -----------
type NumberNode struct {
	valueNodeDefault
}

// StringNode -----------
type StringNode struct {
	valueNodeDefault
}

// BooleanNode -----------
type BooleanNode struct {
	valueNodeDefault
}

// PredicateNode -----------
type PredicateNode struct {
	valueNodeDefault
}

// ValueListNode -----------
type ValueListNode struct {
	valueNodeDefault
}

// NullNode -----------
type NullNode struct {
	valueNodeDefault
}

// UndefinedNode -----------
type UndefinedNode struct {
	valueNodeDefault
}

// ClassNode -----------
type ClassNode struct {
	valueNodeDefault
}

// OffsetDateTimeNode -----------
type OffsetDateTimeNode struct {
	valueNodeDefault
}
