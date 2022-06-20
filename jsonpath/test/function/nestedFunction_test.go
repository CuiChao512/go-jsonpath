package function

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"testing"
)

func TestParameterAverageFunctionCall(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.avg($.numbers.min(), $.numbers.max())", 5.5)
	if !result {
		t.Errorf("not expected")
	}
}

func TestArrayAverageFunctionCallWithParameters(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.numbers.sum($.numbers.min(), $.numbers.max())", 66.0)
	if !result {
		t.Errorf("not expected")
	}
}

func TestJsonInnerArgumentArray(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.sum(5, 3, $.numbers.max(), 2)", 20.0)
	if !result {
		t.Errorf("not expected")
	}
}

func TestSimpleLiteralArgument(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.sum(5)", 5.0)
	if !result {
		t.Errorf("not expected")
	}

	result, _ = verifyMathFunction(conf, "$.sum(50)", 50.0)
	if !result {
		t.Errorf("not expected")
	}
}

func TestStringConcat(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyTextFunction(conf, "$.text.concat()", "abcdef")
	if !result {
		t.Errorf("not expected")
	}
}

func TestStringAndNumberConcat(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyTextAndNumberFunction(conf, "$.concat($.text[0], $.numbers[0])", "a1")
	if !result {
		t.Errorf("not expected")
	}
}

func TestStringConcatWithJSONParameter(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyTextFunction(conf, "$.text.concat(\"-\", \"ghijk\")", "abcdef-ghijk")
	if !result {
		t.Errorf("not expected")
	}
}

func TestAppendNumber(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.numbers.append(11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 0).avg()", 10.0)
	if !result {
		t.Errorf("not expected")
	}
}
