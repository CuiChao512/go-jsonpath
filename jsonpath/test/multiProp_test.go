package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
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
