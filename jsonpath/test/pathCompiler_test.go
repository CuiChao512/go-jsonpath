package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/filter"
	"testing"
)

var invalidPathTestPath = []string{
	"$X",       //a_root_path_must_be_followed_by_period_or_bracket
	"$.",       //a_path_may_not_end_with_period
	"$.prop.",  //a_path_may_not_end_with_period_2
	"$..",      //a_path_may_not_end_with_scan
	"$.prop..", //a_path_may_not_end_with_scan_2
	"$..",      //a_path_may_not_end_with_scan
	"$.prop..", //a_path_may_not_end_with_scan_2
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
	////a_root_path_can_be_compiled
	//{
	//	PathString:     "$",
	//	ToStringExpect: "$",
	//},
	//{
	//	PathString:     "@",
	//	ToStringExpect: "@",
	//},
	////a_property_token_can_be_compiled
	//{
	//	PathString:     "$.prop",
	//	ToStringExpect: "$['prop']",
	//},
	//{
	//	PathString:     "$.1prop",
	//	ToStringExpect: "$['1prop']",
	//},
	//{
	//	PathString:     "$.@prop",
	//	ToStringExpect: "$['@prop']",
	//},
	////a_bracket_notation_property_token_can_be_compiled
	//{
	//	PathString:     "$['prop']",
	//	ToStringExpect: "$['prop']",
	//},
	//{
	//	PathString:     "$['1prop']",
	//	ToStringExpect: "$['1prop']",
	//},
	//{
	//	PathString:     "$['@prop']",
	//	ToStringExpect: "$['@prop']",
	//},
	{
		PathString:     "$[  '@prop'  ]",
		ToStringExpect: "$['@prop']",
	},
	//{
	//	PathString:     "$[\"prop\"]",
	//	ToStringExpect: "$[\"prop\"]",
	//},
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
