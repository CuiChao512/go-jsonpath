package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"strings"
)

type Filter interface {
	common.Predicate
	Or(other common.Predicate) *OrFilter
	And(other common.Predicate) *AndFilter
}

type FilterImpl struct {
}

func (filter *FilterImpl) String() string {
	return ""
}

func (filter *FilterImpl) Apply(ctx common.PredicateContext) bool {
	return false
}

func (filter *FilterImpl) And(other common.Predicate) *AndFilter {
	return nil
}

func (filter *FilterImpl) Or(other common.Predicate) *OrFilter {
	return nil
}

type SingleFilter struct {
	predicate common.Predicate
}

func (filter *SingleFilter) Apply(ctx common.PredicateContext) bool {
	return filter.predicate.Apply(ctx)
}

func (filter *SingleFilter) String() string {
	predicateString := filter.predicate.String()
	if strings.HasPrefix(predicateString, "(") {
		return "[?" + predicateString + "]"
	} else {
		return "[?(" + predicateString + ")]"
	}
}

func NewAndFilterByPredicates(predicates []common.Predicate) *AndFilter {
	return &AndFilter{predicates: predicates}
}

func NewAndFilter(left common.Predicate, right common.Predicate) *AndFilter {
	predicates := []common.Predicate{
		left, right,
	}
	return &AndFilter{predicates: predicates}
}

type AndFilter struct {
	FilterImpl
	predicates []common.Predicate
}

func (filter *AndFilter) Apply(ctx common.PredicateContext) bool {
	for _, predicate0 := range filter.predicates {
		if !(predicate0).Apply(ctx) {
			return false
		}
	}
	return true
}

func (filter *AndFilter) String() string {
	string_ := ""
	lenPredicates := len(filter.predicates)
	for i := 0; i < lenPredicates; i++ {
		p := filter.predicates[i]
		pString := (p).String()
		if strings.HasPrefix(pString, "[?(") {
			pString = pString[3:]
		}
		string_ = string_ + pString
		if i < lenPredicates {
			string_ = string_ + "&&"
		}
	}
	return string_
}

type OrFilter struct {
}
