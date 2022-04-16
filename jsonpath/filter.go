package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/predicate"
	"strings"
)

type Filter interface {
	predicate.Predicate
	Or(other *predicate.Predicate) *OrFilter
	And(other *predicate.Predicate) *AndFilter
}

type FilterImpl struct {
}

func (filter *FilterImpl) String() string {
	return ""
}

func (filter *FilterImpl) Apply(ctx predicate.PredicateContext) bool {
	return false
}

func (filter *FilterImpl) And(other predicate.Predicate) *AndFilter {
	return nil
}

func (filter *FilterImpl) Or(other predicate.Predicate) *OrFilter {
	return nil
}

type SingleFilter struct {
	predicate predicate.Predicate
}

func (filter *SingleFilter) Apply(ctx predicate.PredicateContext) bool {
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

func NewAndFilterByPredicates(predicates []*predicate.Predicate) *AndFilter {
	return &AndFilter{predicates: predicates}
}

func NewAndFilter(left *predicate.Predicate, right *predicate.Predicate) *AndFilter {
	predicates := []*predicate.Predicate{
		left, right,
	}
	return &AndFilter{predicates: predicates}
}

type AndFilter struct {
	FilterImpl
	predicates []*predicate.Predicate
}

func (filter *AndFilter) Apply(ctx predicate.PredicateContext) bool {
	for _, predicate := range filter.predicates {
		if !(*predicate).Apply(ctx) {
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
		pString := (*p).String()
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
