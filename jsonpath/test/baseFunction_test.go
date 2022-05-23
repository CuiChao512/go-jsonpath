package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

const (
	NUMBER_ERIES           = "{\"empty\": [], \"numbers\" : [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}"
	TEXT_SERIES            = "{\"urls\": [\"http://api.worldbank.org/countries/all/?format=json\", \"http://api.worldbank.org/countries/all/?format=json\"], \"text\" : [ \"a\", \"b\", \"c\", \"d\", \"e\", \"f\" ]}"
	TEXT_AND_NUMBER_SERIES = "{\"text\" : [ \"a\", \"b\", \"c\", \"d\", \"e\", \"f\" ], \"numbers\" : [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}"
)

func verifyFunction(conf *common.Configuration, pathExpr string, json string, expectedValue interface{}) (bool, string) {
	_, err := jsonpath.CreateParseContextImplByConfiguration(conf).ParseString(json)
	if err == nil {
		return false, err.Error()
	} else {
		return true, ""
	}
}
