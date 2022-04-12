package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"cuichao.com/go-jsonpath/jsonpath/function"
	"errors"
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
	prevToken() Token
	getPathFragment() string
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

func (t *defaultToken) handleObjectProperty(currentPath string, model interface{}, ctx *jsonpath.EvaluationContextImpl, properties []string) error {

	if len(properties) == 1 {
		property := properties[0]
		evalPath := jsonpath.UtilsConcat(currentPath, "['", property, "']")
		propertyVal := pathTokenReadObjectProperty(property, model, ctx)
		if propertyVal == jsonpath.JsonProviderUndefined {
			// Conditions below heavily depend on current token type (and its logic) and are not "universal",
			// so this code is quite dangerous (I'd rather rewrite it & move to PropertyPathToken and implemented
			// WildcardPathToken as a dynamic multi prop case of PropertyPathToken).
			// Better safe than sorry.
			switch jsonpath.UtilsGetPtrElem(t).(type) {
			case PropertyPathToken:
			default:
				return errors.New("only PropertyPathToken is supported")
			}

			if t.isLeaf() {

				if jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_DEFAULT_PATH_LEAF_TO_NULL) {
					propertyVal = nil
				} else {
					if jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_SUPPRESS_EXCEPTIONS) ||
						!jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_REQUIRE_PROPERTIES) {
						return nil
					} else {
						return &jsonpath.PathNotFoundError{Message: "No results for path: " + evalPath}
					}
				}
			} else {
				if !(t.IsUpstreamDefinite() && t.IsTokenDefinite()) &&
					!jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_REQUIRE_PROPERTIES) ||
					jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_SUPPRESS_EXCEPTIONS) {
					// If there is some indefiniteness in the path and properties are not required - we'll ignore
					// absent property. And also in case of exception suppression - so that other path evaluation
					// branches could be examined.
					return nil
				} else {
					return &jsonpath.PathNotFoundError{Message: "Missing property in path " + evalPath}
				}
			}
		}

		var ref Ref

		if ctx.ForUpdate() {
			ref = CreateObjectPropertyPathRef(model, property)
		} else {
			ref = PathRefNoOp
		}
		if t.isLeaf() {
			idx := "[" + jsonpath.UtilsToString(t.upstreamArrayIndex) + "]"
			if idx == "[-1]" || ctx.GetRoot().GetTail().prevToken().getPathFragment() == idx {
				ctx.AddResult(evalPath, ref, propertyVal)
			}
		} else {
			next, _ := t.nextToken()
			next.Evaluate(evalPath, ref, propertyVal, ctx)
		}
	} else {
		evalPath := currentPath + "[" + jsonpath.UtilsJoin(", ", "'", properties) + "]"

		if !t.isLeaf() {
			return errors.New("non-leaf multi props handled elsewhere")
		}

		merged := ctx.JsonProvider().CreateMap()
		for _, property := range properties {
			var propertyVal interface{}
			if pathTokenHasProperty(property, model, ctx) {
				propertyVal = pathTokenReadObjectProperty(property, model, ctx)
				if propertyVal == jsonpath.JsonProviderUndefined {
					if jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_DEFAULT_PATH_LEAF_TO_NULL) {
						propertyVal = nil
					} else {
						continue
					}
				}
			} else {
				if jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_DEFAULT_PATH_LEAF_TO_NULL) {
					propertyVal = nil
				} else if jsonpath.UtilsSliceContains(ctx.Options(), jsonpath.OPTION_REQUIRE_PROPERTIES) {
					return &jsonpath.PathNotFoundError{Message: "Missing property in path " + evalPath}
				} else {
					continue
				}
			}
			ctx.JsonProvider().SetProperty(merged, property, propertyVal)
		}
		var pathRef Ref
		if ctx.ForUpdate() {
			pathRef = CreateObjectMultiPropertyPathRef(model, properties)
		} else {
			pathRef = PathRefNoOp
		}
		ctx.AddResult(evalPath, pathRef, merged)
	}
	return nil
}

func pathTokenHasProperty(property string, model interface{}, impl *jsonpath.EvaluationContextImpl) bool {
	return jsonpath.UtilsSliceContains(impl.JsonProvider().GetPropertyKeys(model), property)
}

func pathTokenReadObjectProperty(property string, model interface{}, ctx *jsonpath.EvaluationContextImpl) interface{} {
	return ctx.JsonProvider().GetMapValue(model, property)
}

func (t *defaultToken) handleArrayIndex(index int, currentPath string, model interface{}, ctx *jsonpath.EvaluationContextImpl) {
	evalPath := jsonpath.UtilsConcat(currentPath, "[", strconv.FormatInt(int64(index), 10), "]")
	var pathRef Ref
	if ctx.ForUpdate() {
		pathRef = CreateArrayIndexPathRef(model, index)
	} else {
		pathRef = PathRefNoOp
	}

	var effectiveIndex int
	if index < 0 {
		effectiveIndex = ctx.JsonProvider().Length(model) + index
	} else {
		effectiveIndex = index
	}

	evalHit := ctx.JsonProvider().GetArrayIndex(model, effectiveIndex)

	if t.isLeaf() {
		ctx.AddResult(evalPath, pathRef, evalHit)
	} else {
		next, err := t.nextToken()
		if err != nil {
			next.Evaluate(evalPath, pathRef, evalHit, ctx)
		}
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

//RootPathToken ----
type RootPathToken struct {
	*defaultToken
}

func (r *RootPathToken) GetTail() Token {
	//TODO:
	return nil
}

func CreateRootPathToken(token rune) *RootPathToken {
	return &RootPathToken{}
}

//PropertyPathToken

type PropertyPathToken struct {
	*defaultToken
}
