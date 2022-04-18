package test

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/common"
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
	{
		Key:      "int-key",
		Operator: eq,
		Value:    666,
		Expected: false,
	},
	{
		Key:      "int-key",
		Operator: eq,
		Value:    "1",
		Expected: true,
	},
	{
		Key:      "int-key",
		Operator: eq,
		Value:    "666",
		Expected: true,
	},
	{
		Expression: "[?(1 == '1')]",
		Expected:   true,
	},
	{
		Expression: "[?('1' == 1)]",
		Expected:   true,
	},
	{
		Expression: "[?(1 === '1')]",
		Expected:   false,
	},
	{
		Expression: "[?('1' === 1)]",
		Expected:   false,
	},
	{
		Expression: "[?(1 === 1)]",
		Expected:   true,
	},
	{
		Key:      "long-key",
		Operator: eq,
		Value:    3000000000,
		Expected: true,
	},
	{
		Key:      "long-key",
		Operator: eq,
		Value:    666,
		Expected: true,
	},
	{
		Key:      "float-key",
		Operator: eq,
		Value:    10.1,
		Expected: true,
	},
	{
		Key:      "float-key",
		Operator: eq,
		Value:    10.10,
		Expected: true,
	},
	{
		Key:      "float-key",
		Operator: eq,
		Value:    10.11,
		Expected: false,
	},
}

func TestIntEqEvals(t *testing.T) {
	for _, row := range testDataTable {
		if row.Type == "parse" {

		} else {
			criteria := jsonpath.WhereString(row.Key)
			switch row.Operator {
			case eq:
				criteria = criteria.Eq(row.Value)
			case lt:
				criteria = criteria.Lt(row.Value)
			case lte:
				criteria = criteria.Lte(row.Value)
			case gt:
				criteria = criteria.Gt(row.Value)
			case gte:
				criteria = criteria.Gte(row.Value)
			case ne:
				criteria = criteria.Ne(row.Value)
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

			result, err := criteria.Apply(createPredicateContext(FilterTestJson))
			if err != nil {
				t.Errorf("filter by key failed, err: %s", err)
			}
			if result != row.Expected {
				t.Errorf("%s = %t; expected %t", row.Key, result, row.Expected)
			}
		}
	}
}
