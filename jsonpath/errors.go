package jsonpath

type InvalidPathError struct {
	Message string
}

func (e *InvalidPathError) Error() string {
	return e.Message
}
