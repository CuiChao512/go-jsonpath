package test

import (
	"encoding/json"
	"fmt"
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/filter"
	"regexp"
	"testing"
)

var FilterTestJson, _ = common.DefaultConfiguration().JsonProvider().Parse("{" +
	"  \"int-key\" : 1, " +
	"  \"long-key\" : 3000000000, " +
	"  \"float-key\" : 10.1, " +
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
	eq       relationOperator = 1
	ne       relationOperator = 2
	lt       relationOperator = 3
	lte      relationOperator = 4
	gt       relationOperator = 5
	gte      relationOperator = 6
	regex    relationOperator = 7
	in       relationOperator = 8
	nin      relationOperator = 9
	all      relationOperator = 10
	size     relationOperator = 11
	subSetOf relationOperator = 12
	anyOf    relationOperator = 13
	noneOf   relationOperator = 14
	exists   relationOperator = 15
	typeOf   relationOperator = 16
	notEmpty relationOperator = 17
	empty    relationOperator = 18
)

type testDataRow struct {
	Type       string
	Expression string
	Key        string
	Operator   relationOperator
	Value      interface{}
	Expected   bool
}

func (t testDataRow) String() string {
	str, _ := json.Marshal(t)
	return string(str)
}

//equals
var (
	testMetaDataEqual = []testDataRow{
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
			Expected: false,
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
			Expected: false,
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
		{
			Key:      "string-key",
			Operator: eq,
			Value:    "string",
			Expected: true,
		},
		{
			Key:      "string-key",
			Operator: eq,
			Value:    "666",
			Expected: false,
		},
		{
			Key:      "boolean-key",
			Operator: eq,
			Value:    true,
			Expected: true,
		},
		{
			Key:      "boolean-key",
			Operator: eq,
			Value:    false,
			Expected: false,
		},
		{
			Key:      "null-key",
			Operator: eq,
			Value:    nil,
			Expected: true,
		},
		{
			Key:      "null-key",
			Operator: eq,
			Value:    "666",
			Expected: false,
		},
		{
			Key:      "string-key",
			Operator: eq,
			Value:    nil,
			Expected: false,
		},
		{
			Key:      "arr-empty",
			Operator: eq,
			Value:    "[]",
			Expected: true,
		},
		{
			Key:      "int-arr",
			Operator: eq,
			Value:    "[0,1,2,3,4]",
			Expected: true,
		},
		{
			Key:      "int-arr",
			Operator: eq,
			Value:    "[0,1,2,3]",
			Expected: false,
		},
		{
			Key:      "int-arr",
			Operator: eq,
			Value:    "[0,1,2,3,4,5]",
			Expected: false,
		},
	}

	testMetaDataNotEquals = []testDataRow{
		{
			Key:      "int-key",
			Operator: ne,
			Value:    1,
			Expected: false,
		},
		{
			Key:      "int-key",
			Operator: ne,
			Value:    666,
			Expected: true,
		},
		{
			Key:      "long-key",
			Operator: ne,
			Value:    3000000000,
			Expected: false,
		},
		{
			Key:      "long-key",
			Operator: ne,
			Value:    666,
			Expected: true,
		},
		{
			Key:      "float-key",
			Operator: ne,
			Value:    10.1,
			Expected: false,
		},
		{
			Key:      "float-key",
			Operator: ne,
			Value:    10.10,
			Expected: false,
		},
		{
			Key:      "float-key",
			Operator: ne,
			Value:    10.11,
			Expected: true,
		},
		{
			Key:      "string-key",
			Operator: ne,
			Value:    "string",
			Expected: false,
		},
		{
			Key:      "string-key",
			Operator: ne,
			Value:    "666",
			Expected: true,
		},
		{
			Key:      "boolean-key",
			Operator: ne,
			Value:    true,
			Expected: false,
		},
		{
			Key:      "boolean-key",
			Operator: ne,
			Value:    false,
			Expected: true,
		},
		{
			Key:      "null-key",
			Operator: ne,
			Value:    nil,
			Expected: false,
		},
		{
			Key:      "null-key",
			Operator: ne,
			Value:    "666",
			Expected: true,
		},
		{
			Key:      "string-key",
			Operator: ne,
			Value:    nil,
			Expected: true,
		},
	}

	testMetaDataLt = []testDataRow{
		{
			Key:      "int-key",
			Operator: lt,
			Value:    10,
			Expected: true,
		},
		{
			Key:      "int-key",
			Operator: lt,
			Value:    0,
			Expected: false,
		},
		{
			Key:      "long-key",
			Operator: lt,
			Value:    4000000000,
			Expected: true,
		},
		{
			Key:      "long-key",
			Operator: lt,
			Value:    666,
			Expected: false,
		},
		{
			Key:      "float-key",
			Operator: lt,
			Value:    100.0,
			Expected: true,
		},
		{
			Key:      "float-key",
			Operator: lt,
			Value:    1.1,
			Expected: false,
		},
		{
			Key:      "char-key",
			Operator: lt,
			Value:    "x",
			Expected: true,
		},
		{
			Key:      "char-key",
			Operator: lt,
			Value:    "a",
			Expected: false,
		},
	}

	testMetaDataLte = []testDataRow{
		{
			Key:      "int-key",
			Operator: lte,
			Value:    10,
			Expected: true,
		},
		{
			Key:      "int-key",
			Operator: lte,
			Value:    1,
			Expected: true,
		},
		{
			Key:      "int-key",
			Operator: lte,
			Value:    0,
			Expected: false,
		},
		{
			Key:      "long-key",
			Operator: lte,
			Value:    4000000000,
			Expected: true,
		},
		{
			Key:      "long-key",
			Operator: lte,
			Value:    3000000000,
			Expected: true,
		},
		{
			Key:      "long-key",
			Operator: lte,
			Value:    666,
			Expected: false,
		},
		{
			Key:      "float-key",
			Operator: lte,
			Value:    100.0,
			Expected: true,
		},
		{
			Key:      "float-key",
			Operator: lte,
			Value:    10.1,
			Expected: true,
		},
		{
			Key:      "float-key",
			Operator: lte,
			Value:    1.1,
			Expected: false,
		},
	}

	testMetaDataGt = []testDataRow{
		{
			Key:      "int-key",
			Operator: gt,
			Value:    10,
			Expected: false,
		},
		{
			Key:      "int-key",
			Operator: gt,
			Value:    0,
			Expected: true,
		},
		{
			Key:      "long-key",
			Operator: gt,
			Value:    4000000000,
			Expected: false,
		},
		{
			Key:      "long-key",
			Operator: gt,
			Value:    666,
			Expected: true,
		},
		{
			Key:      "float-key",
			Operator: gt,
			Value:    100.0,
			Expected: false,
		},
		{
			Key:      "float-key",
			Operator: gt,
			Value:    1.1,
			Expected: true,
		},
		{
			Key:      "char-key",
			Operator: gt,
			Value:    "x",
			Expected: false,
		},
		{
			Key:      "char-key",
			Operator: gt,
			Value:    "a",
			Expected: true,
		},
	}

	testMetaDataGte = []testDataRow{
		{
			Key:      "int-key",
			Operator: gte,
			Value:    10,
			Expected: false,
		},
		{
			Key:      "int-key",
			Operator: gte,
			Value:    1,
			Expected: true,
		},
		{
			Key:      "int-key",
			Operator: gte,
			Value:    0,
			Expected: true,
		},
		{
			Key:      "long-key",
			Operator: gte,
			Value:    4000000000,
			Expected: false,
		},
		{
			Key:      "long-key",
			Operator: gte,
			Value:    3000000000,
			Expected: true,
		},
		{
			Key:      "long-key",
			Operator: gte,
			Value:    666,
			Expected: true,
		},
		{
			Key:      "float-key",
			Operator: gte,
			Value:    100.0,
			Expected: false,
		},
		{
			Key:      "float-key",
			Operator: gte,
			Value:    10.1,
			Expected: true,
		},
		{
			Key:      "float-key",
			Operator: gte,
			Value:    1.1,
			Expected: true,
		},
	}

	testMetaDataRegex = []testDataRow{
		{
			Key:        "string-key",
			Expression: "^string$",
			Operator:   regex,
			Expected:   true,
		},
		{
			Key:        "string-key",
			Expression: "^tring$",
			Operator:   regex,
			Expected:   false,
		},
		{
			Key:        "null-key",
			Expression: "^string$",
			Operator:   regex,
			Expected:   false,
		},
		{
			Key:        "int-key",
			Expression: "^string$",
			Operator:   regex,
			Expected:   false,
		},
	}

	testMetaDataStringIn = []testDataRow{
		{
			Key:      "string-key",
			Value:    []string{"a", "", "string"},
			Operator: in,
			Expected: true,
		},
		{
			Key:      "string-key",
			Value:    []string{"a", ""},
			Operator: in,
			Expected: false,
		},
		{
			Key:      "null-key",
			Value:    []interface{}{"a", nil},
			Operator: in,
			Expected: true,
		},
		{
			Key:      "null-key",
			Value:    []interface{}{"a", "b"},
			Operator: in,
			Expected: false,
		},
		{
			Key:      "string-arr",
			Value:    []interface{}{"a"},
			Operator: in,
			Expected: false,
		},
	}
	testMetaDataStringNin = []testDataRow{
		{
			Key:      "string-key",
			Value:    []string{"a", "", "string"},
			Operator: nin,
			Expected: false,
		},
		{
			Key:      "string-key",
			Value:    []string{"a", ""},
			Operator: nin,
			Expected: true,
		},
		{
			Key:      "null-key",
			Value:    []interface{}{"a", nil},
			Operator: nin,
			Expected: false,
		},
		{
			Key:      "null-key",
			Value:    []interface{}{"a", "b"},
			Operator: nin,
			Expected: true,
		},
		{
			Key:      "string-arr",
			Value:    []interface{}{"a"},
			Operator: nin,
			Expected: true,
		},
	}

	testMetaDataAll = []testDataRow{
		{
			Key:      "int-arr",
			Value:    []interface{}{0, 1},
			Operator: all,
			Expected: true,
		},
		{
			Key:      "int-arr",
			Value:    []interface{}{0, 7},
			Operator: all,
			Expected: false,
		},
		{
			Key:      "string-arr",
			Value:    []interface{}{"a", "b"},
			Operator: all,
			Expected: true,
		},
		{
			Key:      "string-arr",
			Value:    []interface{}{"a", "x"},
			Operator: all,
			Expected: false,
		},
		{
			Key:      "string-key",
			Value:    []interface{}{"a", "x"},
			Operator: all,
			Expected: false,
		},
	}

	testMetaDataSize = []testDataRow{
		{
			Key:      "string-arr",
			Value:    5,
			Operator: size,
			Expected: true,
		},
		{
			Key:      "string-arr",
			Value:    7,
			Operator: size,
			Expected: false,
		},
		{
			Key:      "string-key",
			Value:    6,
			Operator: size,
			Expected: true,
		},
		{
			Key:      "string-key",
			Value:    7,
			Operator: size,
			Expected: false,
		},
		{
			Key:      "int-key",
			Value:    6,
			Operator: size,
			Expected: false,
		},
		{
			Key:      "null-key",
			Value:    6,
			Operator: size,
			Expected: false,
		},
	}

	testMetaDataSubSetOf = []testDataRow{
		{
			Key:      "string-arr",
			Value:    []interface{}{"a", "b", "c", "d", "e", "f", "g"},
			Operator: subSetOf,
			Expected: true,
		},
		{
			Key:      "string-arr",
			Value:    []interface{}{"e", "d", "b", "c", "a"},
			Operator: subSetOf,
			Expected: true,
		},
		{
			Key:      "string-arr",
			Value:    []interface{}{"a", "b", "c", "d"},
			Operator: subSetOf,
			Expected: false,
		},
	}
	testMetaData = [][]testDataRow{
		//testMetaDataEqual,
		//testMetaDataNotEquals,
		//testMetaDataLt,
		//testMetaDataLte,
		//testMetaDataGt,
		//testMetaDataGte,
		//testMetaDataRegex,
		//testMetaDataStringIn,
		//testMetaDataStringNin,
		//testMetaDataAll,
		//testMetaDataSize,
		testMetaDataSubSetOf,
	}
)

func TestFilterEvaluations(t *testing.T) {
	total, passed := 0, 0
	for _, data := range testMetaData {
		for _, row := range data {
			total++
			var predicate common.Predicate
			if row.Expression != "" && row.Operator == 0 {
				var err error
				predicate, err = filter.Compile(row.Expression)
				if err != nil {
					t.Errorf(err.Error())
				}
			} else {
				//fmt.Println(row.Key)
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
				case in:
					criteria, err = criteria.InSlice(row.Value)
				case nin:
					criteria, err = criteria.NinSlice(row.Value)
				case all:
					criteria, err = criteria.AllSlice(row.Value)
				case size:
					size, ok := row.Value.(int)
					if !ok {
						t.Errorf("size should be a int")
					}
					criteria, err = criteria.Size(size)
				case subSetOf:
					criteria, err = criteria.SubSetOfSlice(row.Value)
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
				predicate = jsonpath.CreateSingleFilter(criteria)
			}
			result, err := predicate.Apply(createPredicateContext(FilterTestJson))
			if err != nil {
				t.Errorf("filter by key failed, err: %s", err)
			}
			if result != row.Expected {
				t.Errorf("filter by key %s = %s actual %t expected %t", row.Key, row.Value, result, row.Expected)
			} else {
				passed++
			}
		}
	}
	println()
	fmt.Printf("=========== Total:%d  pass:%d  fail:%d ===========", total, passed, total-passed)
	println()
}
