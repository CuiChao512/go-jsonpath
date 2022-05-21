package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/filter"
	"reflect"
	"regexp"
	"testing"
)

var aFilterCanBeParsedTestDataSlice = []string{
	"[?(@.foo)]",
	"[?(@.foo == 1)]",
	"[?(@.foo == 1 || @['bar'])]",
	"[?(@.foo == 1 && @['bar'])]",
}

func Test_a_filter_can_be_parsed(t *testing.T) {
	for _, filterString := range aFilterCanBeParsedTestDataSlice {
		_, err := filter.Compile(filterString)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

var invalidPathErrorTestDataSlice = []string{
	"[?(@.foo == 1)",
	"[?(@.foo == 1) ||]",
	"[(@.foo == 1)]",
	"[?@.foo == 1)]",
}

func Test_an_invalid_filter_can_not_be_parsed(t *testing.T) {
	for _, filterString := range invalidPathErrorTestDataSlice {
		_, err := filter.Compile(filterString)
		if err == nil {
			t.Errorf("invalid path error expect")
		} else {
			switch err.(type) {
			case *common.InvalidPathError:
			default:
				t.Errorf("invalid path error expect")
			}
		}
	}
}

type canBeSerializedTestData struct {
	WhereString string
	Operator    relationOperator
	Value       interface{}
	ParseString string
}

var canBeSerializedTestDataSlice = []canBeSerializedTestData{
	//a_gte_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    gte,
		Value:       1,
		ParseString: "[?(@['a'] >= 1)]",
	},
	//a_lte_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    lte,
		Value:       1,
		ParseString: "[?(@['a'] <= 1)]",
	},
	//a_eq_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    eq,
		Value:       1,
		ParseString: "[?(@['a'] == 1)]",
	},
	//a_ne_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    ne,
		Value:       1,
		ParseString: "[?(@['a'] != 1)]",
	},
	//a_lt_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    lt,
		Value:       1,
		ParseString: "[?(@['a'] < 1)]",
	},
	//a_gt_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    gt,
		Value:       1,
		ParseString: "[?(@['a'] > 1)]",
	},
	//a_nin_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    nin,
		Value:       1,
		ParseString: "[?(@['a'] NIN [1])]",
	},
	//a_in_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    in,
		Value:       "a",
		ParseString: "[?(@['a'] IN ['a'])]",
	},
	//a_contains_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    contains,
		Value:       "a",
		ParseString: "[?(@['a'] CONTAINS 'a')]",
	},
	//a_all_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    all,
		Value:       []interface{}{"a", "b"},
		ParseString: "[?(@['a'] ALL ['a','b'])]",
	},
	//a_size_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    size,
		Value:       5,
		ParseString: "[?(@['a'] SIZE 5)]",
	},
	//a_subsetof_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    anyOf,
		Value:       []interface{}{},
		ParseString: "[?(@['a'] ANYOF [])]",
	},
	//a_noneof_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    noneOf,
		Value:       []interface{}{},
		ParseString: "[?(@['a'] NONEOF [])]",
	},
	//a_exists_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    exists,
		Value:       true,
		ParseString: "[?(@['a'])]",
	},
	//a_not_exists_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    exists,
		Value:       false,
		ParseString: "[?(!@['a'])]",
	},
	//TODO: a_type_filter_can_be_serialized
	//a_matches_filter_can_be_serialized --
	//a_not_empty_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    empty,
		Value:       false,
		ParseString: "[?(@['a'] EMPTY false)]",
	},
	//and_filter_can_be_serialized --
	//in_string_filter_can_be_serialized
	{
		WhereString: "a",
		Operator:    in,
		Value:       []interface{}{"1", "2"},
		ParseString: "[?(@['a'] IN ['1','2'])]",
	},
	//a_deep_path_filter_can_be_serialized
	{
		WhereString: "a.b.c",
		Operator:    in,
		Value:       []interface{}{"1", "2"},
		ParseString: "[?(@['a']['b']['c'] IN ['1','2'])]",
	},
	//a_regex_filter_can_be_serialized --
	//a_doc_ref_filter_can_be_serialized --
	//and_combined_filters_can_be_serialized --
	//or_combined_filters_can_be_serialized --
}

func Test_filter_parse(t *testing.T) {
	for _, testData := range canBeSerializedTestDataSlice {
		c, err := jsonpath.WhereString(testData.WhereString)
		if err == nil {
			switch testData.Operator {
			case gte:
				c, err = c.Gte(testData.Value)
			case lte:
				c, err = c.Lte(testData.Value)
			case eq:
				c, err = c.Eq(testData.Value)
			case ne:
				c, err = c.Ne(testData.Value)
			case lt:
				c, err = c.Lt(testData.Value)
			case gt:
				c, err = c.Gt(testData.Value)
			case nin:
				c, err = c.Nin(testData.Value)
			case in:
				if reflect.ValueOf(testData.Value).Kind() == reflect.Slice {
					c, err = c.InSlice(testData.Value)
				} else {
					c, err = c.In(testData.Value)
				}
			case contains:
				c, err = c.Contains(testData.Value)
			case all:
				c, err = c.AllSlice(testData.Value)
			case size:
				s, ok := testData.Value.(int)
				if !ok {
					t.Errorf("size should be a int")
				}
				c, err = c.Size(s)
			case subSetOf:
				c, err = c.SubSetOfSlice(testData.Value)
			case anyOf:
				c, err = c.AnyOfSlice(testData.Value)
			case noneOf:
				c, err = c.NoneOfSlice(testData.Value)
			case exists:
				expected, _ := testData.Value.(bool)
				c, err = c.Exists(expected)
			case empty:
				expected, _ := testData.Value.(bool)
				c = c.Empty(expected)
			}

			if err == nil {
				f := jsonpath.CreateSingleFilter(c)
				fString := f.String()
				parsed, err1 := filter.Compile(testData.ParseString)
				if err1 == nil {
					parsedString := parsed.String()
					if fString != parsedString {
						t.Errorf("failed")
					}
				} else {
					t.Errorf(err1.Error())
				}
			} else {
				t.Errorf(err.Error())
			}
		} else {
			t.Errorf(err.Error())
		}
	}
}

func Test_a_matches_filter_can_be_serialized(t *testing.T) {
	x, err0 := jsonpath.WhereString("x")
	if err0 == nil {
		x, err0 = x.Eq(1000)
		c, err := jsonpath.WhereString("a")
		if err == nil {
			c = c.Matches(x)
			f := jsonpath.CreateSingleFilter(c)
			fString := f.String()

			parsedString := "[?(@['a'] MATCHES [?(@['x'] == 1000)])]"
			if fString != parsedString {
				t.Errorf("failed")
			}
		} else {
			t.Errorf(err.Error())
		}
	} else {
		t.Errorf(err0.Error())
	}
}

func Test_and_filter_can_be_serialized(t *testing.T) {
	c, err := jsonpath.WhereString("a")
	if err == nil {
		c, err = c.Eq(1)
		if err == nil {
			c, err = c.And("b")
			if err == nil {
				c, err = c.Eq(2)
				if err == nil {
					f := jsonpath.CreateSingleFilter(c)
					fString := f.String()
					parsed, err1 := filter.Compile("[?(@['a'] == 1 && @['b'] == 2)]")
					if err1 == nil {
						parsedString := parsed.String()
						if fString != parsedString {
							t.Errorf("failed")
						}
					} else {
						t.Errorf(err1.Error())
					}
				} else {
					t.Errorf(err.Error())
				}
			} else {
				t.Errorf(err.Error())
			}
		} else {
			t.Errorf(err.Error())
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_a_regex_filter_can_be_serialized(t *testing.T) {
	c, err := jsonpath.WhereString("a")
	if err == nil {
		c, err = c.Regex(regexp.MustCompile("(?i).*?"))
		if err == nil {
			f := jsonpath.CreateSingleFilter(c)
			fString := f.String()
			parsedString := "[?(@['a'] =~ /.*?/i)]"

			if fString != parsedString {
				t.Errorf("failed")
			}
		} else {
			t.Errorf(err.Error())
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_a_doc_ref_filter_can_be_serialized(t *testing.T) {
	parsed, err := filter.Compile("[?(@.display-price <= $.max-price)]")
	if err == nil {
		parsedString := parsed.String()
		if parsedString != "[?(@['display-price'] <= $['max-price'])]" {
			t.Errorf("failed")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_and_combined_filters_can_be_serialized(t *testing.T) {
	b, err0 := jsonpath.WhereString("b")
	if err0 == nil {
		b, err0 = b.Eq(2)
		a, err := jsonpath.WhereString("a")
		if err == nil {
			a, err = a.Eq(1)
			if err == nil {
				f := jsonpath.CreateSingleFilter(a)
				c := f.And(b)
				cString := c.String()
				parsed, err1 := filter.Compile("[?(@['a'] == 1 && @['b'] == 2)]")
				if err1 == nil {
					parsedString := parsed.String()
					if cString != parsedString {
						t.Errorf("failed")
					}
				} else {
					t.Errorf(err1.Error())
				}
			} else {
				t.Errorf(err.Error())
			}
		} else {
			t.Errorf(err.Error())
		}
	} else {
		t.Errorf(err0.Error())
	}
}
func Test_or_combined_filters_can_be_serialized(t *testing.T) {
	b, err0 := jsonpath.WhereString("b")
	if err0 == nil {
		b, err0 = b.Eq(2)
		a, err := jsonpath.WhereString("a")
		if err == nil {
			a, err = a.Eq(1)
			if err == nil {
				f := jsonpath.CreateSingleFilter(a)
				c := f.Or(b)
				cString := c.String()
				parsed, err1 := filter.Compile("[?(@['a'] == 1 || @['b'] == 2)]")
				if err1 == nil {
					parsedString := parsed.String()
					if cString != parsedString {
						t.Errorf("failed")
					}
				} else {
					t.Errorf(err1.Error())
				}
			} else {
				t.Errorf(err.Error())
			}
		} else {
			t.Errorf(err.Error())
		}
	} else {
		t.Errorf(err0.Error())
	}
}
