package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"reflect"
	"testing"
)

func TestDeepScan(t *testing.T) {
	documentCtx, err := jsonpath.JsonpathParseString("{\"x\": [0,1,[0,1,2,3,null],null]}")
	if err != nil {
		t.Errorf(err.Error())
	}
	result, err := documentCtx.Read("$..[2][3]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, []int{3}) {
		t.Errorf("fail")
	}
}
