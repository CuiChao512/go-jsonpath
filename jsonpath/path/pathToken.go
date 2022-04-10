package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/function"
)

type Token interface {
	GetTokenCount() (int, error)
	IsPathDefinite() bool
	IsUpstreamDefinite() bool
	IsTokenDefinite() bool
	String() string
	Invoke(pathFunction function.PathFunction, currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl)
}

type defaultToken struct {
	prev               Token
	next               Token
	definite           bool
	definiteUpdated    bool
	upstreamDefinite   bool
	upstreamUpdated    bool
	upstreamArrayIndex int
}

func (t *defaultToken) SetUpstreamArrayIndex(idx int) {
	t.upstreamArrayIndex = idx
}

func (t *defaultToken) handleObjectProperty(currentPath string, model interface{}, ctx *jsonpath.EvaluationContextImpl, properties []string) {

}

func (t *defaultToken) handleArrayIndex(index int, currentPath string, model interface{}, ctx *jsonpath.EvaluationContextImpl) {

}

func (t *defaultToken) prevToken() Token {
	return t.prev
}

func (t *defaultToken) isLeaf() bool {
	return t.next == nil
}

func (t *defaultToken) isRoot() bool {
	return t.prev == nil
}

func (t *defaultToken) IsTokenDefinite() bool {
	return false
}

func (t *defaultToken) IsPathDefinite() bool {
	if !t.definiteUpdated {
		return t.definite
	}

	isDefinite := t.IsTokenDefinite()
	if isDefinite && !t.isLeaf() {
		isDefinite = t.next.IsPathDefinite()
	}
	t.definite = isDefinite
	t.definiteUpdated = true
	return isDefinite
}

func (t *defaultToken) IsUpstreamDefinite() bool {
	if t.upstreamUpdated == false {
		t.upstreamUpdated = true
		t.upstreamDefinite = t.isRoot() || t.prev.IsPathDefinite() && t.prev.IsUpstreamDefinite()
	}
	return t.upstreamDefinite
}

func (t *defaultToken) GetTokenCount() (int, error) {
	cnt := 1
	token := t
	for token.isLeaf() {
		next, err := token.nextToken()
		if err != nil {
			return -1, err
		}
		token = next
		cnt++
	}
	return cnt, nil
}

func (t *defaultToken) String() string {
	return ""
}

func (t *defaultToken) Invoke(pathFunction function.PathFunction, currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) {

}

func (t *defaultToken) nextToken() (Token, error) {
	if t.isLeaf() {
		return nil, &jsonpath.IllegalStateException{Message: "Current path token is a leaf"}
	}
	return t.next, nil
}

func CreateRootPathToken(token rune) *RootPathToken {

}
