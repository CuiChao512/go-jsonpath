package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
	"testing"
)

func Test_multi_prop_can_be_read_from_root(t *testing.T) {
	model := map[string]interface{}{"a": "a-val", "b": "b-val", "c": "c-val"}
	if documentContext, err := getParseContextUsingDefaultConf().ParseAny(model); err != nil {
		t.Errorf(err.Error())
	} else {
		if result, err1 := documentContext.Read("$['a', 'b']"); err1 != nil {
			t.Errorf(err1.Error())
		} else {
			mapResult, ok := result.(map[string]interface{})
			if !ok {
				t.Errorf("result should be a map")
			} else {
				if mapResult["a"] != "a-val" || mapResult["b"] != "b-val" {
					t.Errorf("result should be a map, and contanis a and b")
				}
			}
		}
		if result, err1 := documentContext.Read("$['a', 'd']"); err1 != nil {
			t.Errorf(err1.Error())
		} else {
			mapResult, ok := result.(map[string]interface{})
			if !ok {
				t.Errorf("result should be a map")
			} else {
				if mapResult == nil || len(mapResult) != 1 || mapResult["a"] != "a-val" {
					t.Errorf("result should be a map, and contanis a")
				}
			}
		}
	}
}

//multi_props_can_be_defaulted_to_null ignore

func Test_multi_props_can_be_required(t *testing.T) {
	model := map[string]interface{}{"a": "a-val", "b": "b-val", "c": "c-val"}

	if documentContext, err := jsonpath.CreateParseContextImplByConfiguration(common.DefaultConfiguration().AddOptions(common.OPTION_REQUIRE_PROPERTIES)).ParseAny(model); err != nil {
		t.Errorf(err.Error())
	} else {
		if _, err1 := documentContext.Read("$['a', 'x']"); err1 != nil {
			switch err1.(type) {
			case *common.PathNotFoundError:
			default:
				t.Errorf("shuould throw path not found error")
			}
		} else {
			t.Errorf("shuould throw path not found error")
		}
	}
}

type multiPropsTestMetaData struct {
	ModelString string
	PathString  string
	Expect      interface{}
}

var multiPropsTestMetaDatas = []multiPropsTestMetaData{
	//multi_props_can_be_non_leafs
	{
		ModelString: "{\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}",
		PathString:  "$['a', 'c'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
	//nonexistent_non_leaf_multi_props_ignored
	{
		ModelString: "{\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}",
		PathString:  "$['d', 'a', 'c', 'm'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
	//multi_props_with_post_filter
	{
		ModelString: "{\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1, \"flag\": true}}",
		PathString:  "$['a', 'c'][?(@.flag)].v",
		Expect:      []interface{}{float64(1)},
	},
	//deep_scan_does_not_affect_non_leaf_multi_props
	{
		ModelString: "{\"v\": [[{}, 1, {\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1, \"flag\": true}}]]}",
		PathString:  "$..['a', 'c'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
	{
		ModelString: "{\"v\": [[{}, 1, {\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1, \"flag\": true}}]]}",
		PathString:  "$..['a', 'c'][?(@.flag)].v",
		Expect:      []interface{}{float64(1)},
	},
	//multi_props_can_be_in_the_middle
	{
		ModelString: "{\"x\": [null, {\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}]}",
		PathString:  "$.x[1]['a', 'c'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
	{
		ModelString: "{\"x\": [null, {\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}]}",
		PathString:  "$.x[*]['a', 'c'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
	{
		ModelString: "{\"x\": [null, {\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}]}",
		PathString:  "$[*][*]['a', 'c'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
	{
		ModelString: "{\"x\": [null, {\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}]}",
		PathString:  "$.x[1]['d', 'a', 'c', 'm'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
	{
		ModelString: "{\"x\": [null, {\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}]}",
		PathString:  "$.x[*]['d', 'a', 'c', 'm'].v",
		Expect:      []interface{}{float64(5), float64(1)},
	},
}

func Test_multi_props_cases(t *testing.T) {
	for _, testData := range multiPropsTestMetaDatas {
		if documentContext, err := jsonpath.JsonpathParseString(testData.ModelString); err != nil {
			t.Errorf(err.Error())
		} else {
			if result, err1 := documentContext.Read(testData.PathString); err1 != nil {
				t.Errorf(err1.Error())
			} else {
				if !reflect.DeepEqual(result, testData.Expect) {
					t.Errorf("result should be a map, and contanis a")
				}
			}
		}
	}
}

func Test_non_leaf_multi_props_can_be_required(t *testing.T) {
	json := "{\"a\": {\"v\": 5}, \"b\": {\"v\": 4}, \"c\": {\"v\": 1}}"
	if documentContext, err := jsonpath.CreateParseContextImplByConfiguration(common.DefaultConfiguration().AddOptions(common.OPTION_REQUIRE_PROPERTIES)).ParseString(json); err != nil {
		t.Errorf(err.Error())
	} else {
		if result, err1 := documentContext.Read("$['a', 'c'].v"); err1 != nil {
			t.Errorf(err1.Error())
		} else {
			if !reflect.DeepEqual(result, []interface{}{float64(5), float64(1)}) {
				t.Errorf("result should be a map, and contanis a")
			}
		}
		if _, err1 := documentContext.Read("$['d', 'a', 'c', 'm'].v"); err1 != nil {
			switch err1.(type) {
			case *common.PathNotFoundError:
			default:
				t.Errorf("shuould throw path not found error")
			}
		} else {
			t.Errorf("shuould throw path not found error")
		}
	}
}
