package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
	"testing"
)

type inlineTestMetaData struct {
	PathString string
	Function   func(interface{}) interface{}
	Expected   interface{}
}

var (
	bookCount = 4

	inlineTestMetaDataTable = []inlineTestMetaData{
		////root_context_can_be_referred_in_predicate
		//{
		//	PathString: "store.book[?(@.display-price <= $.max-price)].display-price",
		//	Expected:   []interface{}{8.95, 8.99},
		//},
		////multiple_context_object_can_be_referred
		//{
		//	PathString: "store.book[ ?(@.category == @.category) ]",
		//	Function:   sizeOf,
		//	Expected:   bookCount,
		//},
		//{
		//	PathString: "store.book[ ?(@.category == @['category']) ]",
		//	Function:   sizeOf,
		//	Expected:   bookCount,
		//},
		//{
		//	PathString: "store.book[ ?(@ == @) ]",
		//	Function:   sizeOf,
		//	Expected:   bookCount,
		//},
		//{
		//	PathString: "store.book[ ?(@.category != @.category) ]",
		//	Function:   sizeOf,
		//	Expected:   0,
		//},
		//{
		//	PathString: "store.book[ ?(@.category != @) ]",
		//	Function:   sizeOf,
		//	Expected:   bookCount,
		//},
		////simple_inline_or_statement_evaluates
		//{
		//	PathString: "store.book[ ?(@.author == 'Nigel Rees' || @.author == 'Evelyn Waugh') ].author",
		//	Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh"},
		//},
		//{
		//	PathString: "store.book[ ?((@.author == 'Nigel Rees' || @.author == 'Evelyn Waugh') && @.display-price < 15) ].author",
		//	Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh"},
		//},
		//{
		//	PathString: "store.book[ ?((@.author == 'Nigel Rees' || @.author == 'Evelyn Waugh') && @.category == 'reference') ].author",
		//	Expected:   []interface{}{"Nigel Rees"},
		//},
		//{
		//	PathString: "store.book[ ?((@.author == 'Nigel Rees') || (@.author == 'Evelyn Waugh' && @.category != 'fiction')) ].author",
		//	Expected:   []interface{}{"Nigel Rees"},
		//},
		////no_path_ref_in_filter_hit_all
		//{
		//	PathString: "$.store.book[?('a' == 'a')].author",
		//	Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien"},
		//},
		//no_path_ref_in_filter_hit_none
		//{
		//	PathString: "$.store.book[?('a' == 'b')].author",
		//	Expected:   []interface{}{},
		//},
		////path_can_be_on_either_side_of_operator
		//{
		//	PathString: "$.store.book[?(@.category == 'reference')].author",
		//	Expected:   []interface{}{"Nigel Rees"},
		//},
		//{
		//	PathString: "$.store.book[?('reference' == @.category)].author",
		//	Expected:   []interface{}{"Nigel Rees"},
		//},
		////path_can_be_on_both_side_of_operator
		//{
		//	PathString: "$.store.book[?(@.category == @.category)].author",
		//	Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien"},
		//},
		////patterns_can_be_evaluated
		//{
		//	PathString: "$.store.book[?(@.category =~ /reference/)].author",
		//	Expected:   []interface{}{"Nigel Rees"},
		//},
		//{
		//	PathString: "$.store.book[?(/reference/ =~ @.category)].author",
		//	Expected:   []interface{}{"Nigel Rees"},
		//},
		//patterns_can_be_evaluated_with_ignore_case
		{
			PathString: "$.store.book[?(@.category =~ /REFERENCE/)].author",
			Expected:   []interface{}{},
		},
		{
			PathString: "$.store.book[?(@.category =~ /REFERENCE/i)].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		//
	}
)

func Test_inline_filters(t *testing.T) {
	for _, testData := range inlineTestMetaDataTable {
		if document, err := jsonpath.CreateParseContextImplByConfiguration(common.DefaultConfiguration()).ParseString(TestJsonDocument); err == nil {
			if result, err := document.Read(testData.PathString); err == nil {
				if testData.Function != nil {
					result = testData.Function(result)
				}
				if !reflect.DeepEqual(result, testData.Expected) {
					t.Errorf("fail")
				}
			} else {
				t.Errorf(err.Error())
			}
		} else {
			t.Errorf(err.Error())
		}

	}
}
