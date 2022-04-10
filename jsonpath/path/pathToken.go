package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/function"
	"strconv"
)

type Token interface {
	GetTokenCount() (int, error)
	IsPathDefinite() bool
	IsUpstreamDefinite() bool
	IsTokenDefinite() bool
	String() string
	Invoke(pathFunction function.PathFunction, currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl)
	Evaluate(currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl)
	SetNext(next Token)
	GetNext() Token
	isLeaf() bool
	nextToken() (Token, error)
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
	evalPath := jsonpath.UtilsConcat(currentPath, "[", strconv.FormatInt(int64(index), 10), "]")
	var pathRef Ref
	if ctx.ForUpdate() {
		pathRef = CreatePathRef(model, index)
	} else {
		pathRef = Ref_NO_OP
	}
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
	var token Token
	token = t
	for token.isLeaf() {
		next1, err := token.nextToken()
		if err != nil {
			return -1, err
		}
		token = next1
		cnt++
	}
	return cnt, nil
}

func (t *defaultToken) String() string {
	if t.isLeaf() {
		return t.getPathFragment()
	} else {
		token, _ := t.nextToken()
		return t.getPathFragment() + token.String()
	}
}

func (t *defaultToken) Invoke(pathFunction function.PathFunction, currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) {
	ctx.AddResult(currentPath, parent, pathFunction.Invoke(currentPath, parent, model, ctx, nil))
}

func (t *defaultToken) nextToken() (Token, error) {
	if t.isLeaf() {
		return nil, &jsonpath.IllegalStateException{Message: "Current path token is a leaf"}
	}
	return t.next, nil
}

func (t *defaultToken) getPathFragment() string {
	return ""
}

func (t *defaultToken) SetNext(next Token) {
	t.next = next
}

func (t *defaultToken) GetNext() Token {
	return t.next
}

func (t *defaultToken) Evaluate(currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) {

}

type RootPathToken struct {
	*defaultToken
}

func CreateRootPathToken(token rune) *RootPathToken {
	return &RootPathToken{}
}
