package numeric

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/function"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type abstractAggregation struct {
}

func (*abstractAggregation) next(value interface{}) {}

func (*abstractAggregation) GetValue() interface{} { return nil }

func (a *abstractAggregation) Invoke(currentPath string, parent path.Ref, model interface{}, ctx jsonpath.EvaluationContext, parameters []*function.Parameter) (interface{}, error) {
	count := 0
	if ctx.Configuration().JsonProvider().IsArray(model) {

		objects := ctx.Configuration().JsonProvider().ToIterable(model)
		for _, obj := range objects {
			isNumber := false
			switch obj.(type) {
			case int:
				isNumber = true
			case float64:
				isNumber = true
			case float32:
				isNumber = true
			}
			if isNumber {
				count++
				a.next(obj)
			}
		}
	}
	if parameters != nil {
		values, err := function.ParametersToList(jsonpath.TYPE_NUMBER, ctx, parameters)
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			count++
			a.next(value)
		}
	}
	if count != 0 {
		return a.GetValue(), nil
	}
	return nil, &jsonpath.JsonPathError{Message: "Aggregation function attempted to calculate value using empty array"}
}
