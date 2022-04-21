package function

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"math"
)

type abstractAggregation struct {
}

func (*abstractAggregation) Next(value interface{}) {}

func (*abstractAggregation) GetValue() interface{} { return nil }

func (a *abstractAggregation) Invoke(currentPath string, parent common.PathRef, model interface{}, ctx common.EvaluationContext, parameters *[]*Parameter) (interface{}, error) {
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
				a.Next(obj)
			}
		}
	}
	if parameters != nil {
		values, err := ParametersToList(common.TYPE_NUMBER, ctx, *parameters)
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			count++
			a.Next(value)
		}
	}
	if count != 0 {
		return a.GetValue(), nil
	}
	return nil, &common.JsonPathError{Message: "Aggregation function attempted to calculate value using empty array"}
}

// Average function

type Average struct {
	*abstractAggregation
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
	*abstractAggregation
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
	*abstractAggregation
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
	sum float64
}

func (s *Sum) Next(value interface{}) {
	v := common.UtilsNumberToFloat64Force(value)
	s.sum += v
}

func (s *Sum) GetValue() interface{} {
	return s.sum
}
