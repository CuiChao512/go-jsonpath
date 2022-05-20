package test

import (
	"fmt"
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
	"testing"
)

type deepScanTestData struct {
	Options    []common.Option
	JsonObject interface{}
	JsonString string
	PathString string
	Function   func(interface{}) interface{}
	Expected   interface{}
}

type pathNotFoundErrorTestData struct {
	Options    []common.Option
	JsonString string
	PathString string
}

var (
	deepScanTestMetaDates = []deepScanTestData{
		//when_deep_scanning_non_array_subscription_is_ignored
		{
			JsonString: "{\"x\": [0,1,[0,1,2,3,null],null]}",
			PathString: "$..[2][3]",
			Expected:   []interface{}{float64(3)},
		},
		{
			JsonString: "{\"x\": [0,1,[0,1,2,3,null],null], \"y\": [0,1,2]}",
			PathString: "$..[2][3]",
			Expected:   []interface{}{float64(3)},
		},
		{
			JsonString: "{\"x\": [0,1,[0,1,2],null], \"y\": [0,1,2]}",
			PathString: "$..[2][3]",
			Expected:   []interface{}{},
		},

		//when_deep_scanning_null_subscription_is_ignored
		{
			JsonString: "{\"x\": [null,null,[0,1,2,3,null],null]}",
			PathString: "$..[2][3]",
			Expected:   []interface{}{float64(3)},
		},
		{
			JsonString: "{\"x\": [null,null,[0,1,2,3,null],null], \"y\": [0,1,null]}",
			PathString: "$..[2][3]",
			Expected:   []interface{}{float64(3)},
		},

		//when_deep_scanning_array_index_oob_is_ignored
		{
			JsonString: "{\"x\": [0,1,[0,1,2,3,10],null]}",
			PathString: "$..[4]",
			Expected:   []interface{}{float64(10)},
		},
		{
			JsonString: "{\"x\": [null,null,[0,1,2,3]], \"y\": [null,null,[0,1]]}",
			PathString: "$..[2][3]",
			Expected:   []interface{}{float64(3)},
		},

		//when_deep_scanning_illegal_property_access_is_ignored
		{
			JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
			PathString: "$..foo",
			Function: func(data interface{}) interface{} {
				val := reflect.ValueOf(data)
				if val.Kind() == reflect.Slice {
					return val.Len()
				}
				return -1
			},
			Expected: 2,
		},
		{
			JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
			PathString: "$..foo.bar",
			Expected:   []interface{}{float64(4)},
		},
		{
			JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
			PathString: "$..[*].foo.bar",
			Expected:   []interface{}{float64(4)},
		},
		{
			JsonString: "{\"x\": {\"foo\": {\"baz\": 4}}, \"y\": {\"foo\": 1}}",
			PathString: "$..[*].foo.bar",
			Expected:   []interface{}{},
		},
		{
			JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
			PathString: "$..foo[?(@.bar)].bar",
			Expected:   []interface{}{float64(4)},
		},
		{
			JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
			PathString: "$..[*]foo[?(@.bar)].bar",
			Expected:   []interface{}{float64(4)},
		},
		//when_deep_scanning_require_properties_is_ignored_on_scan_target
		{
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			JsonString: "[{\"x\": {\"foo\": {\"x\": 4}, \"x\": null}, \"y\": {\"x\": 1}}, {\"x\": []}]",
			PathString: "$..x",
			Function:   sizeOf,
			Expected:   5,
		},
		{
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			JsonString: "{\"foo\": {\"bar\": 4}}",
			PathString: "$..foo.bar",
			Expected:   []interface{}{float64(4)},
		},
		//when_deep_scanning_leaf_multi_props_work
		{
			JsonString: "[{\"a\": \"a-val\", \"b\": \"b-val\", \"c\": \"c-val\"}, [1, 5], {\"a\": \"a-val\"}]",
			PathString: "$..['a', 'c']",
			Expected: []interface{}{map[string]interface{}{
				"a": "a-val",
				"c": "c-val",
			}},
		},
		{
			Options:    []common.Option{common.OPTION_DEFAULT_PATH_LEAF_TO_NULL},
			JsonString: "[{\"a\": \"a-val\", \"b\": \"b-val\", \"c\": \"c-val\"}, [1, 5], {\"a\": \"a-val\"}]",
			PathString: "$..['a', 'c']",
			Expected: []interface{}{
				map[string]interface{}{
					"a": "a-val",
					"c": "c-val",
				},
				map[string]interface{}{
					"a": "a-val",
					"c": nil,
				},
			},
		},
		//require_single_property_ok
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"a": "a0",
				},
				map[string]interface{}{
					"a": "a1",
				},
			},
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			PathString: "$..a",
			Expected:   []interface{}{"a0", "a1"},
		},
		//require_single_property
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"a": "a0",
				},
				map[string]interface{}{
					"b": "b2",
				},
			},
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			PathString: "$..a",
			Expected:   []interface{}{"a0"},
		},
		//require_multi_property_all_match
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"a": "aa",
					"b": "bb",
				},
				map[string]interface{}{
					"a": "aa",
					"b": "bb",
				},
			},
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			PathString: "$..['a', 'b']",
			Expected: []interface{}{
				map[string]interface{}{
					"a": "aa",
					"b": "bb",
				},
				map[string]interface{}{
					"a": "aa",
					"b": "bb",
				},
			},
		},
		//require_multi_property_some_match
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"a": "aa",
					"b": "bb",
				},
				map[string]interface{}{
					"a": "aa",
					"d": "dd",
				},
			},
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			PathString: "$..['a', 'b']",
			Expected: []interface{}{
				map[string]interface{}{
					"a": "aa",
					"b": "bb",
				},
			},
		},
		//scan_for_single_property
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"a": "aa",
				},
				map[string]interface{}{
					"b": "bb",
				},
				map[string]interface{}{
					"b": map[string]interface{}{
						"b": "bb",
					},
					"ab": map[string]interface{}{
						"a": map[string]interface{}{
							"a": "aa",
						},
						"b": map[string]interface{}{
							"b": "bb",
						},
					},
				},
			},
			PathString: "$..['a']",
			Expected: []interface{}{
				"aa",
				map[string]interface{}{
					"a": "aa",
				},
				"aa",
			},
		},
		//scan_for_property_path
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"a": "aa",
				},
				map[string]interface{}{
					"x": "xx",
				},
				map[string]interface{}{
					"a": map[string]interface{}{
						"x": "xx",
					},
				},
				map[string]interface{}{
					"z": map[string]interface{}{
						"a": map[string]interface{}{
							"x": "xx",
						},
					},
				},
			},
			PathString: "$..['a'].x",
			Expected: []interface{}{
				"xx",
				"xx",
			},
		},
		//scan_for_property_path_missing_required_property
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"a": "aa",
				},
				map[string]interface{}{
					"x": "xx",
				},
				map[string]interface{}{
					"a": map[string]interface{}{
						"x": "xx",
					},
				},
				map[string]interface{}{
					"z": map[string]interface{}{
						"a": map[string]interface{}{
							"x": "xx",
						},
					},
				},
			},
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			PathString: "$..['a'].x",
			Expected: []interface{}{
				"xx",
				"xx",
			},
		},
		//scans_can_be_filtered
		{
			JsonObject: []interface{}{
				map[string]interface{}{
					"mammal": true,
					"color": map[string]interface{}{
						"val": "brown",
					},
				},
				map[string]interface{}{
					"mammal": true,
					"color": map[string]interface{}{
						"val": "white",
					},
				},
				map[string]interface{}{
					"mammal": false,
				},
			},
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			PathString: "$..[?(@.mammal == true)].color",
			Expected: []interface{}{
				map[string]interface{}{
					"val": "brown",
				},
				map[string]interface{}{
					"val": "white",
				},
			},
		},
		//scan_with_a_function_filter
		{
			JsonString: TestJsonDocument,
			PathString: "$..*[?(@.length() > 5)]",
			Function:   sizeOf,
			Expected:   1,
		},
		//deepScanPathDefault
		{
			JsonString: `{"index": "index", "data": {"array": [{ "object1": { "name": "robert"} }]}}`,
			PathString: "$..array[0]",
			Expected: []interface{}{
				map[string]interface{}{
					"object1": map[string]interface{}{
						"name": "robert",
					},
				},
			},
		},
		//deepScanPathRequireProperties
		{
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			JsonString: `{"index": "index", "data": {"array": [{ "object1": { "name": "robert"} }]}}`,
			PathString: "$..array[0]",
			Expected: []interface{}{
				map[string]interface{}{
					"object1": map[string]interface{}{
						"name": "robert",
					},
				},
			},
		},
	}

	testMetaDataPathNotFoundError = []pathNotFoundErrorTestData{
		{
			JsonString: "{\"foo\": {\"bar\": null}}",
			PathString: "$.foo.bar.[5]",
		},
		{
			JsonString: "{\"foo\": {\"bar\": null}}",
			PathString: "$.foo.bar.[5, 10]",
		},
		{
			JsonString: "{\"foo\": {\"bar\": 4}}",
			PathString: "$.foo.bar.[5]",
		},
		{
			JsonString: "{\"foo\": {\"bar\": 4}}",
			PathString: "$.foo.bar.[5, 10]",
		},
		{
			JsonString: "{\"foo\": {\"bar\": 4}}",
			PathString: "$.foo.bar.[5]",
		},
		//when_deep_scanning_require_properties_is_ignored_on_scan_target
		{
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			JsonString: "{\"foo\": {\"baz\": 4}}",
			PathString: "$..foo.bar",
		},
		//when_deep_scanning_require_properties_is_ignored_on_scan_target_but_not_on_children
		{
			Options:    []common.Option{common.OPTION_REQUIRE_PROPERTIES},
			JsonString: "{\"foo\": {\"baz\": 4}}",
			PathString: "$..foo.bar",
		},
	}
)

func TestDeepScan(t *testing.T) {
	var documentCtx jsonpath.DocumentContext
	var err error
	var result interface{}

	for _, data := range deepScanTestMetaDates {
		if data.Options != nil && len(data.Options) > 0 {
			configuration := common.DefaultConfiguration()
			for _, op := range data.Options {
				configuration.AddOptions(op)
			}
			if data.JsonObject != nil {
				documentCtx, err = jsonpath.CreateParseContextImplByConfiguration(configuration).ParseAny(data.JsonObject)
			} else {
				documentCtx, err = jsonpath.CreateParseContextImplByConfiguration(configuration).ParseString(data.JsonString)
			}
		} else {
			if data.JsonObject != nil {
				documentCtx, err = jsonpath.JsonpathParseObject(data.JsonObject)
			} else {
				documentCtx, err = jsonpath.JsonpathParseString(data.JsonString)
			}
		}
		if err != nil {
			t.Errorf(err.Error())
		} else {
			result, err = documentCtx.Read(data.PathString)
			if err != nil {
				t.Errorf(err.Error())
			} else {
				if data.Function != nil {
					result = data.Function(result)
					fmt.Printf("result after func: %d \n", result)
				}
				if !reflect.DeepEqual(result, data.Expected) {
					t.Errorf("fail")
				}
			}
		}
	}
}

func Test_deepScan_error(t *testing.T) {
	var documentCtx jsonpath.DocumentContext
	var err error
	//errors
	for _, data := range testMetaDataPathNotFoundError {
		configuration := common.DefaultConfiguration()
		if data.Options != nil && len(data.Options) > 0 {
			for _, op := range data.Options {
				configuration.AddOptions(op)
			}
		}

		documentCtx, err = jsonpath.CreateParseContextImplByConfiguration(configuration).ParseString(data.JsonString)
		if err != nil {
			t.Errorf(err.Error())
		}
		_, err = documentCtx.Read(data.PathString)

		if err == nil {
			t.Errorf("path not found error expected")
		} else {
			switch err.(type) {
			case *common.PathNotFoundError:
			default:
				t.Errorf("path not found error expected,actual : %s", err)
			}
		}
	}
}
