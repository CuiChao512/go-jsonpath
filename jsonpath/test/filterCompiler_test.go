package test

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/filter"
	"strconv"
	"testing"
)

type validFilterTestData struct {
	FilterString           string
	FilterToStringExpected string
}

var validFilterTestDataSlice = []validFilterTestData{
	{FilterString: "[?(@)]", FilterToStringExpected: "[?(@)]"},
	{FilterString: "[?(@.firstname)]", FilterToStringExpected: "[?(@['firstname'])]"},
	{FilterString: "[?($.firstname)]", FilterToStringExpected: "[?($['firstname'])]"},
	{FilterString: "[?(@['firstname'])]", FilterToStringExpected: "[?(@['firstname'])]"},
	{FilterString: "[?($['firstname'].lastname)]", FilterToStringExpected: "[?($['firstname']['lastname'])]"},
	{FilterString: "[?($['firstname']['lastname'])]", FilterToStringExpected: "[?($['firstname']['lastname'])]"},
	{FilterString: "[?($['firstname']['lastname'].*)]", FilterToStringExpected: "[?($['firstname']['lastname'][*])]"},
	{FilterString: "[?($['firstname']['num_eq'] == 1)]", FilterToStringExpected: "[?($['firstname']['num_eq'] == 1)]"},
	{FilterString: "[?($['firstname']['num_gt'] > 1.1)]", FilterToStringExpected: "[?($['firstname']['num_gt'] > 1.1)]"},
	{FilterString: "[?($['firstname']['num_lt'] < 11.11)]", FilterToStringExpected: "[?($['firstname']['num_lt'] < 11.11)]"},
	{FilterString: "[?($['firstname']['str_eq'] == 'hej')]", FilterToStringExpected: "[?($['firstname']['str_eq'] == 'hej')]"},
	{FilterString: "[?($['firstname']['str_eq'] == '')]", FilterToStringExpected: "[?($['firstname']['str_eq'] == '')]"},
	{FilterString: "[?($['firstname']['str_eq'] == null)]", FilterToStringExpected: "[?($['firstname']['str_eq'] == null)]"},
	{FilterString: "[?($['firstname']['str_eq'] == true)]", FilterToStringExpected: "[?($['firstname']['str_eq'] == true)]"},
	{FilterString: "[?($['firstname']['str_eq'] == false)]", FilterToStringExpected: "[?($['firstname']['str_eq'] == false)]"},
	{FilterString: "[?(@.firstname && @.lastname)]", FilterToStringExpected: "[?(@['firstname'] && @['lastname'])]"},
	{FilterString: "[?((@.firstname || @.lastname) && @.and)]", FilterToStringExpected: "[?((@['firstname'] || @['lastname']) && @['and'])]"},
	{FilterString: "[?((@.a || @.b || @.c) && @.x)]", FilterToStringExpected: "[?((@['a'] || @['b'] || @['c']) && @['x'])]"},
	{FilterString: "[?((@.a && @.b && @.c) || @.x)]", FilterToStringExpected: "[?((@['a'] && @['b'] && @['c']) || @['x'])]"},
	{FilterString: "[?((@.a && @.b || @.c) || @.x)]", FilterToStringExpected: "[?(((@['a'] && @['b']) || @['c']) || @['x'])]"},
	{FilterString: "[?((@.a && @.b) || (@.c && @.d))]", FilterToStringExpected: "[?((@['a'] && @['b']) || (@['c'] && @['d']))]"},
	{FilterString: "[?(@.a IN [1,2,3])]", FilterToStringExpected: "[?(@['a'] IN [1,2,3])]"},
	{FilterString: "[?(@.a IN {'foo':'bar'})]", FilterToStringExpected: "[?(@['a'] IN {'foo':'bar'})]"},
	{FilterString: "[?(@.value<'7')]", FilterToStringExpected: "[?(@['value'] < '7')]"},
	{FilterString: "[?(@.message == 'it\\\\')]", FilterToStringExpected: "[?(@['message'] == 'it\\\\')]"},
	{FilterString: "[?(@.message.min() > 10)]", FilterToStringExpected: "[?(@['message'].min() > 10)]"},
	{FilterString: "[?(@.message.min()==10)]", FilterToStringExpected: "[?(@['message'].min() == 10)]"},
	{FilterString: "[?(10 == @.message.min())]", FilterToStringExpected: "[?(10 == @['message'].min())]"},
	{FilterString: "[?(((@)))]", FilterToStringExpected: "[?(@)]"},
	{FilterString: "[?(@.name =~ /.*?/i)]", FilterToStringExpected: "[?(@['name'] =~ /.*?/i)]"},
	{FilterString: "[?(@.name =~ /.*?/)]", FilterToStringExpected: "[?(@['name'] =~ /.*?/)]"},
	{FilterString: "[?($[\"firstname\"][\"lastname\"])]", FilterToStringExpected: "[?($[\"firstname\"][\"lastname\"])]"},
	{FilterString: "[?($[\"firstname\"].lastname)]", FilterToStringExpected: "[?($[\"firstname\"]['lastname'])]"},
	{FilterString: "[?($[\"firstname\", \"lastname\"])]", FilterToStringExpected: "[?($[\"firstname\",\"lastname\"])]"},
	{FilterString: "[?(((@.a && @.b || @.c)) || @.x)]", FilterToStringExpected: "[?(((@['a'] && @['b']) || @['c']) || @['x'])]"},
	//string_quote_style_is_serialized
	{FilterString: "[?('apa' == 'apa')]", FilterToStringExpected: "[?('apa' == 'apa')]"},
	{FilterString: "[?('apa' == \"apa\")]", FilterToStringExpected: "[?('apa' == \"apa\")]"},
	//string_can_contain_path_chars
	{FilterString: "[?(@[')]@$)]'] == ')]@$)]')]", FilterToStringExpected: "[?(@[')]@$)]'] == ')]@$)]')]"},
	{FilterString: "[?(@[\")]@$)]\"] == \")]@$)]\")]", FilterToStringExpected: "[?(@[\")]@$)]\"] == \")]@$)]\")]"},
	//or_has_lower_priority_than_and
	{
		FilterString:           "[?(@.category == 'fiction' && @.author == 'Evelyn Waugh' || @.price > 15)]",
		FilterToStringExpected: "[?((@['category'] == 'fiction' && @['author'] == 'Evelyn Waugh') || @['price'] > 15)]",
	},
	//compile_and_serialize_not_exists_filter
	{
		FilterString:           "[?(!@.foo)]",
		FilterToStringExpected: "[?(!@['foo'])]",
	},
}

func Test_valid_filters_compile(t *testing.T) {
	for i, testData := range validFilterTestDataSlice {
		compiledFilter, err := filter.Compile(testData.FilterString)
		if err != nil {
			println(strconv.Itoa(i) + " failed")
			t.Errorf(err.Error())
		} else {
			str := compiledFilter.String()
			if str != testData.FilterToStringExpected {
				println(strconv.Itoa(i) + " failed")
				t.Errorf("failed, expected:%s actual:%s", testData.FilterToStringExpected, str)
			}
		}
	}
}

var invalidFilterTestDataSlice = []string{
	"[?(@.foo == x)]",
	"[?(@))]",
	"[?(@ FOO 1)]",
	"[?(@ || )]",
	"[?(@ == 'foo )]",
	"[?(@ == 1' )]",
	"[?(@.foo bar == 1)]",
	"[?(@.i == 5 @.i == 8)]",
	"[?(!5)]",
	"[?(!'foo')]",
}

func Test_invalid_filter(t *testing.T) {
	for _, testFilterString := range invalidFilterTestDataSlice {
		if fc, err := filter.Compile(testFilterString); err == nil {
			println("filterCompiled:" + fc.String())
			t.Errorf("shuould throw invalid path error")
		} else {
			switch err.(type) {
			case *common.InvalidPathError:
			default:
				t.Errorf("shuould throw path not found error")
			}
		}
	}
}
