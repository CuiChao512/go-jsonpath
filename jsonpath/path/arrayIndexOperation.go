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
