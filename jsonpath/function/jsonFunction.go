package function

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
)

type KeySetFunction struct {
}

func (*KeySetFunction) Invoke(currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters *[]*Parameter) (interface{}, error) {
	if ctx.Configuration().JsonProvider().IsMap(model) {
		return ctx.Configuration().JsonProvider().GetPropertyKeys(model)
	}
	return nil, nil
}

type Append struct {
}

func (*Append) Invoke(currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters []*Parameter) (interface{}, error) {
	jsonProvider := ctx.Configuration().JsonProvider()
	if parameters != nil && len(parameters) > 0 {
		for _, param := range parameters {
			if jsonProvider.IsArray(model) {
				l, err := jsonProvider.Length(model)
				if err != nil {
					return nil, err
				}
				val, err := param.GetValue()
				if err != nil {
					return nil, err
				}
				err = jsonProvider.SetArrayIndex(model, l, val)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return model, nil
}
