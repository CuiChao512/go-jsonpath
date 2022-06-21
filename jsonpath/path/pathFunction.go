package path

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/function"
)

type PathFunction interface {
	PathFunctionNextAndGet
	PathFunctionInvoker
}

type PathFunctionNextAndGet interface {
	Next(value interface{})
	GetValue() interface{}
}

type PathFunctionInvoker interface {
	Invoke(nextAndGet PathFunctionNextAndGet, currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters []*function.Parameter) (interface{}, error)
}

func GetFunctionByName(name string) (PathFunction, error) {
	var f PathFunction
	switch name {
	case "avg":
		f = &Average{}
	case "stddev":
		f = &StandardDeviation{}
	case "sum":
		f = &Sum{}
	case "min":
		f = CreateMinFunction()
	case "max":
		f = CreateMaxFunction()
	case "concat":
		f = &Concatenate{}
	case "length":
		f = &Length{}
	case "size":
		f = &Length{}
	case "append":
		f = &Append{}
	case "keys":
		f = &KeySetFunction{}
	default:
		return nil, &common.InvalidPathError{Message: "Function with name: " + name + " does not exist."}
	}

	return f, nil
}
