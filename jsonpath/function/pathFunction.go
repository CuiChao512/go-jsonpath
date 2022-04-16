package function

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
)

type PathFunction interface {
	Invoke(currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters *[]*Parameter) (interface{}, error)
}

var functions map[string]PathFunction

func init() {
	//TODO:
}

func GetFunctionByName(name string) (PathFunction, error) {
	f := functions[name]
	if f == nil {
		return nil, &common.InvalidPathError{Message: "Function with name: " + name + " does not exist."}
	}
	return f, nil
}
