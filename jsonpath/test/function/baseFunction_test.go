package function

import (
	"errors"
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
)

const (
	NUMBER_ERIES           = "{\"empty\": [], \"numbers\" : [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}"
	TEXT_SERIES            = "{\"urls\": [\"http://api.worldbank.org/countries/all/?format=json\", \"http://api.worldbank.org/countries/all/?format=json\"], \"text\" : [ \"a\", \"b\", \"c\", \"d\", \"e\", \"f\" ]}"
	TEXT_AND_NUMBER_SERIES = "{\"text\" : [ \"a\", \"b\", \"c\", \"d\", \"e\", \"f\" ], \"numbers\" : [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}"
)

func verifyFunction(conf *common.Configuration, pathExpr string, json string, expectedValue interface{}) (bool, error) {
	parseContext, err := jsonpath.CreateParseContextImplByConfiguration(conf).ParseString(json)
	if err != nil {
		return false, err
	} else {
		if result, err1 := parseContext.Read(pathExpr); err1 != nil {
			return false, err1
		} else {
			if reflect.DeepEqual(conf.JsonProvider().Unwrap(result), expectedValue) {
				return true, nil
			} else {
				return false, errors.New(" not match")
			}
		}
	}
}

func verifyMathFunction(conf *common.Configuration, pathExpr string, expectedValue interface{}) (bool, error) {
	return verifyFunction(conf, pathExpr, NUMBER_ERIES, expectedValue)
}

func verifyTextFunction(conf *common.Configuration, pathExpr string, expectedValue interface{}) (bool, error) {
	return verifyFunction(conf, pathExpr, TEXT_SERIES, expectedValue)
}

func verifyTextAndNumberFunction(conf *common.Configuration, pathExpr string, expectedValue interface{}) (bool, error) {
	return verifyFunction(conf, pathExpr, TEXT_AND_NUMBER_SERIES, expectedValue)
}
