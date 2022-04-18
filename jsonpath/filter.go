package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"strings"
)

type Filter interface {
	common.Predicate
	Or(other common.Predicate) Filter
	And(other common.Predicate) Filter
}

type FilterImpl struct {
}

func (filter *FilterImpl) String() string {
	return ""
}

func (filter *FilterImpl) Apply(ctx common.PredicateContext) (bool, error) {
	return false, nil
}

func (filter *FilterImpl) And(other common.Predicate) Filter {
	return nil
}

func (filter *FilterImpl) Or(other common.Predicate) Filter {
	return nil
}

type SingleFilter struct {
	*FilterImpl
	predicate common.Predicate
}

func (filter *SingleFilter) Apply(ctx common.PredicateContext) (bool, error) {
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

func CreateAndFilterByPredicates(predicates []common.Predicate) *AndFilter {
	return &AndFilter{predicates: predicates}
}

func createAndFilter(left common.Predicate, right common.Predicate) *AndFilter {
	predicates := []common.Predicate{
		left, right,
	}
	return &AndFilter{predicates: predicates}
}

type AndFilter struct {
	FilterImpl
	predicates []common.Predicate
}

func (filter *AndFilter) Apply(ctx common.PredicateContext) (bool, error) {
	for _, predicate0 := range filter.predicates {
		result, err := predicate0.Apply(ctx)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func (filter *AndFilter) String() string {
	lenPredicates := len(filter.predicates)
	sb := new(strings.Builder)
	sb.WriteString("[?(")
	for i := 0; i < lenPredicates; i++ {
		p := filter.predicates[i]
		pString := (p).String()
		if strings.HasPrefix(pString, "[?(") {
			pString = pString[3 : len(pString)-2]
		}
		sb.WriteString(pString)
		if i < lenPredicates {
			sb.WriteString("&&")
		}
	}
	sb.WriteString(")]")
	return sb.String()
}

type OrFilter struct {
	*FilterImpl
	left  common.Predicate
	right common.Predicate
}

func (o *OrFilter) And(other common.Predicate) Filter {
	return createOrFilter(o.left, createAndFilter(o.right, other))
}

func (o *OrFilter) Apply(ctx common.PredicateContext) (bool, error) {
	l, err := o.left.Apply(ctx)
	if err != nil {
		return false, err
	}
	r, err := o.right.Apply(ctx)
	return l || r, err
}

func (o *OrFilter) String() string {
	sb := new(strings.Builder)
	sb.WriteString("[?(")

	l := o.left.String()
	r := o.right.String()

	if strings.HasPrefix(l, "[?(") {
		l = l[3 : len(l)-2]
	}

	if strings.HasPrefix(r, "[?(") {
		r = r[3 : len(r)-2]
	}

	sb.WriteString(l)
	sb.WriteString(" || ")
	sb.WriteString(r)
	sb.WriteString(")]")
	return sb.String()
}

func createOrFilter(left common.Predicate, right common.Predicate) *OrFilter {
	return &OrFilter{left: left, right: right}
}

func CreateSingleFilter(p common.Predicate) Filter {
	return &SingleFilter{predicate: p}
}
