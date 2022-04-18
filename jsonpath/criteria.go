package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/filter"
	"errors"
	"regexp"
	"strings"
)

type Criteria struct {
	criteriaChain []*Criteria
	left          filter.ValueNode
	criteriaType  string
	right         filter.ValueNode
}

func (c *Criteria) Apply(ctx common.PredicateContext) (bool, error) {
	for _, expressionNode := range c.toRelationalExpressionNodes() {
		result, err := expressionNode.Apply(ctx)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func (c *Criteria) String() string {
	return common.UtilsJoin("&&", "", c.toRelationalExpressionNodes())
}

func (c *Criteria) toRelationalExpressionNodes() []*filter.RelationExpressionNode {
	nodes := make([]*filter.RelationExpressionNode, len(c.criteriaChain))
	for _, criteria := range c.criteriaChain {
		nodes = append(nodes, filter.CreateRelationExpressionNode(criteria.left, criteria.criteriaType, criteria.right))
	}
	return nodes
}

func WherePath(key common.Path) *Criteria {
	return createCriteriaByValueNode(filter.CreatePathNode(key, false, false))
}

func WhereString(key string) *Criteria {
	return createCriteriaByValueNode(filter.CreateValueNode(prefixPath(key)))
}

func (c *Criteria) And(key string) (*Criteria, error) {
	err := c.checkComplete()
	if err != nil {
		return nil, err
	}
	criteria := &Criteria{}
	valueNode := filter.CreateValueNode(key)
	criteria.left = valueNode
	criteria.criteriaChain = c.criteriaChain
	criteria.criteriaChain = append(c.criteriaChain, c)
	return criteria, nil
}

func (c *Criteria) Is(o interface{}) *Criteria {
	return c.Eq(o)
}

func (c *Criteria) Eq(o interface{}) *Criteria {
	c.criteriaType = filter.RelationalOperator_EQ
	c.right = filter.CreateValueNode(o)
	return c
}

func (c *Criteria) Ne(o interface{}) *Criteria {
	c.criteriaType = filter.RelationalOperator_NE
	c.right = filter.CreateValueNode(o)
	return c
}

func (c *Criteria) Lt(o interface{}) *Criteria {
	c.criteriaType = filter.RelationalOperator_LT
	c.right = filter.CreateValueNode(o)
	return c
}

func (c *Criteria) Lte(o interface{}) *Criteria {
	c.criteriaType = filter.RelationalOperator_LTE
	c.right = filter.CreateValueNode(o)
	return c
}

func (c *Criteria) Gt(o interface{}) *Criteria {
	c.criteriaType = filter.RelationalOperator_GT
	c.right = filter.CreateValueNode(o)
	return c
}

func (c *Criteria) Gte(o interface{}) *Criteria {
	c.criteriaType = filter.RelationalOperator_GTE
	c.right = filter.CreateValueNode(o)
	return c
}

func (c *Criteria) Regex(pattern *regexp.Regexp) (*Criteria, error) {
	if pattern == nil {
		return nil, errors.New("pattern can not be null")
	}
	c.criteriaType = filter.RelationalOperator_REGEX
	c.right = filter.CreateValueNode(pattern)
	return c, nil
}

func (c *Criteria) In(o ...interface{}) (*Criteria, error) {
	return c.InSlice(o)
}

func (c *Criteria) InSlice(l interface{}) (*Criteria, error) {
	if l == nil {
		return nil, errors.New("collection can not be null")
	}
	c.criteriaType = filter.RelationalOperator_IN
	c.right = filter.CreateValueListNode(l)
	return c, nil
}

func (c *Criteria) Contains(o interface{}) *Criteria {
	c.criteriaType = filter.RelationalOperator_CONTAINS
	c.right = filter.CreateValueNode(o)
	return c
}

func (c *Criteria) Nin(o ...interface{}) (*Criteria, error) {
	return c.NinSlice(o)
}

func (c *Criteria) NinSlice(l interface{}) (*Criteria, error) {
	if l == nil {
		return nil, errors.New("collection can not be null")
	}
	c.criteriaType = filter.RelationalOperator_NIN
	c.right = filter.CreateValueListNode(l)
	return c, nil
}

func (c *Criteria) SubSetOf(o ...interface{}) (*Criteria, error) {
	return c.SubSetOfSlice(o)
}

func (c *Criteria) SubSetOfSlice(l interface{}) (*Criteria, error) {
	if l == nil {
		return nil, errors.New("collection can not be null")
	}
	c.criteriaType = filter.RelationalOperator_SUBSETOF
	c.right = filter.CreateValueListNode(l)
	return c, nil

}

func (c *Criteria) AnyOf(o ...interface{}) (*Criteria, error) {
	return c.AnyOfSlice(o)

}

func (c *Criteria) AnyOfSlice(l interface{}) (*Criteria, error) {
	if l == nil {
		return nil, errors.New("collection can not be null")
	}
	c.criteriaType = filter.RelationalOperator_ANYOF
	c.right = filter.CreateValueListNode(l)
	return c, nil

}

func (c *Criteria) NoneOf(o ...interface{}) (*Criteria, error) {
	return c.NoneOfSlice(o)

}

func (c *Criteria) NoneOfSlice(l interface{}) (*Criteria, error) {
	if l == nil {
		return nil, errors.New("collection can not be null")
	}
	c.criteriaType = filter.RelationalOperator_NONEOF
	c.right = filter.CreateValueListNode(l)
	return c, nil

}

func (c *Criteria) All(o ...interface{}) (*Criteria, error) {
	return c.AllSlice(o)
}

func (c *Criteria) AllSlice(l interface{}) (*Criteria, error) {
	if l == nil {
		return nil, errors.New("collection can not be null")
	}
	c.criteriaType = filter.RelationalOperator_ALL
	c.right = filter.CreateValueListNode(l)
	return c, nil
}

func (c *Criteria) Size(size int) *Criteria {
	c.criteriaType = filter.RelationalOperator_SIZE
	c.right = filter.CreateValueNode(size)
	return c
}

func (c *Criteria) Type(typeString string) *Criteria {
	//TODO java class...
	c.criteriaType = filter.RelationalOperator_TYPE
	c.right = filter.CreateValueNode(typeString)
	return c
}

func (c *Criteria) Exists(shouldExist bool) (*Criteria, error) {
	c.criteriaType = filter.RelationalOperator_EXISTS
	c.right = filter.CreateValueNode(shouldExist)
	pathNode, err := c.left.AsPathNode()
	if err != nil {
		return nil, err
	}
	c.left = pathNode.AsExistsCheck(shouldExist)
	return c, nil
}

func (c *Criteria) NotEmpty() *Criteria {
	return c.Empty(false)
}

func (c *Criteria) Empty(empty bool) *Criteria {
	c.criteriaType = filter.RelationalOperator_EMPTY
	if empty {
		c.right = filter.TRUE_NODE
	} else {
		c.right = filter.FALSE_NODE
	}
	return c
}

func (c *Criteria) Matches(p common.Predicate) *Criteria {
	c.criteriaType = filter.RelationalOperator_MATCHES
	c.right = filter.CreatePredicateNode(p)
	return c
}

func (c *Criteria) checkComplete() error {
	if c.left == nil || c.criteriaType == "" || c.right == nil {
		return &common.JsonPathError{Message: "Criteria build exception. Complete on criteria before defining next."}
	}
	return nil
}

func parseCriteria(criteria string) (*Criteria, error) {
	if criteria == "" {
		return nil, &common.InvalidPathError{Message: "Criteria can not be null"}
	}
	split := strings.Split(strings.TrimSpace(criteria), " ")
	if len(split) == 3 {
		return createCriteriaByStrings(split[0], split[1], split[2]), nil
	} else if len(split) == 1 {
		return createCriteriaByStrings(split[0], filter.RelationalOperator_EXISTS, "true"), nil
	} else {
		return nil, &common.InvalidPathError{Message: "Could not parse criteria"}
	}
}

func prefixPath(key string) string {
	if !strings.HasPrefix(key, "$") && !strings.HasPrefix(key, "@") {
		key = "@." + key
	}
	return key
}

func createCriteriaByStrings(left string, operator string, right string) *Criteria {
	criteria := createCriteriaByValueNode(filter.CreateValueNode(left))
	criteria.criteriaType = operator
	criteria.right = filter.CreateValueNode(right)
	return criteria
}

func createCriteriaByValueNode(valueNode filter.ValueNode) *Criteria {
	return &Criteria{left: valueNode}
}
