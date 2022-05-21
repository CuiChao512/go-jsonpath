package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath"
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"reflect"
	"testing"
)

type inlineTestMetaData struct {
	JsonString string
	PathString string
	Function   func(interface{}) interface{}
	Expected   interface{}
}

var (
	bookCount = 4

	inlineTestMetaDataTable = []inlineTestMetaData{
		//root_context_can_be_referred_in_predicate
		{
			PathString: "store.book[?(@.display-price <= $.max-price)].display-price",
			Expected:   []interface{}{8.95, 8.99},
		},
		//multiple_context_object_can_be_referred
		{
			PathString: "store.book[ ?(@.category == @.category) ]",
			Function:   sizeOf,
			Expected:   bookCount,
		},
		{
			PathString: "store.book[ ?(@.category == @['category']) ]",
			Function:   sizeOf,
			Expected:   bookCount,
		},
		{
			PathString: "store.book[ ?(@ == @) ]",
			Function:   sizeOf,
			Expected:   bookCount,
		},
		{
			PathString: "store.book[ ?(@.category != @.category) ]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			PathString: "store.book[ ?(@.category != @) ]",
			Function:   sizeOf,
			Expected:   bookCount,
		},
		//simple_inline_or_statement_evaluates
		{
			PathString: "store.book[ ?(@.author == 'Nigel Rees' || @.author == 'Evelyn Waugh') ].author",
			Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh"},
		},
		{
			PathString: "store.book[ ?((@.author == 'Nigel Rees' || @.author == 'Evelyn Waugh') && @.display-price < 15) ].author",
			Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh"},
		},
		{
			PathString: "store.book[ ?((@.author == 'Nigel Rees' || @.author == 'Evelyn Waugh') && @.category == 'reference') ].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		{
			PathString: "store.book[ ?((@.author == 'Nigel Rees') || (@.author == 'Evelyn Waugh' && @.category != 'fiction')) ].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		//no_path_ref_in_filter_hit_all
		{
			PathString: "$.store.book[?('a' == 'a')].author",
			Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien"},
		},
		//no_path_ref_in_filter_hit_none
		{
			PathString: "$.store.book[?('a' == 'b')].author",
			Expected:   []interface{}{},
		},
		//path_can_be_on_either_side_of_operator
		{
			PathString: "$.store.book[?(@.category == 'reference')].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		{
			PathString: "$.store.book[?('reference' == @.category)].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		//path_can_be_on_both_side_of_operator
		{
			PathString: "$.store.book[?(@.category == @.category)].author",
			Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien"},
		},
		//patterns_can_be_evaluated
		{
			PathString: "$.store.book[?(@.category =~ /reference/)].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		{
			PathString: "$.store.book[?(/reference/ =~ @.category)].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		//patterns_can_be_evaluated_with_ignore_case
		{
			PathString: "$.store.book[?(@.category =~ /REFERENCE/)].author",
			Expected:   []interface{}{},
		},
		{
			PathString: "$.store.book[?(@.category =~ /REFERENCE/i)].author",
			Expected:   []interface{}{"Nigel Rees"},
		},
		//negate_exists_check
		{
			PathString: "$.store.book[?(@.isbn)].author",
			Expected:   []interface{}{"Herman Melville", "J. R. R. Tolkien"},
		},
		{
			PathString: "$.store.book[?(!@.isbn)].author",
			Expected:   []interface{}{"Nigel Rees", "Evelyn Waugh"},
		},
		//equality_check_does_not_break_evaluation
		{
			JsonString: "[{\"value\":\"5\"}]",
			PathString: "$[?(@.value=='5')]",
			Function:   sizeOf,
			Expected:   1,
		},
		{
			JsonString: "[{\"value\":5}]",
			PathString: "$[?(@.value==5)]",
			Function:   sizeOf,
			Expected:   1,
		},
		{
			JsonString: "[{\"value\":\"5.1.26\"}]",
			PathString: "$[?(@.value=='5.1.26')]",
			Function:   sizeOf,
			Expected:   1,
		},
		{
			JsonString: "[{\"value\":\"5\"}]",
			PathString: "$[?(@.value=='5.1.26')]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":5}]",
			PathString: "$[?(@.value=='5.1.26')]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":5.1}]",
			PathString: "$[?(@.value=='5.1.26')]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":\"5.1.26\"}]",
			PathString: "$[?(@.value=='5')]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":\"5.1.26\"}]",
			PathString: "$[?(@.value==5)]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":\"5.1.26\"}]",
			PathString: "$[?(@.value==5.1)]",
			Function:   sizeOf,
			Expected:   0,
		},
		//lt_check_does_not_break_evaluation
		{
			JsonString: "[{\"value\":\"5\"}]",
			PathString: "$[?(@.value<'7')]",
			Function:   sizeOf,
			Expected:   1,
		},
		{
			JsonString: "[{\"value\":\"7\"}]",
			PathString: "$[?(@.value<'5')]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":5}]",
			PathString: "$[?(@.value<7)]",
			Function:   sizeOf,
			Expected:   1,
		},
		{
			JsonString: "[{\"value\":7}]",
			PathString: "$[?(@.value<5)]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":5}]",
			PathString: "$[?(@.value<7.1)]",
			Function:   sizeOf,
			Expected:   1,
		},
		{
			JsonString: "[{\"value\":7}]",
			PathString: "$[?(@.value<5.1)]",
			Function:   sizeOf,
			Expected:   0,
		},
		{
			JsonString: "[{\"value\":5.1}]",
			PathString: "$[?(@.value<7)]",
			Function:   sizeOf,
			Expected:   1,
		},
		{
			JsonString: "[{\"value\":7.1}]",
			PathString: "$[?(@.value<5)]",
			Function:   sizeOf,
			Expected:   0,
		},
		//escaped_literals
		{
			JsonString: "[\"'foo\"]",
			PathString: "$[?(@ == '\\'foo')]",
			Function:   sizeOf,
			Expected:   1,
		},
		//escaped_literals2
		{
			JsonString: "[\"\\\\'foo\"]",
			PathString: "$[?(@ == \"\\\\'foo\")]",
			Function:   sizeOf,
			Expected:   1,
		},
		//escape_pattern
		{
			JsonString: "[\"x\"]",
			PathString: "$[?(@ =~ /\\/|x/)]",
			Function:   sizeOf,
			Expected:   1,
		},
		//escape_pattern_after_literal
		{
			JsonString: "[\"x\"]",
			PathString: "$[?(@ == \"abc\" || @ =~ /\\/|x/)]",
			Function:   sizeOf,
			Expected:   1,
		},
		//escape_pattern_before_literal
		{
			JsonString: "[\"x\"]",
			PathString: "$[?(@ =~ /\\/|x/ || @ == \"abc\")]",
			Function:   sizeOf,
			Expected:   1,
		},
		//filter_evaluation_does_not_break_path_evaluation
		{
			JsonString: "[{\"s\": \"fo\", \"expected_size\": \"m\"}, {\"s\": \"lo\", \"expected_size\": 2}]",
			PathString: "$[?(@.s size @.expected_size)]",
			Function:   sizeOf,
			Expected:   1,
		},
	}
)

func Test_inline_filters(t *testing.T) {
	for i, testData := range inlineTestMetaDataTable {
		json := TestJsonDocument
		if testData.JsonString != "" {
			json = testData.JsonString
		}
		if document, err := jsonpath.CreateParseContextImplByConfiguration(common.DefaultConfiguration()).ParseString(json); err == nil {
			if result, err := document.Read(testData.PathString); err == nil {
				if testData.Function != nil {
					result = testData.Function(result)
				}
				if !reflect.DeepEqual(result, testData.Expected) {
					t.Errorf("case No.%d failed message:%s", i, err.Error())
				}
			} else {
				t.Errorf("case No.%d failed message:%s", i, err.Error())
			}
		} else {
			t.Errorf("case No.%d failed message:%s", i, err.Error())
		}

	}
}

func Test_negate_exists_check_primitive(t *testing.T) {
	ints := []interface{}{0, 1, nil, 2, 3}
	parsed, err := jsonpath.JsonpathParseObject(ints)
	if err != nil {
		t.Errorf(err.Error())
	}
	hits, err := parsed.Read("$[?(@)]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(hits, []interface{}{0, 1, nil, 2, 3}) {
		t.Errorf("fail")
	}

	parsed, err = jsonpath.JsonpathParseObject(ints)
	if err != nil {
		t.Errorf(err.Error())
	}
	hits, err = parsed.Read("$[?(@ != null)]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(hits, []interface{}{0, 1, 2, 3}) {
		t.Errorf("fail")
	}

	parsed, err = jsonpath.JsonpathParseObject(ints)
	if err != nil {
		t.Errorf(err.Error())
	}
	hits, err = parsed.Read("$[?(!@)]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(hits, []interface{}{}) {
		t.Errorf("fail")
	}

}
