package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

var NULL_NODE = CreateNullNode()
var TRUE_NODE = CreateBooleanNode(true)
var FALSE_NODE = CreateBooleanNode(false)
var UNDEFINED_NODE = &UndefinedNode{}

// PatternNode -------patternNode------
type PatternNode struct {
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
	pattern         string
	compiledPattern *regexp.Regexp
}

func CreatePatternNodeByString(pattern string) (*PatternNode, error) {

	begin := strings.Index(pattern, "/")
	end := strings.LastIndex(pattern, "/")
	purePattern := pattern[begin+1 : end]
	compiledPattern, err := regexp.Compile(purePattern)
	if err != nil {
		return nil, err
	}
	return &PatternNode{pattern: purePattern, compiledPattern: compiledPattern}, nil
}

func CreatePatternNodeByRegexp(pattern *regexp.Regexp) *PatternNode {
	return &PatternNode{pattern: pattern.String(), compiledPattern: pattern}
}

func (pn *PatternNode) GetCompiledPattern() *regexp.Regexp {
	return pn.compiledPattern
}

func (pn *PatternNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
	return reflect.Invalid
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

func (pn *PatternNode) Equals(o interface{}) bool {
	if pn == o {
		return true
	}
	switch o.(type) {
	case *PatternNode:
		that, _ := o.(*PatternNode)
		return pn.compiledPattern == that.compiledPattern
	default:
		return false
	}

}

// PathNode ------PathNode-----
type PathNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
	path        common.Path
	existsCheck bool
	shouldExist bool
}

func CreatePathNodeWithString(pathString string, existsCheck bool, shouldExist bool) (*PathNode, error) {
	compiledPath, err := PathCompile(pathString)
	if err != nil {
		return nil, err
	}
	return &PathNode{path: compiledPath, existsCheck: existsCheck, shouldExist: shouldExist}, nil
}

func CreatePathNode(path common.Path, existsCheck bool, shouldExist bool) *PathNode {
	return &PathNode{path: path, existsCheck: existsCheck, shouldExist: shouldExist}
}

func (pn *PathNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
	return reflect.Invalid
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

func (pn *PathNode) GetPath() common.Path {
	return pn.path
}

func (pn *PathNode) Evaluate(ctx common.PredicateContext) (ValueNode, error) {
	if pn.IsExistsCheck() {
		c := common.CreateConfiguration(ctx.Configuration().JsonProvider(), []common.Option{common.OPTION_ALWAYS_RETURN_LIST}, ctx.Configuration().MappingProvider())
		evaluationCtx, err := pn.path.Evaluate(ctx.Item(), ctx.Root(), c)
		if err == nil {
			if result, err := evaluationCtx.GetValueUnwrap(false); err == nil {
				if result == common.JsonProviderUndefined {
					return FALSE_NODE, nil
				} else {
					return TRUE_NODE, nil
				}
			}
		}
		return FALSE_NODE, nil
	} else {
		var res interface{}
		switch ctx.(type) {
		case *common.PredicateContextImpl:
			ctxi, _ := ctx.(*common.PredicateContextImpl)
			var err error
			res, err = ctxi.Evaluate(pn.path)
			if err != nil {
				return UNDEFINED_NODE, nil
			}
		default:
			var doc interface{}
			if pn.path.IsRootPath() {
				doc = ctx.Root()
			} else {
				doc = ctx.Item()
			}

			evaCtx, _ := pn.path.Evaluate(doc, ctx.Root(), ctx.Configuration())
			var err error
			res, err = evaCtx.GetValue()
			if err != nil {
				return nil, err
			}
		}

		res = ctx.Configuration().JsonProvider().Unwrap(res)
		resString := common.UtilsToString(res)

		switch res.(type) {
		case int:
			return CreateNumberNodeByString(resString)
		case float32:
			return CreateNumberNodeByString(resString)
		case float64:
			return CreateNumberNodeByString(resString)
		case string:
			return CreateStringNode(resString, false)
		case bool:
			resBool := false
			if resString == "true" {
				resBool = true
			}
			return CreateBooleanNode(resBool), nil
		case *OffsetDateTimeNode:
			return CreateOffsetDateTimeNode(resString), nil
		}

		if res == nil {
			return NULL_NODE, nil
		} else if ctx.Configuration().JsonProvider().IsArray(res) {
			return CreateJsonNodeByObject(ctx.Configuration().MappingProvider().MapSlice(res, ctx.Configuration())), nil
		} else if ctx.Configuration().JsonProvider().IsMap(res) {
			return CreateJsonNodeByObject(ctx.Configuration().MappingProvider().MapMap(res, ctx.Configuration())), nil
		} else {
			return nil, &common.JsonPathError{Message: fmt.Sprintf("Could not convert %t: %s to a ValueNode", res, resString)}
		}
	}
}

func (pn *PathNode) Equals(o interface{}) bool {
	return false
}

// NumberNode -----------
type NumberNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
	number *decimal.Decimal
}

func (n *NumberNode) IsStringNode() bool {
	return false
}
func (n *NumberNode) AsStringNode() (*StringNode, error) {
	return CreateStringNode(n.number.String(), false)
}

func (n *NumberNode) GetNumber() *decimal.Decimal {
	return n.number
}

func (n *NumberNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
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

func CreateNumberNode(decimal2 *decimal.Decimal) *NumberNode {
	return &NumberNode{
		number: decimal2,
	}
}

func CreateNumberNodeByString(str string) (*NumberNode, error) {
	decimal2, err := decimal.NewFromString(str)
	if err == nil {
		return &NumberNode{
			number: &decimal2,
		}, nil
	} else {
		return nil, err
	}

}

// StringNode -----------
type StringNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
	str            string
	useSingleQuote bool
}

func (n *StringNode) IsNumberNode() bool {
	return false
}

func (n *StringNode) AsNumberNode() (*NumberNode, error) {
	number, err := decimal.NewFromString(n.str)
	if err != nil {
		return nil, err
	} else {
		return CreateNumberNode(&number), nil
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

func (n *StringNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
	return reflect.String
}

func (n *StringNode) IsStringNode() bool {
	return true
}

func (n *StringNode) AsStringNode() (*StringNode, error) {
	return n, nil
}

func CreateStringNode(str string, escape bool) (*StringNode, error) {
	runes := []rune(str)
	useSingleQuote := true
	if escape && len(str) > 1 {
		open := runes[0]
		closeC := runes[len(runes)-1]
		if open == '\'' && closeC == '\'' {
			str = str[1 : len(str)-1]
		} else if open == '"' && closeC == '"' {
			str = str[1 : len(str)-1]
			useSingleQuote = false
		}
		var err error
		str, err = common.UtilsStringUnescape(str)
		if err != nil {
			return nil, err
		}
	}

	return &StringNode{str: str, useSingleQuote: useSingleQuote}, nil
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
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
	value bool
}

func (*BooleanNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
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
		that, _ := o.(*BooleanNode)
		return n.value == that.value
	default:
		return false
	}
}

func CreateBooleanNodeByString(str string) *BooleanNode {
	if str == "true" {
		return TRUE_NODE
	}
	return FALSE_NODE
}

func CreateBooleanNode(value bool) *BooleanNode {
	return &BooleanNode{
		value: value,
	}
}

// PredicateNode -----------
type PredicateNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
	predicate common.Predicate
}

func (n *PredicateNode) GetPredicate() common.Predicate {
	return n.predicate
}

func (n *PredicateNode) AsPredicateNode() (*PredicateNode, error) {
	return n, nil
}

func (n *PredicateNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
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

func CreatePredicateNode(p common.Predicate) *PredicateNode {
	return &PredicateNode{predicate: p}
}

// ValueListNode -----------
type ValueListNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
	nodes []ValueNode
}

func (v *ValueListNode) Contains(node ValueNode) bool {
	return valueNodeSliceContains(v.nodes, node)
}

func (v *ValueListNode) SubSetOf(right *ValueListNode) bool {
	for _, leftNode := range v.nodes {
		if !valueNodeSliceContains(right.nodes, leftNode) {
			return false
		}
	}
	return true
}

func (v *ValueListNode) AsValueListNode() (*ValueListNode, error) {
	return v, nil
}

func (v *ValueListNode) GetNodes() []ValueNode {
	return v.nodes
}

func (v *ValueListNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
	return reflect.Slice
}

func (v *ValueListNode) IsValueListNode() bool {
	return true
}

func (v *ValueListNode) String() string {
	return "[" + common.UtilsJoin(",", "", v.nodes) + "]"
}

func (v *ValueListNode) Equals(o interface{}) bool {
	if v == o {
		return true
	}
	switch o.(type) {
	case *ValueListNode:
		that, _ := o.(ValueListNode)
		return common.UtilsSliceEquals(v.nodes, that.nodes)
	default:
		return false
	}
}

func CreateValueListNode(list interface{}) (*ValueListNode, error) {
	l, err := common.ConvertToAnySlice(list)
	if err != nil {
		return nil, err
	}
	var nodes = make([]ValueNode, 0)
	for _, value := range l {
		if vn, err := CreateValueNode(value); err != nil {
			return nil, err
		} else {
			nodes = append(nodes, vn)
		}
	}
	return &ValueListNode{nodes: nodes}, nil
}

// NullNode -----------
type NullNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
}

func (n *NullNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
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

func CreateNullNode() *NullNode {
	return &NullNode{}
}

// UndefinedNode -----------
type UndefinedNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
}

func (n *UndefinedNode) AsUndefinedNode() (*UndefinedNode, error) {
	return n, nil
}

func (n *UndefinedNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
	return reflect.Invalid
}

func (n *UndefinedNode) IsUndefinedNode() bool {
	return true
}

func (n *UndefinedNode) Equals(o interface{}) bool {
	return false
}

func (n *UndefinedNode) String() string {
	return common.UtilsToString(n)
}
func NewUndefinedNode() *UndefinedNode {
	return &UndefinedNode{}
}

// ClassNode -----------
type ClassNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultOffsetDateTimeNode
	*defaultJsonNode
}

// OffsetDateTime -----
type OffsetDateTime struct {
}

func (o *OffsetDateTime) String() string {
	return ""
}

// OffsetDateTimeNode -----------
type OffsetDateTimeNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultJsonNode
	dateTime *OffsetDateTime
}

func (n *OffsetDateTimeNode) IsStringNode() bool {
	return false
}

func (n *OffsetDateTimeNode) AsStringNode() (*StringNode, error) {
	return CreateStringNode(n.dateTime.String(), false)
}

func (n *OffsetDateTimeNode) GetDate() *OffsetDateTime {
	return n.dateTime
}

func (n *OffsetDateTimeNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
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
	//TODO:
	return 0
}

func CreateOffsetDateTimeNode(str string) *OffsetDateTimeNode {
	//TODO
	return &OffsetDateTimeNode{dateTime: nil}
}

// JsonNode --------
type JsonNode struct {
	*defaultPatternNode
	*defaultPathNode
	*defaultNumberNode
	*defaultStringNode
	*defaultBooleanNode
	*defaultPredicateNode
	*defaultValueListNode
	*defaultNullNode
	*defaultUndefinedNode
	*defaultClassNode
	*defaultOffsetDateTimeNode
	json   interface{}
	parsed bool
}

func (n *JsonNode) TypeOf(ctx common.PredicateContext) reflect.Kind {
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

func (n *JsonNode) IsEmpty(ctx common.PredicateContext) (bool, error) {
	if n.IsArray(ctx) || n.IsMap(ctx) {
		parseResult, err := n.Parse(ctx)
		if err != nil {
			return false, err
		}
		switch reflect.ValueOf(parseResult).Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Map:
			return reflect.ValueOf(parseResult).Len() == 0, nil
		}
	} else {
		parseResult, err := n.Parse(ctx)
		if err != nil {
			return false, err
		}
		switch parseResult.(type) {
		case string:
			str, _ := parseResult.(string)
			return len(str) == 0, nil
		}
	}
	return true, nil
}

func (n *JsonNode) IsParsed() bool {
	return n.parsed
}

func (n *JsonNode) GetJson() interface{} {
	return n.json
}

func (n *JsonNode) IsArray(ctx common.PredicateContext) bool {
	parsedObj, _ := n.Parse(ctx)
	return common.UtilsIsSlice(parsedObj)
}

func (n *JsonNode) IsMap(ctx common.PredicateContext) bool {
	parsedObj, _ := n.Parse(ctx)
	return common.UtilsIsMap(parsedObj)
}

func (n *JsonNode) Parse(ctx common.PredicateContext) (interface{}, error) {
	if n.parsed {
		return n.json, nil
	} else {
		//TODO:new JSONParser(JSONParser.MODE_PERMISSIVE).parse(json.toString());
		jsonString, ok := n.json.(string)
		if !ok {
			return nil, errors.New("json should be a string")
		}
		var result interface{}
		if err := json.Unmarshal([]byte(jsonString), &result); err == nil {
			return result, nil
		} else {
			return nil, err
		}
	}
}

func (n *JsonNode) EqualsByPredicateContext(jsonNode *JsonNode, ctx common.PredicateContext) (bool, error) {
	if n == jsonNode {
		return true, nil
	}

	if n.json != nil {
		parseResult, err := jsonNode.Parse(ctx)
		if err != nil {
			return false, err
		}
		return reflect.DeepEqual(n.json, parseResult), nil
	} else {
		return jsonNode.json != nil, nil
	}
}

func (n *JsonNode) AsValueListNodeByPredicateContext(ctx common.PredicateContext) (ValueNode, error) {
	if !n.IsArray(ctx) {
		return UNDEFINED_NODE, nil
	} else {
		parsedObj, _ := n.Parse(ctx)
		list, _ := parsedObj.([]interface{})
		return CreateValueListNode(list)
	}
}

func (n *JsonNode) String() string {
	return common.UtilsToString(n.json)
}

func (n *JsonNode) Equals(o interface{}) bool {
	if n == o {
		return true
	}
	switch o.(type) {
	case *JsonNode:
		v, _ := o.(JsonNode)
		return n.json == v.json
	default:
		return false
	}
}

func (n *JsonNode) Length(ctx common.PredicateContext) int {
	if n.IsArray(ctx) {
		if parsed, err := n.Parse(ctx); err == nil {
			val := reflect.ValueOf(parsed)
			if val.Kind() == reflect.Slice {
				return val.Len()
			}
		}
	}
	return -1
}

func CreateJsonNodeByString(json string) *JsonNode {
	return &JsonNode{json: json}
}

func CreateJsonNodeByObject(json interface{}) *JsonNode {
	return &JsonNode{json: json, parsed: true}
}

func isPath(o interface{}) bool {
	if o == nil || reflect.TypeOf(o).Kind() != reflect.String {
		return false
	}
	str := strings.TrimSpace(common.UtilsToString(o))
	if len(str) <= 0 {
		return false
	}
	c0 := []rune(str)[0]
	if c0 == '@' || c0 == '$' {
		_, err := PathCompile(str)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func isJson(o interface{}) bool {
	if o == nil || reflect.TypeOf(o).Kind() != reflect.String {
		return false
	}
	str := strings.TrimSpace(common.UtilsToString(o))
	if len(str) <= 1 {
		return false
	}
	runes := []rune(str)
	c0 := runes[0]
	c1 := runes[len(runes)-1]
	if (c0 == '[' && c1 == ']') || (c0 == '{' && c1 == '}') {
		var i interface{}
		err := json.Unmarshal([]byte(str), &i)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func valueNodeSliceContains(valueNodes []ValueNode, node ValueNode) bool {
	for _, vn := range valueNodes {
		if vn.Equals(node) {
			return true
		}
	}
	return false
}

func CreateValueNode(o interface{}) (ValueNode, error) {
	if o == nil {
		return NULL_NODE, nil
	}
	switch o.(type) {
	case ValueNode:
		vn, _ := o.(ValueNode)
		return vn, nil
	}

	if isPath(o) {
		return CreatePathNodeWithString(common.UtilsToString(o), false, false)
	} else if isJson(o) {
		return CreateJsonNodeByString(common.UtilsToString(o)), nil
	}

	switch o.(type) {
	case string:
		return CreateStringNode(common.UtilsToString(o), false)
	case rune:
		return CreateStringNode(common.UtilsToString(o), true)
	case int:
		return CreateNumberNodeByString(common.UtilsToString(o))
	case float64:
		return CreateNumberNodeByString(common.UtilsToString(o))
	case bool:
		return CreateBooleanNodeByString(common.UtilsToString(o)), nil
	case *regexp.Regexp:
		r, _ := o.(*regexp.Regexp)
		return CreatePatternNodeByRegexp(r), nil
	case OffsetDateTime:
		return CreateOffsetDateTimeNode(common.UtilsToString(o)), nil
	}
	return nil, &common.JsonPathError{Message: "Could not determine value type"}
}
