package function

import (
	"cuichao.com/go-jsonpath/jsonpath/evaluationContext"
	"cuichao.com/go-jsonpath/jsonpath/path"
)

type KeySetFunction struct {
}

func (*KeySetFunction) Invoke(currentPath string, parent path.Ref, model interface{}, ctx evaluationContext.EvaluationContext, parameters *[]*Parameter) (interface{}, error) {
	if ctx.Configuration().JsonProvider().IsMap(model) {
		return ctx.Configuration().JsonProvider().GetPropertyKeys(model), nil
	}
	return nil, nil
}

type Append struct {
}

func (*Append) Invoke(currentPath string, parent path.Ref, model interface{}, ctx evaluationContext.EvaluationContext, parameters *[]*Parameter) (interface{}, error) {
	jsonProvider := ctx.Configuration().JsonProvider()
	if parameters != nil && len(*parameters) > 0 {
		for _, param := range *parameters {
			if jsonProvider.IsArray(model) {
				l := jsonProvider.Length(model)
				val, err := param.GetValue()
				if err != nil {
					return nil, err
				}
				jsonProvider.SetArrayIndex(model, l, val)
			}
		}
	}
	return model, nil
}
