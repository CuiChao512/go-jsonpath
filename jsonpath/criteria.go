package jsonpath

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/filter"
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

func (c *Criteria) WherePath(key common.Path) *Criteria {
	return createCriteria(filter.CreatePathNode(key, false, false))
}

func prefixPath(key string) string {
	if !strings.HasPrefix(key, "$") && !strings.HasPrefix(key, "@") {
		key = "@." + key
	}
	return key
}

func (c *Criteria) WhereString(key string) *Criteria {
	return createCriteria()
}

func (c *Criteria) checkComplete() error {
	if c.left == nil || c.criteriaType == "" || c.right == nil {
		return &common.JsonPathError{Message: "Criteria build exception. Complete on criteria before defining next."}
	}
	return nil
}

func createCriteria(pathNode *filter.PathNode) *Criteria {
	return &Criteria{left: pathNode}
}
