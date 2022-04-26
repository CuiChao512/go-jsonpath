package path

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
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
	return "[" + common.UtilsJoin(",", "", a.indexes) + "]"
}

func ParseArrayIndexOperation(operation string) (*ArrayIndexOperation, error) {
	for _, c := range []rune(operation) {
		if !common.UtilsCharIsDigit(c) && c != ',' && c != ' ' && c != '-' {
			return nil, &common.InvalidPathError{Message: "Failed to parse ArrayIndexOperation: " + operation}
		}
	}
	tokens := strings.Split(operation, ",")

	var tempIndexes []int
	for _, token := range tokens {
		i, err := strconv.Atoi(strings.TrimSpace(token))
		if err != nil {
			return nil, &common.InvalidPathError{Message: "Failed to parse token in ArrayIndexOperation: " + token}
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

func tryRead(tokens []string, idx int) (bool, int, error) {
	if len(tokens) > idx {
		if tokens[idx] == "" {
			return false, 0, nil
		}
		intR, err := strconv.Atoi(tokens[idx])
		return true, intR, err
	} else {
		return false, 0, nil
	}
}

func ParseArraySliceOperation(operation string) (*ArraySliceOperation, error) {
	for i := 0; i < len(operation); i++ {
		c := []rune(operation)[i]
		if !common.UtilsCharIsDigit(c) && c != '-' && c != ':' {
			return nil, &common.InvalidPathError{Message: "Failed to parse SliceOperation: " + operation}
		}
	}
	tokens := strings.Split(operation, ":")

	tempFromSuccess, tempFrom, err := tryRead(tokens, 0)
	if err != nil {
		return nil, err
	}
	tempToSuccess, tempTo, err := tryRead(tokens, 1)
	if err != nil {
		return nil, err
	}
	var tempOperation ArraySliceOperationType

	if tempFromSuccess && !tempToSuccess {
		tempOperation = SLICE_FROM
	} else if tempFromSuccess {
		tempOperation = SLICE_BETWEEN
	} else if tempToSuccess {
		tempOperation = SLICE_TO
	} else {
		return nil, &common.InvalidPathError{Message: "Failed to parse SliceOperation: " + operation}
	}

	return &ArraySliceOperation{from: tempFrom, to: tempTo, operationType: tempOperation}, nil
}
