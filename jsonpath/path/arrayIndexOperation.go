package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"regexp"
	"strconv"
	"strings"
)

var REGEXP_COMMA = regexp.MustCompile("\\s*,\\s*")

type ArrayIndexOperation struct {
	indexes []int
}

func (a *ArrayIndexOperation) Indexes() []int {
	return a.indexes
}

func (a *ArrayIndexOperation) IsSingleIndexOperation() bool {
	return len(a.indexes) == 1
}

func (a *ArrayIndexOperation) String() string {
	return "[" + jsonpath.UtilsJoin(",", "", a.indexes) + "]"
}

func ParseArrayIndexOperation(operation string) (*ArrayIndexOperation, error) {
	for _, c := range []rune(operation) {
		if !jsonpath.UtilsCharIsDigit(c) && c != ',' && c != ' ' && c != '-' {
			return nil, &jsonpath.InvalidPathError{Message: "Failed to parse ArrayIndexOperation: " + operation}
		}
	}
	tokens := strings.Split(operation, ",")

	var tempIndexes []int
	for _, token := range tokens {
		i, err := strconv.Atoi(token)
		if err != nil {
			return nil, &jsonpath.InvalidPathError{Message: "Failed to parse token in ArrayIndexOperation: " + token}
		}
		tempIndexes = append(tempIndexes, i)
	}
	a := &ArrayIndexOperation{indexes: tempIndexes}
	return a, nil
}

type ArraySliceOperationType int

const (
	SLICE_FROM    ArraySliceOperationType = 0
	SLICE_TO      ArraySliceOperationType = 1
	SLICE_BETWEEN ArraySliceOperationType = 2
)

type ArraySliceOperation struct {
	from          int
	to            int
	operationType ArraySliceOperationType
}

func (a *ArraySliceOperation) From() int {
	return a.from
}

func (a *ArraySliceOperation) To() int {
	return a.to
}

func (a *ArraySliceOperation) OperationType() ArraySliceOperationType {
	return a.operationType
}
