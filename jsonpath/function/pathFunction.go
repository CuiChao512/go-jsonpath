package function

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type PathFunction interface {
	Invoke(currentPath string, parent path.Ref, model interface{}, ctx jsonpath.EvaluationContext, parameters *[]*Parameter) (interface{}, error)
}

var functions map[string]PathFunction

func init() {
	//TODO:
}

func GetFunctionByName(name string) (PathFunction, error) {
	f := functions[name]
	if f == nil {
		return nil, &jsonpath.InvalidPathError{Message: "Function with name: " + name + " does not exist."}
	}
	return f, nil
}
