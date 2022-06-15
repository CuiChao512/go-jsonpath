package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"testing"
)

func TestSumOfDouble(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.numbers.sum()", (10*(10+1))/float64(2))
	if !result {
		t.Errorf("not expected")
	}
}

func TestAverageOfEmptyListNegative(t *testing.T) {
	conf := common.DefaultConfiguration()

	_, err := verifyMathFunction(conf, "$.empty.avg()", nil)

	if err == nil {
		t.Errorf("not expected")
	} else {
		switch err.(type) {
		case *common.JsonPathError:
		default:
			t.Errorf("not expected")
		}
	}
}
