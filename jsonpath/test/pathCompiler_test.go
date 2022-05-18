package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/filter"
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
