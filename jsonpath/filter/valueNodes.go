package filter

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/path"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
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

func (pn *PatternNode) AsPatternNode() (*PatternNode, error) {
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
	path        path.Path
	existsCheck bool
	shouldExist bool
}

func NewPathNodeWithString(pathString string, existsCheck bool, shouldExist bool) *PathNode {
	compiledPath := path.Compile(pathString)
	return &PathNode{path: compiledPath, existsCheck: existsCheck, shouldExist: shouldExist}
}

func NewPathNode(path path.Path, existsCheck bool, shouldExist bool) *PathNode {
	return &PathNode{path: path, existsCheck: existsCheck, shouldExist: shouldExist}
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

func (pn *PathNode) AsPathNode() (*PathNode, error) {
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
		return "!" + pn.path.String()
	} else {
		return pn.path.String()
	}
}

func (pn *PathNode) GetPath() path.Path {
	return pn.path
}

func (pn *PathNode) Evaluate(ctx jsonpath.PredicateContext) (ValueNode, error) {
	if pn.IsExistsCheck() {
		c := &jsonpath.Configuration{} //TODO
		result, err := pn.path.Evaluate(ctx.Item(), ctx.Root(), c)
		if err == nil {
			if result == jsonpath.JsonProviderUndefined {
				return FALSE_NODE, nil
			} else {
				return TRUE_NODE, nil
			}
		} else {
			return FALSE_NODE, nil
		}
	} else {
		var res interface{}
		switch ctx.(type) {
		case *jsonpath.PredicateContextImpl:
			ctxi, _ := ctx.(*jsonpath.PredicateContextImpl)
			res = ctxi.Evaluate(pn.path)
		default:
			var doc interface{}
			if pn.path.IsRootPath() {
				doc = ctx.Root()
			} else {
				doc = ctx.Item()
			}

			evaCtx, _ := pn.path.Evaluate(doc, ctx.Root(), ctx.Configuration())
			res = evaCtx.GetValue()
		}

		res = ctx.Configuration().JsonProvider().Unwrap(res)
		resString := ""
		if res == nil {
			return NULL_NODE, nil
		} else if ctx.Configuration().JsonProvider().IsArray(res) {
			return NewJsonNode(resString), nil
		} else if ctx.Configuration().JsonProvider().IsMap(res) {
			return NewJsonNode(resString), nil
		}
		switch res.(type) {
		case int:
			return NewNumberNode(resString), nil
		case float32:
		case float64:
		case string:
		case bool:
			resBool := false
			if resString == "true" {
				resBool = true
			}
			return NewBooleanNode(resBool), nil
		case *OffsetDateTimeNode:
		default:
			return nil, &jsonpath.JsonPathError{Message: fmt.Sprintf("Could not convert %t: %s to a ValueNode", res, resString)}
		}

		return UNDEFINED_NODE, nil
	}
}

// NumberNode -----------
type NumberNode struct {
	*valueNodeDefault
	number *decimal.Decimal
}

func (n *NumberNode) AsStringNode() (*StringNode, error) {
	return NewStringNode(n.number.String(), false), nil
}

func (n *NumberNode) GetNumber() *decimal.Decimal {
	return n.number
}

func (n *NumberNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.Float64
}

func (n *NumberNode) IsNumberNode() bool {
	return true
}

func (n *NumberNode) AsNumberNode() (*NumberNode, error) {
	return n, nil
}

func (n *NumberNode) String() string {
	return n.number.String()
}

func (n *NumberNode) Equals(o interface{}) bool {
	if n == o {
		return true
	}
	switch o.(type) {
	case *NumberNode:
		that, _ := o.(*NumberNode)
		if that.number == nil {
			return false
		} else {
			return n.number.Equals(*that.number)
		}
	case *StringNode:
		v, _ := o.(*StringNode)
		that, _ := v.AsNumberNode()
		if that.number == nil {
			return false
		} else {
			return n.number.Equals(*that.number)
		}
	default:
		return false
	}
}

func NewNumberNode(decimal2 *decimal.Decimal) *NumberNode {
	return &NumberNode{
		number: decimal2,
	}
}

func NewNumberNodeByString(str string) *NumberNode {
	return &NumberNode{}
}

// StringNode -----------
type StringNode struct {
	*valueNodeDefault
}

func NewStringNode(str string, escape bool) *StringNode {
	return &StringNode{}
}

// BooleanNode -----------
type BooleanNode struct {
	*valueNodeDefault
	value bool
}

func (*BooleanNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.Bool
}

func (*BooleanNode) IsBooleanNode() bool {
	return true
}

func (n *BooleanNode) AsBooleanNode() (*BooleanNode, error) {
	return n, nil
}

func (n *BooleanNode) GetBoolean() bool {
	return n.value
}

func (n *BooleanNode) String() string {
	if n.value {
		return "true"
	} else {
		return "false"
	}
}

func (n *BooleanNode) Equals(o interface{}) bool {
	if n == o {
		return true
	}
	switch o.(type) {
	case *BooleanNode:
		that, _ := o.(bool)
		return n.value == that
	default:
		return false
	}
}

func NewBooleanNode(value bool) *BooleanNode {
	return &BooleanNode{
		value: value,
	}
}

// PredicateNode -----------
type PredicateNode struct {
	*valueNodeDefault
	predicate jsonpath.Predicate
}

func (n *PredicateNode) GetPredicate() jsonpath.Predicate {
	return n.predicate
}

func (n *PredicateNode) AsPredicateNode() (*PredicateNode, error) {
	return n, nil
}

func (n *PredicateNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.Invalid
}

func (n *PredicateNode) IsPredicateNode() bool {
	return true
}

func (n *PredicateNode) Equals(o interface{}) bool {
	return false
}

func (n *PredicateNode) String() string {
	return n.predicate.String()
}

// ValueListNode -----------
type ValueListNode struct {
	*valueNodeDefault
	nodes []ValueNode
}

func (v *ValueListNode) Contains(node ValueNode) bool {
	return jsonpath.UtilsSliceContains(v.nodes, node)
}

func (v *ValueListNode) SubSetOf(right *ValueListNode) bool {
	for _, leftNode := range v.nodes {
		if !jsonpath.UtilsSliceContains(right, leftNode) {
			return false
		}
	}
	return true
}

func (v *ValueListNode) GetNodes() []ValueNode {
	return v.nodes
}

func (v *ValueListNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.Slice
}

func (v *ValueListNode) IsValueListNode() bool {
	return true
}

// NullNode -----------
type NullNode struct {
	*valueNodeDefault
}

func (n *NullNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.Invalid
}

func (n *NullNode) IsNullNode() bool {
	return true
}

func (n *NullNode) AsNullNode() (*NullNode, error) {
	return n, nil
}

func (n *NullNode) String() string {
	return "null"
}

func (n *NullNode) Equals(o interface{}) bool {
	if n == o {
		return true
	}
	switch o.(type) {
	case *NullNode:
		return true
	default:
		return false
	}
}

func NewNullNode() *NullNode {
	return &NullNode{}
}

// UndefinedNode -----------
type UndefinedNode struct {
	*valueNodeDefault
}

func (n *UndefinedNode) AsUndefinedNode() (*UndefinedNode, error) {
	return n, nil
}

func (n *UndefinedNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.Invalid
}

func (n *UndefinedNode) IsUndefinedNode() bool {
	return true
}

func (n *UndefinedNode) Equals(o interface{}) bool {
	return false
}
func NewUndefinedNode() *UndefinedNode {
	return &UndefinedNode{}
}

// ClassNode -----------
type ClassNode struct {
	*valueNodeDefault
}

// OffsetDateTime -----
type OffsetDateTime struct {
}

func (o *OffsetDateTime) String() string {
	return ""
}

// OffsetDateTimeNode -----------
type OffsetDateTimeNode struct {
	*valueNodeDefault
	dateTime *OffsetDateTime
}

func (n *OffsetDateTimeNode) AsStringNode() (*StringNode, error) {
	return NewStringNode(n.dateTime.String(), false), nil
}

func (n *OffsetDateTimeNode) GetDate() *OffsetDateTime {
	return n.dateTime
}

func (n *OffsetDateTimeNode) TypeOf(ctx *jsonpath.PredicateContext) reflect.Kind {
	return reflect.Interface
}

func (n *OffsetDateTimeNode) IsOffsetDateTimeNode() bool {
	return true
}

func (n *OffsetDateTimeNode) AsOffsetDateTimeNode() (*OffsetDateTimeNode, error) {
	return n, nil
}

func (n *OffsetDateTimeNode) String() string {
	return n.dateTime.String()
}

func (n *OffsetDateTimeNode) Equals(o interface{}) bool {
	if n == o {
		return true
	}
	switch o.(type) {
	case *OffsetDateTimeNode:
		v, _ := o.(ValueNode)
		that, _ := v.AsOffsetDateTimeNode()
		return OffsetDateTimeCompare(n.dateTime, that.dateTime) == 0
	case *StringNode:
		v, _ := o.(ValueNode)
		that, _ := v.AsOffsetDateTimeNode()
		return OffsetDateTimeCompare(n.dateTime, that.dateTime) == 0
	default:
		return false
	}
}

func OffsetDateTimeCompare(this *OffsetDateTime, that *OffsetDateTime) int {
	return 0
}

// JsonNode --------
type JsonNode struct {
	*valueNodeDefault
}

func NewJsonNode(json string) *JsonNode {
	return &JsonNode{}
}
