package jsonpath

type EvaluationContext interface {
	Configuration() *Configuration
	RootDocument() interface{}
	GetValue() interface{}
	GetValueUnwrap(unwrap bool) interface{}
}

type EvaluationContextImpl struct {
}

func (e *EvaluationContextImpl) Configuration() *Configuration {
	return nil
}
func (e *EvaluationContextImpl) RootDocument() interface{} {
	return nil
}

func (e *EvaluationContextImpl) GetValue() interface{} {
	return nil
}

func (e *EvaluationContextImpl) GetValueUnwrap(unwrap bool) interface{} {
	return nil
}
