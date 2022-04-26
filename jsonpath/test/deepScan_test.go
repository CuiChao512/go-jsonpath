package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
	"testing"
)

type deepScanTestData struct {
	JsonString string
	PathString string
	Function   func(interface{}) interface{}
	Expected   interface{}
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
	}
	testMetaDataPathNotFoundError = [][]string{
		//{
		//	"{\"foo\": {\"bar\": null}}",
		//	"$.foo.bar.[5]",
		//},
		{
			"{\"foo\": {\"bar\": null}}",
			"$.foo.bar.[5, 10]",
		},
		//{
		//	"{\"foo\": {\"bar\": 4}}",
		//	"$.foo.bar.[5]",
		//},
		//{
		//	"{\"foo\": {\"bar\": 4}}",
		//	"$.foo.bar.[5, 10]",
		//},
		//{
		//	"{\"foo\": {\"bar\": 4}}",
		//	"$.foo.bar.[5]",
		//},
	}
)

func TestDeepScan(t *testing.T) {
	var documentCtx jsonpath.DocumentContext
	var err error
	var result interface{}

	for _, data := range deepScanTestMetaDates {
		documentCtx, err = jsonpath.JsonpathParseString(data.JsonString)
		if err != nil {
			t.Errorf(err.Error())
		}
		result, err = documentCtx.Read(data.PathString)
		if err != nil {
			t.Errorf(err.Error())
		}
		if data.Function != nil {
			result = data.Function(result)
		}
		if !reflect.DeepEqual(result, data.Expected) {
			t.Errorf("fail")
		}
	}

	//errors
	for _, data := range testMetaDataPathNotFoundError {
		documentCtx, err = jsonpath.CreateParseContextImplByConfiguration(common.DefaultConfiguration()).ParseString(data[0])
		if err != nil {
			t.Errorf(err.Error())
		}
		result, err = documentCtx.Read(data[1])

		if err == nil {
			t.Errorf("path not found error ecpected")
		} else {
			switch err.(type) {
			case *common.PathNotFoundError:
			default:
				t.Errorf("path not found error expected,actual : %s", err)
			}
		}
	}
}
