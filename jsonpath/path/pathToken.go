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
	Evaluate(currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) error
	SetNext(next Token)
	SetPrev(prev Token)
	GetNext() Token
	isLeaf() bool
	nextToken() (Token, error)
	prevToken() Token
	GetPathFragment() string
	appendTailToken(next Token) Token
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

func (t *defaultToken) appendTailToken(next Token) Token {
	t.next = next
	t.next.SetPrev(t)
	return next
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
			if idx == "[-1]" || ctx.GetRoot().GetTail().prevToken().GetPathFragment() == idx {
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
		return t.GetPathFragment()
	} else {
		token, _ := t.nextToken()
		return t.GetPathFragment() + token.String()
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

func (t *defaultToken) GetPathFragment() string {
	return ""
}

func (t *defaultToken) SetPrev(prev Token) {
	t.prev = prev
}

func (t *defaultToken) SetNext(next Token) {
	t.next = next
}

func (t *defaultToken) GetNext() Token {
	return t.next
}

func (t *defaultToken) Evaluate(currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) error {
	return nil
}

//RootPathToken ----
type RootPathToken struct {
	*defaultToken
	tail       Token
	tokenCount int
	rootToken  string
}

func (r *RootPathToken) GetTail() Token {
	return r.tail
}

func (r *RootPathToken) GetTokenCount() (int, error) {
	return r.tokenCount, nil
}

func (r *RootPathToken) Append(next Token) *RootPathToken {
	r.tail = r.tail.appendTailToken(next)
	r.tokenCount++
	return r
}

func (r *RootPathToken) AppendPathToken(next Token) TokenAppender {
	r.Append(next)
	return r
}

func (r *RootPathToken) GetPathTokenAppender() TokenAppender {
	return r
}

func (r *RootPathToken) Evaluate(currentPath string, ref Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) error {
	if r.isLeaf() {
		var op Ref
		if ctx.ForUpdate() {
			op = ref
		} else {
			op = PathRefNoOp
		}
		ctx.AddResult(r.rootToken, op, model)
	} else {
		next, _ := r.nextToken()
		err := next.Evaluate(r.rootToken, ref, model, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RootPathToken) GetPathFragment() string {
	return r.rootToken
}

func (r *RootPathToken) IsTokenDefinite() bool {
	return true
}

func (r *RootPathToken) IsFunctionPath() bool {
	switch jsonpath.UtilsGetPtrElem(r.tail).(type) {
	case FunctionPathToken:
		return true
	default:
		return false
	}
}

func (r *RootPathToken) SetTail(token Token) {
	r.tail = token
}

func CreateRootPathToken(token rune) *RootPathToken {
	root := &RootPathToken{}
	root.rootToken = string(token)
	root.tail = root
	root.tokenCount = 1
	return root
}

// PathTokenAppender

type TokenAppender interface {
	AppendPathToken(next Token) TokenAppender
}

// FunctionPathToken

type FunctionPathToken struct {
	*defaultToken
	functionName   string
	pathFragment   string
	functionParams []*function.Parameter
}

func (f *FunctionPathToken) Evaluate(currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) error {
	pathFunction, err := function.GetFunctionByName(f.functionName)
	if err != nil {
		return err
	}
	err = f.evaluateParameters(currentPath, parent, model, ctx)
	if err != nil {
		return err
	}
	result := pathFunction.Invoke(currentPath, parent, model, ctx, f.functionParams)
	ctx.AddResult(currentPath+"."+f.functionName, parent, result)
	f.cleanWildcardPathToken()
	if !f.isLeaf() {
		next, _ := f.nextToken()
		err = next.Evaluate(currentPath, parent, result, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FunctionPathToken) evaluateParameters(currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) error {
	if f.functionParams != nil {
		for _, param := range f.functionParams {
			switch param.GetType() {
			case function.PATH:
				pathLateBindingValue, err := function.CreateLateBindingValue(param.GetPath(), ctx.RootDocument(), ctx.Configuration())
				if err != nil {
					return err
				}
				if !param.HasEvaluated() || !pathLateBindingValue.Equals(param.GetILateBindingValue()) {
					param.SetLateBinding(pathLateBindingValue)
					param.SetEvaluated(true)
				}
			case function.JSON:
				if !param.HasEvaluated() {
					param.SetLateBinding(function.CreateJsonLateBindingValue(ctx.Configuration().JsonProvider(), param))
					param.SetEvaluated(true)
				}
			}
		}
	}
	return nil
}

func getNextTokenSuppressError(token Token) Token {
	next, _ := token.nextToken()
	return next
}

func (f *FunctionPathToken) cleanWildcardPathToken() {
	if nil != f.functionParams && len(f.functionParams) > 0 {
		path := f.functionParams[0].GetPath()
		switch jsonpath.UtilsGetPtrElem(path).(type) {
		case CompiledPath:
			if nil != path && !path.IsFunctionPath() {
				compiledPath, _ := jsonpath.UtilsGetPtrElem(path).(CompiledPath)
				root := compiledPath.GetRoot()
				tail := root.GetNext()
				for tail != nil && getNextTokenSuppressError(tail) != nil {
					switch jsonpath.UtilsGetPtrElem(tail.GetNext()).(type) {
					case WildcardPathToken:
						tail.SetNext(tail.GetNext().GetNext())
						break
					default:
						tail = tail.GetNext()
					}
				}
			}
		default:
		}
	}
}

func CreateFunctionPathToken(pathFragment string, parameters []*function.Parameter) *FunctionPathToken {
	functionPathToken := &FunctionPathToken{}
	if parameters != nil && len(parameters) > 0 {
		functionPathToken.pathFragment = pathFragment + "(...)"
	} else {
		functionPathToken.pathFragment = pathFragment + "()"
	}
	if pathFragment != "" {
		functionPathToken.functionName = pathFragment
		functionPathToken.functionParams = parameters
	} else {
		functionPathToken.functionName = pathFragment
		functionPathToken.functionParams = nil
	}
	return functionPathToken
}

//PropertyPathToken

type PropertyPathToken struct {
	*defaultToken
	properties      []string
	stringDelimiter string
}

func (p *PropertyPathToken) GetProperties() []string {
	return p.properties
}

func (p *PropertyPathToken) SinglePropertyCase() bool {
	return len(p.properties) == 1
}

func (p *PropertyPathToken) MultiPropertyMergeCase() bool {
	return p.isLeaf() && len(p.properties) > 1
}

func (p *PropertyPathToken) MultiPropertyIterationCase() bool {
	return !p.isLeaf() && len(p.properties) > 1
}

func (p *PropertyPathToken) IsTokenDefinite() bool {
	return p.SinglePropertyCase() || p.MultiPropertyMergeCase()
}

func (p *PropertyPathToken) GetPathFragment() string {
	return "[" + jsonpath.UtilsJoin(",", p.stringDelimiter, p.properties) + "]"
}

func (p *PropertyPathToken) Evaluate(currentPath string, parent Ref, model interface{}, ctx *jsonpath.EvaluationContextImpl) error {
	return nil
}

//WildCardPathToken

type WildcardPathToken struct {
	*defaultToken
}
