package test

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/common"
	"fmt"
	"regexp"
	"testing"
)

var FilterTestJson, _ = common.DefaultConfiguration().JsonProvider().Parse("{" +
	"  \"int-key\" : 1, " +
	"  \"long-key\" : 3000000000, " +
	"  \"double-key\" : 10.1, " +
	"  \"boolean-key\" : true, " +
	"  \"null-key\" : null, " +
	"  \"string-key\" : \"string\", " +
	"  \"string-key-empty\" : \"\", " +
	"  \"char-key\" : \"c\", " +
	"  \"arr-empty\" : [], " +
	"  \"int-arr\" : [0,1,2,3,4], " +
	"  \"string-arr\" : [\"a\",\"b\",\"c\",\"d\",\"e\"] " +
	"}")

type relationOperator int

const (
	eq    relationOperator = 1
	ne    relationOperator = 2
	lt    relationOperator = 3
	lte   relationOperator = 4
	gt    relationOperator = 5
	gte   relationOperator = 6
	regex relationOperator = 7
)

type testDataRow struct {
	Type       string
	Expression string
	Key        string
	Operator   relationOperator
	Value      interface{}
	Expected   bool
}

var testDataTable = []testDataRow{
	{
		Key:      "int-key",
		Operator: eq,
		Value:    1,
		Expected: true,
	},
	//{
	//	Key:      "int-key",
	//	Operator: eq,
	//	Value:    666,
	//	Expected: false,
	//},
	//{
	//	Key:      "int-key",
	//	Operator: eq,
	//	Value:    "1",
	//	Expected: true,
	//},
	//{
	//	Key:      "int-key",
	//	Operator: eq,
	//	Value:    "666",
	//	Expected: false,
	//},
	//{
	//	Expression: "[?(1 == '1')]",
	//	Expected:   true,
	//},
	//{
	//	Expression: "[?('1' == 1)]",
	//	Expected:   true,
	//},
	//{
	//	Expression: "[?(1 === '1')]",
	//	Expected:   false,
	//},
	//{
	//	Expression: "[?('1' === 1)]",
	//	Expected:   false,
	//},
	//{
	//	Expression: "[?(1 === 1)]",
	//	Expected:   true,
	//},
	//{
	//	Key:      "long-key",
	//	Operator: eq,
	//	Value:    3000000000,
	//	Expected: true,
	//},
	//{
	//	Key:      "long-key",
	//	Operator: eq,
	//	Value:    666,
	//	Expected: false,
	//},
	//{
	//	Key:      "float-key",
	//	Operator: eq,
	//	Value:    10.1,
	//	Expected: true,
	//},
	//{
	//	Key:      "float-key",
	//	Operator: eq,
	//	Value:    10.10,
	//	Expected: true,
	//},
	//{
	//	Key:      "float-key",
	//	Operator: eq,
	//	Value:    10.11,
	//	Expected: false,
	//},
}

func TestFilterEvals(t *testing.T) {
	for _, row := range testDataTable {
		if row.Type == "parse" {

		} else {
			fmt.Println(row.Key)
			criteria, err := jsonpath.WhereString(row.Key)
			if err != nil {
				t.Errorf(err.Error())
			}
			switch row.Operator {
			case eq:
				criteria, err = criteria.Eq(row.Value)
			case lt:
				criteria, err = criteria.Lt(row.Value)
			case lte:
				criteria, err = criteria.Lte(row.Value)
			case gt:
				criteria, err = criteria.Gt(row.Value)
			case gte:
				criteria, err = criteria.Gte(row.Value)
			case ne:
				criteria, err = criteria.Ne(row.Value)
			case regex:
				compiledRegexp, err := regexp.Compile(row.Expression)
				if err != nil {
					t.Errorf("expreesion=%s not a regex string", row.Expression)
				}
				criteria, err = criteria.Regex(compiledRegexp)
				if err != nil {
					t.Errorf("expreesion=%s not a regex string", row.Expression)
				}
			}
			if err != nil {
				t.Errorf("%s", err)
			}
			filter := jsonpath.CreateSingleFilter(criteria)
			result, err := filter.Apply(createPredicateContext(FilterTestJson))
			if err != nil {
				t.Errorf("filter by key failed, err: %s", err)
			}
			if result != row.Expected {
				t.Errorf("filter by key %s = %s actual %t expected %t", row.Key, row.Value, result, row.Expected)
			}
		}
	}
}
