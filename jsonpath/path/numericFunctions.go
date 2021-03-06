package path

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"github.com/CuiChao512/go-jsonpath/jsonpath/function"
	"math"
	"strings"
)

type defaultInvoker struct {
}

func (a *defaultInvoker) Invoke(nextAndGet PathFunctionNextAndGet, currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters []*function.Parameter) (interface{}, error) {
	count := 0
	if ctx.Configuration().JsonProvider().IsArray(model) {

		objects, err := ctx.Configuration().JsonProvider().ToArray(model)
		if err != nil {
			return nil, err
		}
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
				nextAndGet.Next(obj)
			}
		}
	}
	if parameters != nil {
		values, err := function.ParametersToList(common.TYPE_NUMBER, ctx, parameters)
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			count++
			nextAndGet.Next(value)
		}
	}
	if count != 0 {
		return nextAndGet.GetValue(), nil
	}
	return nil, &common.JsonPathError{Message: "Aggregation function attempted to calculate value using empty array"}
}

// Average function

type Average struct {
	*defaultInvoker
	summation float64
	count     int
}

func (a *Average) Next(value interface{}) {
	a.count++
	v, _ := common.UtilsNumberToFloat64(value)
	a.summation += v
}

func (a *Average) GetValue() interface{} {
	if a.count != 0 {
		return a.summation / float64(a.count)
	}
	return 0
}

//Max function
type Max struct {
	*defaultInvoker
	max float64
}

func (m *Max) Next(value interface{}) {
	v := common.UtilsNumberToFloat64Force(value)
	if m.max < v {
		m.max = v
	}
}

func (m *Max) GetValue() interface{} {
	return m.max
}

func CreateMaxFunction() *Max {
	return &Max{max: math.MinInt64}
}

// Min function
type Min struct {
	*defaultInvoker
	min float64
}

func (m *Min) Next(value interface{}) {
	v := common.UtilsNumberToFloat64Force(value)
	if m.min > v {
		m.min = v
	}
}

func (m *Min) GetValue() interface{} {
	return m.min
}

func CreateMinFunction() *Min {
	return &Min{min: math.MaxInt64}
}

// StandardDeviation ---
type StandardDeviation struct {
	*defaultInvoker
	sumSq float64
	sum   float64
	count int64
}

func (s *StandardDeviation) Next(value interface{}) {
	v := common.UtilsNumberToFloat64Force(value)
	s.sum += v
	s.sumSq += v * v
	s.count++
}

func (s *StandardDeviation) GetValue() interface{} {
	count := float64(s.count)
	return math.Sqrt(s.sumSq/count - s.sum*s.sum/count/count)
}

// Sum function
type Sum struct {
	*defaultInvoker
	sum float64
}

func (s *Sum) Next(value interface{}) {
	v := common.UtilsNumberToFloat64Force(value)
	s.sum += v
}

func (s *Sum) GetValue() interface{} {
	return s.sum
}

type Length struct {
}

func (*Length) Next(value interface{}) {}

func (*Length) GetValue() interface{} { return nil }

func (a *Length) Invoke(nextAndGet PathFunctionNextAndGet, currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters []*function.Parameter) (interface{}, error) {
	if parameters != nil && len(parameters) > 0 {

		// Set the tail of the first parameter, when its not a function path parameter (which wouldn't make sense
		// for length - to the wildcard such that we request all of its children so we can get back an array and
		// take its length.
		lengthOfParameter := parameters[0]
		if !lengthOfParameter.GetPath().IsFunctionPath() {
			path := lengthOfParameter.GetPath()
			switch path.(type) {
			case *CompiledPath:
				cp, _ := path.(*CompiledPath)
				root := cp.GetRoot()
				tail := root.GetNext()
				for nil != tail && nil != tail.GetNext() {
					tail = tail.GetNext()
				}
				if nil != tail {
					tail.SetNext(CreateWildcardPathToken())
				}
			}
		}
		result, err := parameters[0].GetPath().Evaluate(model, model, ctx.Configuration())
		if err != nil {
			return nil, err
		}
		innerModel, err1 := result.GetValue()
		if err1 != nil {
			return nil, err1
		}
		if ctx.Configuration().JsonProvider().IsArray(innerModel) {
			return ctx.Configuration().JsonProvider().Length(innerModel)
		}
	}
	if ctx.Configuration().JsonProvider().IsArray(model) {
		return ctx.Configuration().JsonProvider().Length(model)
	} else if ctx.Configuration().JsonProvider().IsMap(model) {
		return ctx.Configuration().JsonProvider().Length(model)
	}
	return nil, nil
}

type Concatenate struct {
}

func (*Concatenate) Next(value interface{}) {}

func (*Concatenate) GetValue() interface{} { return nil }

func (a *Concatenate) Invoke(nextAndGet PathFunctionNextAndGet, currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters []*function.Parameter) (interface{}, error) {
	result := &strings.Builder{}
	if ctx.Configuration().JsonProvider().IsArray(model) {
		objects, _ := ctx.Configuration().JsonProvider().ToArray(model)
		for _, obj := range objects {
			if str, ok := obj.(string); ok {
				result.WriteString(str)
			}
		}
	}
	if parameters != nil {
		parametersList, err := function.ParametersToList(common.TYPE_STRING, ctx, parameters)
		if err == nil {
			for _, value := range parametersList {
				result.WriteString(common.UtilsToString(value))
			}
		}
	}
	return result.String(), nil
}
