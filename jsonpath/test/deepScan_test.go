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
		//{
		//	JsonString: "{\"x\": [0,1,[0,1,2,3,null],null]}",
		//	PathString: "$..[2][3]",
		//	Expected:   []interface{}{float64(3)},
		//},
		//{
		//	JsonString: "{\"x\": [0,1,[0,1,2,3,null],null], \"y\": [0,1,2]}",
		//	PathString: "$..[2][3]",
		//	Expected:   []interface{}{float64(3)},
		//},
		//{
		//	JsonString: "{\"x\": [0,1,[0,1,2],null], \"y\": [0,1,2]}",
		//	PathString: "$..[2][3]",
		//	Expected:   []interface{}{},
		//},
		//
		////when_deep_scanning_null_subscription_is_ignored
		//{
		//	JsonString: "{\"x\": [null,null,[0,1,2,3,null],null]}",
		//	PathString: "$..[2][3]",
		//	Expected:   []interface{}{float64(3)},
		//},
		//{
		//	JsonString: "{\"x\": [null,null,[0,1,2,3,null],null], \"y\": [0,1,null]}",
		//	PathString: "$..[2][3]",
		//	Expected:   []interface{}{float64(3)},
		//},
		//
		////when_deep_scanning_array_index_oob_is_ignored
		//{
		//	JsonString: "{\"x\": [0,1,[0,1,2,3,10],null]}",
		//	PathString: "$..[4]",
		//	Expected:   []interface{}{float64(10)},
		//},
		//{
		//	JsonString: "{\"x\": [null,null,[0,1,2,3]], \"y\": [null,null,[0,1]]}",
		//	PathString: "$..[2][3]",
		//	Expected:   []interface{}{float64(3)},
		//},
		//
		////when_deep_scanning_illegal_property_access_is_ignored
		//{
		//	JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
		//	PathString: "$..foo",
		//	Function: func(data interface{}) interface{} {
		//		val := reflect.ValueOf(data)
		//		if val.Kind() == reflect.Slice {
		//			return val.Len()
		//		}
		//		return -1
		//	},
		//	Expected: 2,
		//},
		//{
		//	JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
		//	PathString: "$..foo.bar",
		//	Expected:   []interface{}{float64(4)},
		//},
		//{
		//	JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
		//	PathString: "$..[*].foo.bar",
		//	Expected:   []interface{}{float64(4)},
		//},
		//{
		//	JsonString: "{\"x\": {\"foo\": {\"baz\": 4}}, \"y\": {\"foo\": 1}}",
		//	PathString: "$..[*].foo.bar",
		//	Expected:   []interface{}{},
		//},
		//{
		//	JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
		//	PathString: "$..foo[?(@.bar)].bar",
		//	Expected:   []interface{}{float64(4)},
		//},
		//{
		//	JsonString: "{\"x\": {\"foo\": {\"bar\": 4}}, \"y\": {\"foo\": 1}}",
		//	PathString: "$..[*]foo[?(@.bar)].bar",
		//	Expected:   []interface{}{float64(4)},
		//},
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
			documentCtx, err = jsonpath.CreateParseContextImplByConfiguration(configuration).ParseString(data.JsonString)
		} else {
			documentCtx, err = jsonpath.JsonpathParseString(data.JsonString)

		}
		if err != nil {
			t.Errorf(err.Error())
		}
		result, err = documentCtx.Read(data.PathString)
		if err != nil {
			t.Errorf(err.Error())
		}
		if data.Function != nil {
			result = data.Function(result)
			fmt.Printf("result after func: %d \n", result)
		}
		if !reflect.DeepEqual(result, data.Expected) {
			t.Errorf("fail")
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
