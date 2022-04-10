package jsonpath

type InvalidPathError struct {
	Message string
}

func (e *InvalidPathError) Error() string {
	return e.Message
}

type InvalidJsonError struct {
	Message string
}

func (e *InvalidJsonError) Error() string {
	return e.Message
}

type PathNotFoundError struct {
	Message string
}

func (e *PathNotFoundError) Error() string {
	return e.Message
}

type EvaluationAbortError struct {
	Message string
}

func (e *EvaluationAbortError) Error() string {
	return e.Message
}

type InvalidModificationError struct {
	Message string
}

func (e *InvalidModificationError) Error() string {
	return e.Message
}

type InvalidCriteriaError struct {
	Message string
}

func (e *InvalidCriteriaError) Error() string {
	return e.Message
}

type ValueCompareError struct {
	Message string
}

func (e *ValueCompareError) Error() string {
	return e.Message
}

type JsonPathError struct {
	Message string
}

func (e *JsonPathError) Error() string {
	return e.Message
}

type IllegalStateException struct {
	Message string
}

func (e *IllegalStateException) Error() string {
	return e.Message
}
