package function

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"testing"
)

func TestAverageOfDoubles(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.numbers.avg()", 5.5)
	if !result {
		t.Errorf("not expected")
	}
}

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

func TestMaxOfDoubles(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.numbers.max()", 10.0)
	if !result {
		t.Errorf("not expected")
	}
}

func TestMaxOfEmptyListNegative(t *testing.T) {
	conf := common.DefaultConfiguration()
	_, err := verifyMathFunction(conf, "$.empty.max()", nil)

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

func TestMinOfDoubles(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.numbers.min()", 1.0)
	if !result {
		t.Errorf("not expected")
	}
}

func TestMinOfEmptyListNegative(t *testing.T) {
	conf := common.DefaultConfiguration()
	_, err := verifyMathFunction(conf, "$.empty.min()", nil)

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

func TestStdDevOfDouble(t *testing.T) {
	conf := common.DefaultConfiguration()
	result, _ := verifyMathFunction(conf, "$.numbers.stddev()", 2.8722813232690143)
	if !result {
		t.Errorf("not expected")
	}
}

func TestStdDevOfEmptyListNegative(t *testing.T) {
	conf := common.DefaultConfiguration()
	_, err := verifyMathFunction(conf, "$.empty.stddev()", nil)

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
