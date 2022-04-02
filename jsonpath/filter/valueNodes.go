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

type NumberNode struct {
	valueNodeDefault
}

type StringNode struct {
	valueNodeDefault
}

type BooleanNode struct {
	valueNodeDefault
}

type PredicateNode struct {
	valueNodeDefault
}

type ValueListNode struct {
	valueNodeDefault
}

type NullNode struct {
	valueNodeDefault
}

type UndefinedNode struct {
	valueNodeDefault
}

type ClassNode struct {
	valueNodeDefault
}

type OffsetDateTimeNode struct {
	valueNodeDefault
}
