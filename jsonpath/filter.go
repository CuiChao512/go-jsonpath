package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/path"
	"strings"
)

type Filter interface {
	path.Predicate
	Or(other *path.Predicate) *OrFilter
	And(other *path.Predicate) *AndFilter
}

type FilterImpl struct {
}

func (filter *FilterImpl) String() string {
	return ""
}

func (filter *FilterImpl) Apply(ctx path.PredicateContext) bool {
	return false
}

func (filter *FilterImpl) And(other path.Predicate) *AndFilter {
	return nil
}

func (filter *FilterImpl) Or(other path.Predicate) *OrFilter {
	return nil
}

type SingleFilter struct {
	predicate path.Predicate
}

func (filter *SingleFilter) Apply(ctx path.PredicateContext) bool {
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

func NewAndFilterByPredicates(predicates []*path.Predicate) *AndFilter {
	return &AndFilter{predicates: predicates}
}

func NewAndFilter(left *path.Predicate, right *path.Predicate) *AndFilter {
	predicates := []*path.Predicate{
		left, right,
	}
	return &AndFilter{predicates: predicates}
}

type AndFilter struct {
	FilterImpl
	predicates []*path.Predicate
}

func (filter *AndFilter) Apply(ctx path.PredicateContext) bool {
	for _, predicate0 := range filter.predicates {
		if !(*predicate0).Apply(ctx) {
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
