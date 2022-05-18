package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/filter"
	"reflect"
	"testing"
)

var invalidPathTestPath = []string{
	"$X",             //a_root_path_must_be_followed_by_period_or_bracket
	"$.",             //a_path_may_not_end_with_period
	"$.prop.",        //a_path_may_not_end_with_period_2
	"$..",            //a_path_may_not_end_with_scan
	"$.prop..",       //a_path_may_not_end_with_scan_2
	"$..",            //a_path_may_not_end_with_scan
	"$.prop..",       //a_path_may_not_end_with_scan_2
	"$.foo bar",      //a_property_may_not_contain_blanks
	"$[0, 1, 2 4]",   //array_indexes_must_be_separated_by_commas
	"$['1','2',]",    //trailing_comma_after_list_is_not_accepted
	"$['1', ,'3']",   //accept_only_a_single_comma_between_indexes
	"$['aaa'}'bbb']", //property_must_be_separated_by_commas
}

func Test_a_path_must_start_with_dollar_or_at(t *testing.T) {
	for _, testPath := range invalidPathTestPath {
		_, err := filter.PathCompile(testPath)
		if err == nil {
			t.Errorf("shuould throw invalid path error")
		} else {
			switch err.(type) {
			case *common.InvalidPathError:
			default:
				t.Errorf("shuould throw invalid path error")
			}
		}
	}
}

type pathToStringTestData struct {
	PathString     string
	ToStringExpect string
}

var pathStringTestDataSlice = []pathToStringTestData{
	//a_root_path_can_be_compiled
	{
		PathString:     "$",
		ToStringExpect: "$",
	},
	{
		PathString:     "@",
		ToStringExpect: "@",
	},
	//a_property_token_can_be_compiled
	{
		PathString:     "$.prop",
		ToStringExpect: "$['prop']",
	},
	{
		PathString:     "$.1prop",
		ToStringExpect: "$['1prop']",
	},
	{
		PathString:     "$.@prop",
		ToStringExpect: "$['@prop']",
	},
	//a_bracket_notation_property_token_can_be_compiled
	{
		PathString:     "$['prop']",
		ToStringExpect: "$['prop']",
	},
	{
		PathString:     "$['1prop']",
		ToStringExpect: "$['1prop']",
	},
	{
		PathString:     "$['@prop']",
		ToStringExpect: "$['@prop']",
	},
	{
		PathString:     "$[  '@prop'  ]",
		ToStringExpect: "$['@prop']",
	},
	{
		PathString:     "$[\"prop\"]",
		ToStringExpect: "$[\"prop\"]",
	},
	//a_multi_property_token_can_be_compiled
	{
		PathString:     "$['prop0', 'prop1']",
		ToStringExpect: "$['prop0','prop1']",
	},
	{
		PathString:     "$[  'prop0'  , 'prop1'  ]",
		ToStringExpect: "$['prop0','prop1']",
	},
	//a_property_chain_can_be_compiled
	{
		PathString:     "$.abc",
		ToStringExpect: "$['abc']",
	},
	{
		PathString:     "$.aaa.bbb",
		ToStringExpect: "$['aaa']['bbb']",
	},
	{
		PathString:     "$.aaa.bbb.ccc",
		ToStringExpect: "$['aaa']['bbb']['ccc']",
	},
	//a_wildcard_can_be_compiled
	{
		PathString:     "$.*",
		ToStringExpect: "$[*]",
	},
	{
		PathString:     "$[*]",
		ToStringExpect: "$[*]",
	},
	{
		PathString:     "$[ * ]",
		ToStringExpect: "$[*]",
	},
	//a_wildcard_can_follow_a_property
	{
		PathString:     "$.prop[*]",
		ToStringExpect: "$['prop'][*]",
	},
	{
		PathString:     "$['prop'][*]",
		ToStringExpect: "$['prop'][*]",
	},
	//an_array_index_path_can_be_compiled
	{
		PathString:     "$[1]",
		ToStringExpect: "$[1]",
	},
	{
		PathString:     "$[1,2,3]",
		ToStringExpect: "$[1,2,3]",
	},
	{
		PathString:     "$[ 1 , 2 , 3 ]",
		ToStringExpect: "$[1,2,3]",
	},
	//an_array_slice_path_can_be_compiled
	{
		PathString:     "$[-1:]",
		ToStringExpect: "$[-1:]",
	},
	{
		PathString:     "$[1:2]",
		ToStringExpect: "$[1:2]",
	},
	{
		PathString:     "$[:2]",
		ToStringExpect: "$[:2]",
	},
	//an_inline_criteria_can_be_parsed
	{
		PathString:     "$[?(@.foo == 'bar')]",
		ToStringExpect: "$[?]",
	},
	{
		PathString:     "$[?(@.foo == \"bar\")]",
		ToStringExpect: "$[?]",
	},
	//a_scan_token_can_be_parsed
	{
		PathString:     "$..['prop']..[*]",
		ToStringExpect: "$..['prop']..[*]",
	},
	//a_function_can_be_compiled
	{
		PathString:     "$.aaa.foo()",
		ToStringExpect: "$['aaa'].foo()",
	},
	{
		PathString:     "$.aaa.foo(5)",
		ToStringExpect: "$['aaa'].foo(...)",
	},
	{
		PathString:     "$.aaa.foo($.bar)",
		ToStringExpect: "$['aaa'].foo(...)",
	},
	{
		PathString:     "$.aaa.foo(5,10,15)",
		ToStringExpect: "$['aaa'].foo(...)",
	},
}

func Test_a_root_path_can_be_compiled(t *testing.T) {
	for _, pathToStringTestData := range pathStringTestDataSlice {
		if path, err := filter.PathCompile(pathToStringTestData.PathString); err != nil {
			t.Errorf(err.Error())
		} else {
			pathToString := path.String()
			if pathToString != pathToStringTestData.ToStringExpect {
				t.Errorf("path %s 's compiled path to string should be %s, actual is %s", pathToStringTestData.PathString, pathToStringTestData.ToStringExpect, pathToString)
			}
		}
	}
}

type falsePredicate struct{}

func (*falsePredicate) Apply(ctx common.PredicateContext) (bool, error) {
	return false, nil
}

func (*falsePredicate) String() string {
	return ""
}

func Test_a_placeholder_criteria_can_be_parsed(t *testing.T) {
	p := &falsePredicate{}
	if path, err := filter.PathCompile("$[?]", p); err != nil {
		t.Errorf(err.Error())
	} else {
		pathToString := path.String()
		if pathToString != "$[?]" {
			t.Errorf("path %s 's compiled path to string should be %s, actual is %s", "$[?]", "$[?]", pathToString)
		}
	}
	if path, err := filter.PathCompile("$[?,?]", p, p); err != nil {
		t.Errorf(err.Error())
	} else {
		pathToString := path.String()
		if pathToString != "$[?,?]" {
			t.Errorf("path %s 's compiled path to string should be %s, actual is %s", "$[?,?]", "$[?,?]", pathToString)
		}
	}
	if path, err := filter.PathCompile("$[?,?,?]", p, p, p); err != nil {
		t.Errorf(err.Error())
	} else {
		pathToString := path.String()
		if pathToString != "$[?,?,?]" {
			t.Errorf("path %s 's compiled path to string should be %s, actual is %s", "$[?,?,?]", "$[?,?,?]", pathToString)
		}
	}
}

type issuePredicateTestData struct {
	Json       string
	PathString string
	Expected   interface{}
}

var issuePredicateTestDataSlice = []issuePredicateTestData{
	//issue_predicate_can_have_escaped_backslash_in_prop
	{
		Json:       "{\n    \"logs\": [\n        {\n            \"message\": \"it\\\\\",\n            \"id\": 2\n        }\n    ]\n}",
		PathString: "$.logs[?(@.message == 'it\\\\')].message",
		Expected:   []interface{}{"it\\"},
	},
	//issue_predicate_can_have_bracket_in_regex
	{
		Json:       "{\n    \"logs\": [\n        {\n            \"message\": \"(it\",\n            \"id\": 2\n        }\n    ]\n}",
		PathString: "$.logs[?(@.message =~ /\\(it/)].message",
		Expected:   []interface{}{"(it"},
	},
	//issue_predicate_can_have_and_in_regex
	{
		Json:       "{\n    \"logs\": [\n        {\n            \"message\": \"it\",\n            \"id\": 2\n        }\n    ]\n}",
		PathString: "$.logs[?(@.message =~ /&&|it/)].message",
		Expected:   []interface{}{"it"},
	},
	//issue_predicate_can_have_and_in_prop
	{
		Json:       "{\n    \"logs\": [\n        {\n            \"message\": \"&& it\",\n            \"id\": 2\n        }\n    ]\n}",
		PathString: "$.logs[?(@.message == '&& it')].message",
		Expected:   []interface{}{"&& it"},
	},
	//issue_predicate_brackets_must_change_priorities
	{
		Json:       "{\n    \"logs\": [\n        {\n            \"id\": 2\n        }\n    ]\n}",
		PathString: "$.logs[?(@.message && (@.id == 1 || @.id == 2))].id",
		Expected:   []interface{}{},
	},
	{
		Json:       "{\n    \"logs\": [\n        {\n            \"id\": 2\n        }\n    ]\n}",
		PathString: "$.logs[?((@.id == 2 || @.id == 1) && @.message)].id",
		Expected:   []interface{}{},
	},
	//issue_predicate_or_has_lower_priority_than_and
	{
		Json:       "{\n    \"logs\": [\n        {\n            \"id\": 2\n        }\n    ]\n}",
		PathString: "$.logs[?(@.x && @.y || @.id)]",
		Expected:   []interface{}{map[string]interface{}{"id": float64(2)}},
	},
	//issue_predicate_can_have_double_quotes
	{
		Json:       `{"logs": [{ "message": "\"it\""}]}`,
		PathString: "$.logs[?(@.message == '\"it\"')].message",
		Expected:   []interface{}{"\"it\""},
	},
	//issue_predicate_can_have_single_quotes
	{
		Json: `{
					"logs": [
						{
							"message": "'it'"
						}
					]
				}`,
		PathString: "$.logs[?(@.message == \"'it'\")].message",
		Expected:   []interface{}{"'it'"},
	},
	//issue_predicate_can_have_single_quotes_escaped
	{
		Json: `{
		    "logs": [
		        {
		            "message": "'it'"
		        }
		    ]
		}`,
		PathString: "$.logs[?(@.message == '\\'it\\'')].message",
		Expected:   []interface{}{"'it'"},
	},
	//issue_predicate_can_have_square_bracket_in_prop
	{
		Json: `{
                    "logs": [
                        {
                            "message": "] it",
                            "id": 2
                        }
                    ]
                }`,
		PathString: "$.logs[?(@.message == '] it')].message",
		Expected:   []interface{}{"] it"},
	},
}

func Test_issue_predicate_can_have_escaped_backslash_in_prop(t *testing.T) {
	for _, testData := range issuePredicateTestDataSlice {
		if documentContext, err := getParseContextUsingDefaultConf().ParseString(testData.Json); err != nil {
			t.Errorf(err.Error())
		} else {
			if result, err1 := documentContext.Read(testData.PathString); err1 != nil {
				t.Errorf(err.Error())
			} else {
				if !reflect.DeepEqual(result, testData.Expected) {
					t.Errorf("failed")
				}
			}
		}
	}
}
