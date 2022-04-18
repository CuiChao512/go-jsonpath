package test

import "cuichao.com/go-jsonpath/jsonpath/common"

var FilterTestJson, err = common.DefaultConfiguration().JsonProvider().Parse("{" +
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

func intEqEvals_test() {

}
