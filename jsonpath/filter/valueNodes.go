package filter

import "cuichao.com/go-jsonpath/jsonpath"

// PatternNode -------patternNode------
type PatternNode struct {
	valueNodeDefault
	pattern string
	flags   string
}

func (pn *PatternNode) IsPatternNode() bool {
	return true
}

func (pn *PatternNode) AsPatternNode() (*PatternNode, *jsonpath.InvalidPathError) {
	return pn, nil
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
