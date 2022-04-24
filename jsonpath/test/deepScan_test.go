package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
	"testing"
)

func TestDeepScan(t *testing.T) {
	var documentCtx jsonpath.DocumentContext
	var err error
	var result interface{}
	//when_deep_scanning_non_array_subscription_is_ignored
	documentCtx, err = jsonpath.JsonpathParseString("{\"x\": [0,1,[0,1,2,3,null],null]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err = documentCtx.Read("$..[2][3]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, []interface{}{float64(3)}) {
		t.Errorf("fail")
	}

	documentCtx, err = jsonpath.JsonpathParseString("{\"x\": [0,1,[0,1,2,3,null],null], \"y\": [0,1,2]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err = documentCtx.Read("$..[2][3]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, []interface{}{float64(3)}) {
		t.Errorf("fail")
	}

	documentCtx, err = jsonpath.JsonpathParseString("{\"x\": [0,1,[0,1,2],null], \"y\": [0,1,2]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err = documentCtx.Read("$..[2][3]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if val := reflect.ValueOf(result); val.Kind() != reflect.Slice || val.Len() != 0 {
		t.Errorf("fail")
	}

	//when_deep_scanning_null_subscription_is_ignored
	documentCtx, err = jsonpath.JsonpathParseString("{\"x\": [null,null,[0,1,2,3,null],null]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err = documentCtx.Read("$..[2][3]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, []interface{}{float64(3)}) {
		t.Errorf("fail")
	}

	documentCtx, err = jsonpath.JsonpathParseString("{\"x\": [null,null,[0,1,2,3,null],null], \"y\": [0,1,2]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err = documentCtx.Read("$..[2][3]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, []interface{}{float64(3)}) {
		t.Errorf("fail")
	}

	//when_deep_scanning_array_index_oob_is_ignored
	documentCtx, err = jsonpath.JsonpathParseString("{\"x\": [0,1,[0,1,2,3,10],null]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err = documentCtx.Read("$..[4]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, []interface{}{float64(10)}) {
		t.Errorf("fail")
	}

	documentCtx, err = jsonpath.JsonpathParseString("{\"x\": [null,null,[0,1,2,3]], \"y\": [null,null,[0,1]]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err = documentCtx.Read("$..[2][3]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, []interface{}{float64(3)}) {
		t.Errorf("fail")
	}

	//errors
	jsonpath.CreateParseContextImplByConfiguration(common.DefaultConfiguration())
}
