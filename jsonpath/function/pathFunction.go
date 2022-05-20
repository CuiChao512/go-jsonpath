package function

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type PathFunction interface {
	Invoke(currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters *[]*Parameter) (interface{}, error)
}

var functions map[string]PathFunction

func init() {
	functions["avg"] = &Average{}
	//functions["stddev"] = &StandardDeviation{}
	//functions["sum"] = &Sum{}
	functions["min"] = &Min{}
	functions["max"] = &Max{}
	//functions["concat"] = &Concat{}
	functions["length"] = &Length{}
	//functions["size"] = &Size{}
	functions["append"] = &Append{}
	//functions["keys"] = &Keys{}
}

func GetFunctionByName(name string) (PathFunction, error) {
	f := functions[name]
	if f == nil {
		return nil, &common.InvalidPathError{Message: "Function with name: " + name + " does not exist."}
	}
	return f, nil
}
