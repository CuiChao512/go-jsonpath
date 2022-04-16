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
var UNDEFINED_NODE = &UndefinedNode{}

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

func NewPathNodeWithString(pathString string, existsCheck bool, shouldExist bool) (*PathNode, error) {
	compiledPath, err := path.Compile(pathString)
	if err != nil {
		return nil, err
	}
	return &PathNode{path: compiledPath, existsCheck: existsCheck, shouldExist: shouldExist}, nil
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
			return NewNumberNodeByString(resString), nil
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
	decimal2, err := decimal.NewFromString(str)
	if err == nil {
		return &NumberNode{
			number: &decimal2,
		}
	} else {
		return nil
	}

}

// StringNode -----------
type StringNode struct {
	*valueNodeDefault
	str            string
	useSingleQuote bool
}

func (n *StringNode) AsNumberNode() (*NumberNode, error) {
	number, err := decimal.NewFromString(n.str)
	if err != nil {
		return nil, nil
	} else {
		return NewNumberNode(&number), nil
	}
}

func (n *StringNode) GetString() string {
	return n.str
}

func (n *StringNode) Length() int {
	return len(n.str)
}

func (n *StringNode) IsEmpty() bool {
	return len(n.str) == 0
}

func (n *StringNode) Contains(str1 string) bool {
	return strings.Contains(n.str, str1)
}

func (n *StringNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	return reflect.String
}

func (n *StringNode) IsStringNode() bool {
	return true
}

func (n *StringNode) AsStringNode() (*StringNode, error) {
	return n, nil
}

func NewStringNode(str string, escape bool) *StringNode {
	return &StringNode{}
}

func (n *StringNode) String() string {
	quote := "\""
	if n.useSingleQuote {
		quote = "'"
	}
	//TODO: string escape
	return quote + n.str + quote
}

func (n *StringNode) Equals(o interface{}) bool {
	if n == o {
		return true
	}
	switch o.(type) {
	case *NumberNode:
		v, _ := o.(*NumberNode)
		that, _ := v.AsStringNode()
		if len(that.str) == 0 {
			return false
		} else {
			return n.str == that.str
		}
	case *StringNode:
		that, _ := o.(*StringNode)
		if len(that.str) == 0 {
			return false
		} else {
			return n.str == that.str
		}
	default:
		return false
	}
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

func NewValueListNode(list []interface{}) *ValueListNode {
	return nil
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

func (n *OffsetDateTimeNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
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
	json   interface{}
	parsed bool
}

func (n *JsonNode) TypeOf(ctx jsonpath.PredicateContext) reflect.Kind {
	if n.IsArray(ctx) {
		return reflect.Slice
	} else {
		parsedCtx, _ := n.Parse(ctx)
		switch parsedCtx.(type) {
		case decimal.Decimal:
			return reflect.Float64
		case string:
			return reflect.String
		case bool:
			return reflect.Bool
		default:
			return reflect.Invalid
		}
	}
}

func (*JsonNode) IsJsonNode() bool {
	return true
}

func (n *JsonNode) AsJsonNode() (*JsonNode, error) {
	return n, nil
}

func (n *JsonNode) IsParsed() bool {
	return n.parsed
}

func (n *JsonNode) GetJson() interface{} {
	return n.json
}

func (n *JsonNode) IsArray(ctx jsonpath.PredicateContext) bool {
	parsedObj, _ := n.Parse(ctx)
	return jsonpath.UtilsIsSlice(parsedObj)
}

func (n *JsonNode) IsMap(ctx jsonpath.PredicateContext) bool {
	parsedObj, _ := n.Parse(ctx)
	return jsonpath.UtilsIsMap(parsedObj)
}

func (n *JsonNode) Parse(ctx jsonpath.PredicateContext) (interface{}, error) {
	if n.parsed {
		return n.json, nil
	} else {
		//TODO:new JSONParser(JSONParser.MODE_PERMISSIVE).parse(json.toString());
		return nil, nil
	}
}

func (n *JsonNode) AsValueListNodeByPredicateContext(ctx jsonpath.PredicateContext) (ValueNode, error) {
	if !n.IsArray(ctx) {
		return UNDEFINED_NODE, nil
	} else {
		parsedObj, _ := n.Parse(ctx)
		list, _ := parsedObj.([]interface{})
		return NewValueListNode(list), nil
	}
}

func NewJsonNode(json string) *JsonNode {
	return &JsonNode{}
}
