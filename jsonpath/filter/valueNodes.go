package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/path"
	"regexp"
	"strings"
)

var NULL_NODE = NewNullNode()
var TRUE_NODE = NewBooleanNode(true)
var FALSE_NODE = NewBooleanNode(false)
var UNDEFINED_NODE = UndefinedNode{}

// PatternNode -------patternNode------
type PatternNode struct {
	*valueNodeDefault
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
	*valueNodeDefault
	path        *path.Path
	existsCheck bool
	shouldExist bool
}

func NewPathNodeWithString(pathString string, existsCheck bool, shouldExist bool) *PathNode {
	compiledPath := path.Compile(pathString)
	return &PathNode{path: &compiledPath, existsCheck: existsCheck, shouldExist: shouldExist}
}

func NewPathNode(path *path.Path, existsCheck bool, shouldExist bool) *PathNode {
	return &PathNode{path: path, existsCheck: existsCheck, shouldExist: shouldExist}
}

func (pn *PathNode) GetPath() *path.Path {
	return pn.path
}

func (pn *PathNode) IsExistsCheck() bool {
	return pn.existsCheck
}

func (pn *PathNode) ShouldExists() bool {
	return pn.shouldExist
}

func (pn *PathNode) IsPathNode() bool {
	return true
}

func (pn *PathNode) AsPathNode() (*PathNode, *jsonpath.InvalidPathError) {
	return pn, nil
}

func (pn *PathNode) AsExistsCheck(shouldExist bool) *PathNode {
	return &PathNode{
		path:        pn.path,
		existsCheck: true,
		shouldExist: shouldExist,
	}
}

func (pn *PathNode) String() string {
	if pn.existsCheck && !pn.shouldExist {
		return "!" + (*pn.path).String()
	} else {
		return (*pn.path).String()
	}
}

func (pn *PathNode) Evaluate(ctx jsonpath.PredicateContext) ValueNode {
	if pn.IsExistsCheck() {
		return FALSE_NODE
	} else {
		return UNDEFINED_NODE
	}
}

// NumberNode -----------
type NumberNode struct {
	*valueNodeDefault
}

// StringNode -----------
type StringNode struct {
	*valueNodeDefault
}

// BooleanNode -----------
type BooleanNode struct {
	*valueNodeDefault
	value bool
}

func NewBooleanNode(value bool) *BooleanNode {
	return &BooleanNode{
		value: value,
	}
}

// PredicateNode -----------
type PredicateNode struct {
	*valueNodeDefault
}

// ValueListNode -----------
type ValueListNode struct {
	*valueNodeDefault
}

// NullNode -----------
type NullNode struct {
	*valueNodeDefault
}

func NewNullNode() *NullNode {
	return &NullNode{}
}

// UndefinedNode -----------
type UndefinedNode struct {
	*valueNodeDefault
}

func NewUndefinedNode() *UndefinedNode {
	return &UndefinedNode{}
}

// ClassNode -----------
type ClassNode struct {
	*valueNodeDefault
}

// OffsetDateTimeNode -----------
type OffsetDateTimeNode struct {
	*valueNodeDefault
}
