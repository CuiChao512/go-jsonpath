package common

type EvaluationContinuation int

const (
	CONTINUE EvaluationContinuation = 0
	ABORT    EvaluationContinuation = 1
)

type EvaluationListener interface {
	ResultFound(found FoundResult) EvaluationContinuation
}

type FoundResult interface {
	Index() int
	Path() string
	Result() interface{}
}
