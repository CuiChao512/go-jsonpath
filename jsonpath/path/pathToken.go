package path

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/function"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type Token interface {
	GetTokenCount() (int, error)
	IsPathDefinite() bool
	IsUpstreamDefinite() bool
	IsTokenDefinite() bool
	String() string
	Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error
	Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error
	SetNext(next Token)
	SetPrev(prev Token)
	GetNext() Token
	isLeaf() bool
	nextToken() (Token, error)
	prevToken() Token
	GetPathFragment() string
	appendTailToken(next Token) Token
	SetUpstreamArrayIndex(idx int)
	getUpstreamArrayIndex() int
	setDefinite(definite bool)
	isDefinite() bool
	setDefiniteUpdated(definiteUpdated bool)
	isDefiniteUpdated() bool
	setUpstreamUpdated(upstreamUpdated bool)
	isUpstreamUpdated() bool
	setUpstreamDefinite(upstreamDefinite bool)
	isRoot() bool
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

func tokenAppendTailToken(dt Token, next Token) Token {
	dt.SetNext(next)
	dt.GetNext().SetPrev(dt)
	return next
}

func tokenSetUpstreamArrayIndex(dt Token, idx int) {
	dt.SetUpstreamArrayIndex(idx)
}

func tokenHandleObjectProperty(dt Token, currentPath string, model interface{}, ctx *EvaluationContextImpl, properties []string) error {

	if len(properties) == 1 {
		property := properties[0]
		evalPath := common.UtilsConcat(currentPath, "['", property, "']")
		propertyVal := pathTokenReadObjectProperty(property, model, ctx)
		fmt.Println("propertyVal:", common.UtilsToString(propertyVal))
		if propertyVal == common.JsonProviderUndefined {
			// Conditions below heavily depend on current token type (and its logic) and are not "universal",
			// so this code is quite dangerous (I'd rather rewrite it & move to PropertyPathToken and implemented
			// WildcardPathToken as a dynamic multi prop case of PropertyPathToken).
			// Better safe than sorry.
			switch common.UtilsGetPtrElem(dt).(type) {
			case PropertyPathToken:
			default:
				return errors.New("only PropertyPathToken is supported")
			}

			if dt.isLeaf() {

				if common.UtilsSliceContains(ctx.Options(), common.OPTION_DEFAULT_PATH_LEAF_TO_NULL) {
					propertyVal = nil
				} else {
					if common.UtilsSliceContains(ctx.Options(), common.OPTION_SUPPRESS_EXCEPTIONS) ||
						!common.UtilsSliceContains(ctx.Options(), common.OPTION_REQUIRE_PROPERTIES) {
						return nil
					} else {
						return &common.PathNotFoundError{Message: "No results for path: " + evalPath}
					}
				}
			} else {
				if !(dt.IsUpstreamDefinite() && dt.IsTokenDefinite()) &&
					!common.UtilsSliceContains(ctx.Options(), common.OPTION_REQUIRE_PROPERTIES) ||
					common.UtilsSliceContains(ctx.Options(), common.OPTION_SUPPRESS_EXCEPTIONS) {
					// If there is some indefiniteness in the path and properties are not required - we'll ignore
					// absent property. And also in case of exception suppression - so that other path evaluation
					// branches could be examined.
					return nil
				} else {
					return &common.PathNotFoundError{Message: "Missing property in path " + evalPath}
				}
			}
		}

		var ref common.PathRef

		if ctx.ForUpdate() {
			ref = CreateObjectPropertyPathRef(model, property)
		} else {
			ref = PathRefNoOp
		}
		if dt.isLeaf() {
			idx := "[" + common.UtilsToString(dt.getUpstreamArrayIndex()) + "]"
			root, err := ctx.GetRoot()
			if err != nil {
				return err
			}

			if idx == "[-1]" || root.GetTail().prevToken().GetPathFragment() == idx {
				if err = ctx.AddResult(evalPath, ref, propertyVal); err != nil {
					return err
				}
			}
		} else {
			next, _ := dt.nextToken()
			err := next.Evaluate(evalPath, ref, propertyVal, ctx)
			if err != nil {
				return err
			}
		}
	} else {
		evalPath := currentPath + "[" + common.UtilsJoin(", ", "'", properties) + "]"

		if !dt.isLeaf() {
			return errors.New("non-leaf multi props handled elsewhere")
		}

		merged := ctx.JsonProvider().CreateMap()
		for _, property := range properties {
			var propertyVal interface{}
			pathTokenHasProperty, err := pathTokenHasProperty(property, model, ctx)
			if err != nil {
				return err
			}
			if pathTokenHasProperty {
				propertyVal = pathTokenReadObjectProperty(property, model, ctx)
				if propertyVal == common.JsonProviderUndefined {
					if common.UtilsSliceContains(ctx.Options(), common.OPTION_DEFAULT_PATH_LEAF_TO_NULL) {
						propertyVal = nil
					} else {
						continue
					}
				}
			} else {
				if common.UtilsSliceContains(ctx.Options(), common.OPTION_DEFAULT_PATH_LEAF_TO_NULL) {
					propertyVal = nil
				} else if common.UtilsSliceContains(ctx.Options(), common.OPTION_REQUIRE_PROPERTIES) {
					return &common.PathNotFoundError{Message: "Missing property in path " + evalPath}
				} else {
					continue
				}
			}
			var m interface{} = merged
			if err = ctx.JsonProvider().SetProperty(&m, property, propertyVal); err != nil {
				return err
			}
		}
		var pathRef common.PathRef
		if ctx.ForUpdate() {
			pathRef = CreateObjectMultiPropertyPathRef(model, properties)
		} else {
			pathRef = PathRefNoOp
		}
		return ctx.AddResult(evalPath, pathRef, merged)
	}
	return nil
}

func pathTokenHasProperty(property string, model interface{}, impl *EvaluationContextImpl) (bool, error) {
	propertyKeys, err := impl.JsonProvider().GetPropertyKeys(model)
	if err != nil {
		return false, err
	}
	return common.UtilsSliceContains(propertyKeys, property), nil
}

func pathTokenReadObjectProperty(property string, model interface{}, ctx *EvaluationContextImpl) interface{} {
	return ctx.JsonProvider().GetMapValue(model, property)
}

func tokenHandleArrayIndex(dt Token, index int, currentPath string, model interface{}, ctx *EvaluationContextImpl) error {
	evalPath := common.UtilsConcat(currentPath, "[", strconv.FormatInt(int64(index), 10), "]")
	var pathRef common.PathRef
	if ctx.ForUpdate() {
		pathRef = CreateArrayIndexPathRef(model, index)
	} else {
		pathRef = PathRefNoOp
	}

	var effectiveIndex int
	if index < 0 {
		length, err := ctx.JsonProvider().Length(model)
		if err != nil {
			return err
		}
		effectiveIndex = length + index
	} else {
		effectiveIndex = index
	}

	evalHit := ctx.JsonProvider().GetArrayIndex(model, effectiveIndex)

	if dt.isLeaf() {
		if err := ctx.AddResult(evalPath, pathRef, evalHit); err != nil {
			return err
		}
	} else {
		next, err := dt.nextToken()
		if err != nil {
			return err
		}
		return next.Evaluate(evalPath, pathRef, evalHit, ctx)
	}
	return nil
}

func tokenIsPathDefinite(dt Token) bool {
	if dt.isDefiniteUpdated() {
		return dt.isDefinite()
	}

	isDefinite := dt.IsTokenDefinite()
	//isDefinite := true
	if isDefinite && !dt.isLeaf() {
		isDefinite = dt.GetNext().IsPathDefinite()
	}
	dt.setDefinite(isDefinite)
	dt.setDefiniteUpdated(true)
	return isDefinite
}

func tokenIsUpstreamDefinite(dt Token) bool {
	if dt.isUpstreamUpdated() == false {
		dt.setUpstreamUpdated(true)
		dt.setUpstreamDefinite(dt.isRoot() || dt.prevToken().IsPathDefinite() && dt.prevToken().IsUpstreamDefinite())
	}
	return dt.IsUpstreamDefinite()
}

func tokenGetTokenCount(dt Token) (int, error) {
	cnt := 1
	var token Token
	token = dt
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

func tokenString(dt Token) string {
	if dt.isLeaf() {
		return dt.GetPathFragment()
	} else {
		token, _ := dt.nextToken()
		return dt.GetPathFragment() + token.String()
	}
}

func tokenInvoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	result, err := pathFunction.Invoke(currentPath, parent, model, ctx, nil)
	if err != nil {
		return err
	}
	return ctx.AddResult(currentPath, parent, result)
}

func tokenNextToken(dt Token) (Token, error) {
	if dt.isLeaf() {
		return nil, &common.IllegalStateException{Message: "Current path token is a leaf"}
	}
	return dt.GetNext(), nil
}

//RootPathToken ----
type RootPathToken struct {
	*defaultToken
	tail       Token
	tokenCount int
	rootToken  string
}

func (r *RootPathToken) SetPrev(prev Token) {
	r.prev = prev
}

func (r *RootPathToken) SetNext(next Token) {
	r.next = next
}

func (r *RootPathToken) GetNext() Token {
	return r.next
}

func (r *RootPathToken) GetTail() Token {
	return r.tail
}

func (r *RootPathToken) isLeaf() bool {
	return r.next == nil
}

func (r *RootPathToken) isRoot() bool {
	return r.prev == nil
}

func (r *RootPathToken) nextToken() (Token, error) {
	return tokenNextToken(r)
}

func (r *RootPathToken) prevToken() Token {
	return r.prev
}

func (r *RootPathToken) setDefiniteUpdated(definiteUpdated bool) {
	r.definiteUpdated = definiteUpdated
}

func (r *RootPathToken) isDefiniteUpdated() bool {
	return r.definiteUpdated
}

func (r *RootPathToken) setDefinite(definite bool) {
	r.definite = definite
}

func (r *RootPathToken) isDefinite() bool {
	return r.definite
}

func (r *RootPathToken) setUpstreamUpdated(upstreamUpdated bool) {
	r.upstreamUpdated = upstreamUpdated
}

func (r *RootPathToken) isUpstreamUpdated() bool {
	return r.upstreamUpdated
}

func (r *RootPathToken) setUpstreamDefinite(upstreamDefinite bool) {
	r.upstreamDefinite = upstreamDefinite
}

func (r *RootPathToken) IsUpstreamDefinite() bool {
	return r.upstreamDefinite
}

func (r *RootPathToken) String() string {
	return tokenString(r)
}

func (r *RootPathToken) getUpstreamArrayIndex() int {
	return r.upstreamArrayIndex
}

func (r *RootPathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (r *RootPathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(r)
}

func (r *RootPathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(r, idx)
}

func (r *RootPathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(r, next)
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

func (r *RootPathToken) Evaluate(currentPath string, ref common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	if r.isLeaf() {
		var op common.PathRef
		if ctx.ForUpdate() {
			op = ref
		} else {
			op = PathRefNoOp
		}
		return ctx.AddResult(r.rootToken, op, model)
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
	switch common.UtilsGetPtrElem(r.tail).(type) {
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
	root.defaultToken = &defaultToken{upstreamArrayIndex: -1}
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

func (f *FunctionPathToken) SetPrev(prev Token) {
	f.prev = prev
}

func (f *FunctionPathToken) SetNext(next Token) {
	f.next = next
}

func (f *FunctionPathToken) GetNext() Token {
	return f.next
}

func (f *FunctionPathToken) isLeaf() bool {
	return f.next == nil
}

func (f *FunctionPathToken) isRoot() bool {
	return f.prev == nil
}

func (f *FunctionPathToken) nextToken() (Token, error) {
	return tokenNextToken(f)
}

func (f *FunctionPathToken) prevToken() Token {
	return f.prev
}

func (f *FunctionPathToken) setDefiniteUpdated(definiteUpdated bool) {
	f.definiteUpdated = definiteUpdated
}

func (f *FunctionPathToken) isDefiniteUpdated() bool {
	return f.definiteUpdated
}

func (f *FunctionPathToken) setDefinite(definite bool) {
	f.definite = definite
}

func (f *FunctionPathToken) isDefinite() bool {
	return f.definite
}

func (f *FunctionPathToken) setUpstreamUpdated(upstreamUpdated bool) {
	f.upstreamUpdated = upstreamUpdated
}

func (f *FunctionPathToken) isUpstreamUpdated() bool {
	return f.upstreamUpdated
}

func (f *FunctionPathToken) setUpstreamDefinite(upstreamDefinite bool) {
	f.upstreamDefinite = upstreamDefinite
}

func (f *FunctionPathToken) IsUpstreamDefinite() bool {
	return f.upstreamDefinite
}

func (f *FunctionPathToken) String() string {
	return tokenString(f)
}

func (f *FunctionPathToken) getUpstreamArrayIndex() int {
	return f.upstreamArrayIndex
}

func (f *FunctionPathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (f *FunctionPathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(f)
}

func (f *FunctionPathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(f, idx)
}

func (f *FunctionPathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(f, next)
}

func (f *FunctionPathToken) GetTokenCount() (int, error) {
	return tokenGetTokenCount(f)
}

func (f *FunctionPathToken) IsTokenDefinite() bool {
	return true
}

func (f *FunctionPathToken) GetPathFragment() string {
	return "." + f.pathFragment
}

func (f *FunctionPathToken) Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	pathFunction, err := function.GetFunctionByName(f.functionName)
	if err != nil {
		return err
	}
	err = f.evaluateParameters(currentPath, parent, model, ctx)
	if err != nil {
		return err
	}
	result, err := pathFunction.Invoke(currentPath, parent, model, ctx, &f.functionParams)
	if err != nil {
		return err
	}
	if err = ctx.AddResult(currentPath+"."+f.functionName, parent, result); err != nil {
		return err
	}
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

func (f *FunctionPathToken) evaluateParameters(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
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
		switch common.UtilsGetPtrElem(path).(type) {
		case CompiledPath:
			if nil != path && !path.IsFunctionPath() {
				compiledPath, _ := common.UtilsGetPtrElem(path).(CompiledPath)
				root := compiledPath.GetRoot()
				tail := root.GetNext()
				for tail != nil && getNextTokenSuppressError(tail) != nil {
					switch common.UtilsGetPtrElem(tail.GetNext()).(type) {
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
	functionPathToken.defaultToken = &defaultToken{upstreamArrayIndex: -1}

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

func (p *PropertyPathToken) SetPrev(prev Token) {
	p.prev = prev
}

func (p *PropertyPathToken) SetNext(next Token) {
	p.next = next
}

func (p *PropertyPathToken) GetNext() Token {
	return p.next
}

func (p *PropertyPathToken) isLeaf() bool {
	return p.next == nil
}

func (p *PropertyPathToken) isRoot() bool {
	return p.prev == nil
}

func (p *PropertyPathToken) nextToken() (Token, error) {
	return tokenNextToken(p)
}

func (p *PropertyPathToken) prevToken() Token {
	return p.prev
}

func (p *PropertyPathToken) setDefiniteUpdated(definiteUpdated bool) {
	p.definiteUpdated = definiteUpdated
}

func (p *PropertyPathToken) isDefiniteUpdated() bool {
	return p.definiteUpdated
}

func (p *PropertyPathToken) setDefinite(definite bool) {
	p.definite = definite
}

func (p *PropertyPathToken) isDefinite() bool {
	return p.definite
}

func (p *PropertyPathToken) setUpstreamUpdated(upstreamUpdated bool) {
	p.upstreamUpdated = upstreamUpdated
}

func (p *PropertyPathToken) isUpstreamUpdated() bool {
	return p.upstreamUpdated
}

func (p *PropertyPathToken) setUpstreamDefinite(upstreamDefinite bool) {
	p.upstreamDefinite = upstreamDefinite
}

func (p *PropertyPathToken) IsUpstreamDefinite() bool {
	return p.upstreamDefinite
}

func (p *PropertyPathToken) String() string {
	return tokenString(p)
}

func (p *PropertyPathToken) getUpstreamArrayIndex() int {
	return p.upstreamArrayIndex
}

func (p *PropertyPathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (p *PropertyPathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(p)
}

func (p *PropertyPathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(p, idx)
}

func (p *PropertyPathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(p, next)
}

func (p *PropertyPathToken) GetTokenCount() (int, error) {
	return tokenGetTokenCount(p)
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
	return "[" + common.UtilsJoin(",", p.stringDelimiter, p.properties) + "]"
}

func (p *PropertyPathToken) Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	var truthCount int = 0
	if p.SinglePropertyCase() {
		truthCount++
	}
	if p.MultiPropertyIterationCase() {
		truthCount++
	}
	if p.MultiPropertyMergeCase() {
		truthCount++
	}
	if truthCount != 1 {
		return errors.New("")
	}

	if !ctx.JsonProvider().IsMap(model) {
		if !p.IsUpstreamDefinite() || common.UtilsSliceContains(ctx.Options(), common.OPTION_SUPPRESS_EXCEPTIONS) {
			return nil
		} else {
			var m string
			if model == nil {
				m = "null"
			} else {
				m = reflect.TypeOf(common.UtilsGetPtrElem(model)).Name()
			}
			message := fmt.Sprintf("Expected to find an object with property %s in path %s but found '%s'. "+
				"This is not a json object according to the JsonProvider: '%s'.",
				p.GetPathFragment(), currentPath, m, reflect.TypeOf(common.UtilsGetPtrElem(ctx.Configuration().JsonProvider())).Name())
			return &common.PathNotFoundError{Message: message}
		}
	}

	if p.SinglePropertyCase() || p.MultiPropertyMergeCase() {
		return tokenHandleObjectProperty(p, currentPath, model, ctx, p.properties)
	}

	if !p.MultiPropertyIterationCase() {
		return errors.New("")
	}

	for _, property := range p.properties {
		err := tokenHandleObjectProperty(p, currentPath, model, ctx, []string{property})
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePropertyPathToken(properties []string, stringDelimiter string) *PropertyPathToken {
	return &PropertyPathToken{defaultToken: &defaultToken{upstreamArrayIndex: -1}, properties: properties, stringDelimiter: stringDelimiter}
}

//WildCardPathToken

type WildcardPathToken struct {
	*defaultToken
}

func (w *WildcardPathToken) SetPrev(prev Token) {
	w.prev = prev
}

func (w *WildcardPathToken) SetNext(next Token) {
	w.next = next
}

func (w *WildcardPathToken) GetNext() Token {
	return w.next
}

func (w *WildcardPathToken) isLeaf() bool {
	return w.next == nil
}

func (w *WildcardPathToken) isRoot() bool {
	return w.prev == nil
}

func (w *WildcardPathToken) nextToken() (Token, error) {
	return tokenNextToken(w)
}

func (w *WildcardPathToken) prevToken() Token {
	return w.prev
}

func (w *WildcardPathToken) setDefiniteUpdated(definiteUpdated bool) {
	w.definiteUpdated = definiteUpdated
}

func (w *WildcardPathToken) isDefiniteUpdated() bool {
	return w.definiteUpdated
}

func (w *WildcardPathToken) setDefinite(definite bool) {
	w.definite = definite
}

func (w *WildcardPathToken) isDefinite() bool {
	return w.definite
}

func (w *WildcardPathToken) setUpstreamUpdated(upstreamUpdated bool) {
	w.upstreamUpdated = upstreamUpdated
}

func (w *WildcardPathToken) isUpstreamUpdated() bool {
	return w.upstreamUpdated
}

func (w *WildcardPathToken) setUpstreamDefinite(upstreamDefinite bool) {
	w.upstreamDefinite = upstreamDefinite
}

func (w *WildcardPathToken) IsUpstreamDefinite() bool {
	return w.upstreamDefinite
}

func (w *WildcardPathToken) String() string {
	return tokenString(w)
}

func (w *WildcardPathToken) getUpstreamArrayIndex() int {
	return w.upstreamArrayIndex
}

func (w *WildcardPathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (w *WildcardPathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(w)
}

func (w *WildcardPathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(w, idx)
}

func (w *WildcardPathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(w, next)
}

func (w *WildcardPathToken) GetTokenCount() (int, error) {
	return tokenGetTokenCount(w)
}

func (w *WildcardPathToken) IsTokenDefinite() bool {
	return false
}

func (w *WildcardPathToken) GetPathFragment() string {
	return "[*]"
}

func (w *WildcardPathToken) Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	if ctx.JsonProvider().IsMap(model) {
		propertyKeys, err := ctx.JsonProvider().GetPropertyKeys(model)
		if err != nil {
			return err
		}
		for _, property := range propertyKeys {
			err := tokenHandleObjectProperty(w, currentPath, model, ctx, []string{property})
			if err != nil {
				return err
			}
		}
	} else if ctx.JsonProvider().IsArray(model) {
		length, err := ctx.JsonProvider().Length(model)
		if err != nil {
			return err
		}

		for idx := 0; idx < length; idx++ {
			err := tokenHandleArrayIndex(w, idx, currentPath, model, ctx)

			if err != nil && common.UtilsSliceContains(ctx.Options(), common.OPTION_REQUIRE_PROPERTIES) {
				return err
			}
		}
	}
	return nil
}

func CreateWildcardPathToken() *WildcardPathToken {
	return &WildcardPathToken{defaultToken: &defaultToken{upstreamArrayIndex: -1}}
}

// ScanPathToken -----

type defaultScanPredicate struct {
}

func (*defaultScanPredicate) matches(model interface{}) (bool, error) {
	return false, nil
}

type ScanPredicate interface {
	matches(model interface{}) (bool, error)
}

type filterPathTokenPredicate struct {
	ctx                *EvaluationContextImpl
	predicatePathToken *PredicatePathToken
}

func (f *filterPathTokenPredicate) matches(model interface{}) (bool, error) {
	return f.predicatePathToken.accept(model, f.ctx.RootDocument(), f.ctx.Configuration(), f.ctx)
}

func createFilterPathTokenPredicate(target Token, ctx *EvaluationContextImpl) *filterPathTokenPredicate {
	f := &filterPathTokenPredicate{}
	t, _ := target.(*PredicatePathToken)
	f.predicatePathToken = t
	f.ctx = ctx
	return f
}

type wildCardPathTokenPredicate struct {
}

func (*wildCardPathTokenPredicate) matches(model interface{}) (bool, error) {
	return true, nil
}

type arrayPathTokenPredicate struct {
	ctx *EvaluationContextImpl
}

func (a *arrayPathTokenPredicate) matches(model interface{}) (bool, error) {
	return a.ctx.JsonProvider().IsArray(model), nil
}

type propertyPathTokenPredicate struct {
	ctx               *EvaluationContextImpl
	propertyPathToken *PropertyPathToken
}

func (p *propertyPathTokenPredicate) matches(model interface{}) (bool, error) {
	if !p.ctx.JsonProvider().IsMap(model) {
		return false, nil
	}

	if !p.propertyPathToken.IsTokenDefinite() {
		return true, nil
	}

	if p.propertyPathToken.isLeaf() && common.UtilsSliceContains(p.ctx.Options(), common.OPTION_DEFAULT_PATH_LEAF_TO_NULL) {
		return true, nil
	}
	propertyKeys, err := p.ctx.JsonProvider().GetPropertyKeys(model)
	if err != nil {
		return false, err
	}
	return common.UtilsStringSliceContainsAll(propertyKeys, p.propertyPathToken.GetProperties()), nil
}

func createPropertyPathTokenPredicate(target *PropertyPathToken, ctx *EvaluationContextImpl) *propertyPathTokenPredicate {
	return &propertyPathTokenPredicate{propertyPathToken: target, ctx: ctx}
}

var falseScanPredicate = &defaultScanPredicate{}

type ScanPathToken struct {
	*defaultToken
}

func (s *ScanPathToken) SetPrev(prev Token) {
	s.prev = prev
}

func (s *ScanPathToken) SetNext(next Token) {
	s.next = next
}

func (s *ScanPathToken) GetNext() Token {
	return s.next
}

func (s *ScanPathToken) isLeaf() bool {
	return s.next == nil
}

func (s *ScanPathToken) isRoot() bool {
	return s.prev == nil
}

func (s *ScanPathToken) nextToken() (Token, error) {
	return tokenNextToken(s)
}

func (s *ScanPathToken) prevToken() Token {
	return s.prev
}

func (s *ScanPathToken) setDefiniteUpdated(definiteUpdated bool) {
	s.definiteUpdated = definiteUpdated
}

func (s *ScanPathToken) isDefiniteUpdated() bool {
	return s.definiteUpdated
}

func (s *ScanPathToken) setDefinite(definite bool) {
	s.definite = definite
}

func (s *ScanPathToken) isDefinite() bool {
	return s.definite
}

func (s *ScanPathToken) setUpstreamUpdated(upstreamUpdated bool) {
	s.upstreamUpdated = upstreamUpdated
}

func (s *ScanPathToken) isUpstreamUpdated() bool {
	return s.upstreamUpdated
}

func (s *ScanPathToken) setUpstreamDefinite(upstreamDefinite bool) {
	s.upstreamDefinite = upstreamDefinite
}

func (s *ScanPathToken) IsUpstreamDefinite() bool {
	return s.upstreamDefinite
}

func (s *ScanPathToken) String() string {
	return tokenString(s)
}

func (s *ScanPathToken) getUpstreamArrayIndex() int {
	return s.upstreamArrayIndex
}

func (s *ScanPathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (s *ScanPathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(s)
}

func (s *ScanPathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(s, idx)
}

func (s *ScanPathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(s, next)
}

func (s *ScanPathToken) GetTokenCount() (int, error) {
	return tokenGetTokenCount(s)
}

func (s *ScanPathToken) Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	pt, err := s.nextToken()
	if err != nil {
		return err
	}
	return s.walk(pt, currentPath, parent, model, ctx, s.createScanPredicate(pt, ctx))
}

func (s *ScanPathToken) walk(pt Token, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl, predicate ScanPredicate) error {
	if ctx.JsonProvider().IsMap(model) {
		return s.walkObject(pt, currentPath, parent, model, ctx, predicate)
	} else if ctx.JsonProvider().IsArray(model) {
		return s.walkArray(pt, currentPath, parent, model, ctx, predicate)
	}
	return nil
}

func (s *ScanPathToken) walkObject(pt Token, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl, predicate ScanPredicate) error {
	matchesResult, err := predicate.matches(model)
	if err != nil {
		return err
	}
	if matchesResult {
		err := pt.Evaluate(currentPath, parent, model, ctx)
		if err != nil {
			return err
		}
	}
	properties, err := ctx.JsonProvider().GetPropertyKeys(model)
	if err != nil {
		return err
	}
	for _, property := range properties {
		evalPath := currentPath + "['" + property + "']"
		propertyModel := ctx.JsonProvider().GetMapValue(model, property)
		if propertyModel != common.JsonProviderUndefined {
			err := s.walk(pt, evalPath, CreateObjectPropertyPathRef(model, property), propertyModel, ctx, predicate)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ScanPathToken) walkArray(pt Token, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl, predicate ScanPredicate) error {
	matchesResult, err := predicate.matches(model)
	if err != nil {
		return err
	}
	if matchesResult {
		if pt.isLeaf() {
			err = pt.Evaluate(currentPath, parent, model, ctx)
			if err != nil {
				return err
			}
		} else {
			next, err := pt.nextToken()
			if err != nil {
				return err
			}
			models, err := ctx.JsonProvider().ToArray(model)
			if err != nil {
				return err
			}
			idx := 0
			for _, evalModel := range models {
				evalPath := currentPath + "[" + strconv.Itoa(idx) + "]"
				next.SetUpstreamArrayIndex(idx)
				err = next.Evaluate(evalPath, parent, evalModel, ctx)
				if err != nil {
					return err
				}
				idx++
			}
		}
	}

	models, err := ctx.JsonProvider().ToArray(model)
	if err != nil {
		return err
	}
	idx := 0
	for _, evalModel := range models {
		evalPath := currentPath + "[" + strconv.Itoa(idx) + "]"
		err := s.walk(pt, evalPath, CreateArrayIndexPathRef(model, idx), evalModel, ctx, predicate)
		if err != nil {
			return err
		}
		idx++
	}
	return nil
}

func (*ScanPathToken) createScanPredicate(target Token, ctx *EvaluationContextImpl) ScanPredicate {
	switch target.(type) {
	case *PropertyPathToken:
		p, _ := target.(*PropertyPathToken)
		return createPropertyPathTokenPredicate(p, ctx)
	case *ArrayIndexPathToken:
		return &arrayPathTokenPredicate{ctx: ctx}
	case *ArraySlicePathToken:
		return &arrayPathTokenPredicate{ctx: ctx}
	case *WildcardPathToken:
		return &wildCardPathTokenPredicate{}
	case *PredicatePathToken:
		return createFilterPathTokenPredicate(target, ctx)
	default:
		return falseScanPredicate
	}
}

func (*ScanPathToken) IsTokenDefinite() bool {
	return false
}

func (*ScanPathToken) GetPathFragment() string {
	return ".."
}

func CreateScanPathToken() *ScanPathToken {
	return &ScanPathToken{}
}

type ArrayIndexPathToken struct {
	*defaultToken
	arrayIndexOperation *ArrayIndexOperation
}

func (a *ArrayIndexPathToken) SetPrev(prev Token) {
	a.prev = prev
}

func (a *ArrayIndexPathToken) SetNext(next Token) {
	a.next = next
}

func (a *ArrayIndexPathToken) GetNext() Token {
	return a.next
}

func (a *ArrayIndexPathToken) isLeaf() bool {
	return a.next == nil
}

func (a *ArrayIndexPathToken) isRoot() bool {
	return a.prev == nil
}

func (a *ArrayIndexPathToken) nextToken() (Token, error) {
	return tokenNextToken(a)
}

func (a *ArrayIndexPathToken) prevToken() Token {
	return a.prev
}

func (a *ArrayIndexPathToken) setDefiniteUpdated(definiteUpdated bool) {
	a.definiteUpdated = definiteUpdated
}

func (a *ArrayIndexPathToken) isDefiniteUpdated() bool {
	return a.definiteUpdated
}

func (a *ArrayIndexPathToken) setDefinite(definite bool) {
	a.definite = definite
}

func (a *ArrayIndexPathToken) isDefinite() bool {
	return a.definite
}

func (a *ArrayIndexPathToken) setUpstreamUpdated(upstreamUpdated bool) {
	a.upstreamUpdated = upstreamUpdated
}

func (a *ArrayIndexPathToken) isUpstreamUpdated() bool {
	return a.upstreamUpdated
}

func (a *ArrayIndexPathToken) setUpstreamDefinite(upstreamDefinite bool) {
	a.upstreamDefinite = upstreamDefinite
}

func (a *ArrayIndexPathToken) IsUpstreamDefinite() bool {
	return a.upstreamDefinite
}

func (a *ArrayIndexPathToken) String() string {
	return tokenString(a)
}

func (a *ArrayIndexPathToken) getUpstreamArrayIndex() int {
	return a.upstreamArrayIndex
}

func (a *ArrayIndexPathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (a *ArrayIndexPathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(a)
}

func (a *ArrayIndexPathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(a, idx)
}

func (a *ArrayIndexPathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(a, next)
}

func (a *ArrayIndexPathToken) GetTokenCount() (int, error) {
	return tokenGetTokenCount(a)
}

func (a *ArrayIndexPathToken) checkArrayModel(currentPath string, model interface{}, ctx *EvaluationContextImpl) (bool, error) {
	if model == nil {
		if !a.IsTokenDefinite() || common.UtilsSliceContains(ctx.Options(), common.OPTION_SUPPRESS_EXCEPTIONS) {
			return false, nil
		} else {
			return false, &common.PathNotFoundError{Message: "The path " + currentPath + " is null"}
		}
	}

	if ctx.JsonProvider().IsArray(model) {
		if a.IsUpstreamDefinite() || common.UtilsSliceContains(ctx.Options(), common.OPTION_SUPPRESS_EXCEPTIONS) {
			return false, nil
		} else {
			return false, &common.PathNotFoundError{Message: fmt.Sprintf("Filter: %s can only be applied to arrays. Current context is: %s", a, model)}
		}
	}
	return true, nil
}

func (a *ArrayIndexPathToken) Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	checkResult, err := a.checkArrayModel(currentPath, model, ctx)
	if err != nil {
		return err
	}

	if !checkResult {
		return nil
	}

	if a.arrayIndexOperation.IsSingleIndexOperation() {
		return tokenHandleArrayIndex(a, a.arrayIndexOperation.Indexes()[0], currentPath, model, ctx)
	} else {
		for _, idx := range a.arrayIndexOperation.Indexes() {
			err = tokenHandleArrayIndex(a, idx, currentPath, model, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *ArrayIndexPathToken) GetPathFragment() string {
	return a.arrayIndexOperation.String()
}

func (a *ArrayIndexPathToken) IsTokenDefinite() bool {
	return a.arrayIndexOperation.IsSingleIndexOperation()
}

func CreateArrayIndexPathToken(arrayIndexOperation *ArrayIndexOperation) *ArrayIndexPathToken {
	return &ArrayIndexPathToken{defaultToken: &defaultToken{upstreamArrayIndex: -1}, arrayIndexOperation: arrayIndexOperation}
}

// ArraySlicePathToken -----
type ArraySlicePathToken struct {
	*defaultToken
	operation *ArraySliceOperation
}

func (a *ArraySlicePathToken) SetPrev(prev Token) {
	a.prev = prev
}

func (a *ArraySlicePathToken) SetNext(next Token) {
	a.next = next
}

func (a *ArraySlicePathToken) GetNext() Token {
	return a.next
}

func (a *ArraySlicePathToken) isLeaf() bool {
	return a.next == nil
}

func (a *ArraySlicePathToken) isRoot() bool {
	return a.prev == nil
}

func (a *ArraySlicePathToken) nextToken() (Token, error) {
	return tokenNextToken(a)
}

func (a *ArraySlicePathToken) prevToken() Token {
	return a.prev
}

func (a *ArraySlicePathToken) setDefiniteUpdated(definiteUpdated bool) {
	a.definiteUpdated = definiteUpdated
}

func (a *ArraySlicePathToken) isDefiniteUpdated() bool {
	return a.definiteUpdated
}

func (a *ArraySlicePathToken) setDefinite(definite bool) {
	a.definite = definite
}

func (a *ArraySlicePathToken) isDefinite() bool {
	return a.definite
}

func (a *ArraySlicePathToken) setUpstreamUpdated(upstreamUpdated bool) {
	a.upstreamUpdated = upstreamUpdated
}

func (a *ArraySlicePathToken) isUpstreamUpdated() bool {
	return a.upstreamUpdated
}

func (a *ArraySlicePathToken) setUpstreamDefinite(upstreamDefinite bool) {
	a.upstreamDefinite = upstreamDefinite
}

func (a *ArraySlicePathToken) IsUpstreamDefinite() bool {
	return a.upstreamDefinite
}

func (a *ArraySlicePathToken) String() string {
	return tokenString(a)
}

func (a *ArraySlicePathToken) getUpstreamArrayIndex() int {
	return a.upstreamArrayIndex
}

func (a *ArraySlicePathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (a *ArraySlicePathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(a)
}

func (a *ArraySlicePathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(a, idx)
}

func (a *ArraySlicePathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(a, next)
}

func (a *ArraySlicePathToken) GetTokenCount() (int, error) {
	return tokenGetTokenCount(a)
}

func (a *ArraySlicePathToken) checkArrayModel(currentPath string, model interface{}, ctx *EvaluationContextImpl) (bool, error) {
	if model == nil {
		if !a.IsTokenDefinite() || common.UtilsSliceContains(ctx.Options(), common.OPTION_SUPPRESS_EXCEPTIONS) {
			return false, nil
		} else {
			return false, &common.PathNotFoundError{Message: "The path " + currentPath + " is null"}
		}
	}

	if ctx.JsonProvider().IsArray(model) {
		if a.IsUpstreamDefinite() || common.UtilsSliceContains(ctx.Options(), common.OPTION_SUPPRESS_EXCEPTIONS) {
			return false, nil
		} else {
			return false, &common.PathNotFoundError{Message: fmt.Sprintf("Filter: %s can only be applied to arrays. Current context is: %s", a, model)}
		}
	}
	return true, nil
}

func (a *ArraySlicePathToken) Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	checkPass, err := a.checkArrayModel(currentPath, model, ctx)
	if err != nil {
		return err
	}
	if checkPass {
		return nil
	}
	switch a.operation.OperationType() {
	case SLICE_FROM:
		return a.sliceFrom(currentPath, parent, model, ctx)
	case SLICE_TO:
		return a.sliceTo(currentPath, parent, model, ctx)
	case SLICE_BETWEEN:
		return a.sliceBetween(currentPath, parent, model, ctx)
	}
	return nil
}

func (a *ArraySlicePathToken) sliceFrom(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	length, err := ctx.JsonProvider().Length(model)
	if err != nil {
		return err
	}
	from := a.operation.From()
	if from < 0 {
		//calculate slice start from array length
		from = length + from
	}
	from = common.UtilsMaxInt(0, from)

	log.Printf("Slice from index on array with length: %d. From index: %d to: %d. Input: %s", length, from, length-1, common.UtilsToString(a))

	if length == 0 || from >= length {
		return nil
	}
	for i := from; i < length; i++ {
		err := tokenHandleArrayIndex(a, i, currentPath, model, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ArraySlicePathToken) sliceBetween(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	length, err := ctx.JsonProvider().Length(model)
	if err != nil {
		return err
	}
	from := a.operation.From()
	to := a.operation.To()

	to = common.UtilsMinInt(length, to)

	if from >= to || length == 0 {
		return nil
	}

	log.Printf("Slice between indexes on array with length: %d. From index: %d to: %d. Input: %s", length, from, to, common.UtilsToString(a))

	for i := from; i < to; i++ {
		err := tokenHandleArrayIndex(a, i, currentPath, model, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ArraySlicePathToken) sliceTo(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	length, err := ctx.JsonProvider().Length(model)
	if err != nil {
		return err
	}
	if length == 0 {
		return nil
	}
	to := a.operation.To()
	if to < 0 {
		//calculate slice end from array length
		to = length + to
	}
	to = common.UtilsMinInt(length, to)

	log.Printf("Slice to index on array with length: %d. From index: 0 to: %d. Input: %s", length, to, common.UtilsToString(a))

	for i := 0; i < to; i++ {
		err := tokenHandleArrayIndex(a, i, currentPath, model, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ArraySlicePathToken) GetPathFragment() string {
	return common.UtilsToString(a.operation)
}

func (*ArraySlicePathToken) IsTokenDefinite() bool {
	return false
}

func CreateArraySlicePathToken(operation *ArraySliceOperation) *ArraySlicePathToken {
	return &ArraySlicePathToken{
		defaultToken: &defaultToken{upstreamArrayIndex: -1},
		operation:    operation,
	}
}

// PredicatePathToken

type PredicatePathToken struct {
	*defaultToken
	predicates []common.Predicate
}

func (p *PredicatePathToken) SetPrev(prev Token) {
	p.prev = prev
}

func (p *PredicatePathToken) SetNext(next Token) {
	p.next = next
}

func (p *PredicatePathToken) GetNext() Token {
	return p.next
}

func (p *PredicatePathToken) isLeaf() bool {
	return p.next == nil
}

func (p *PredicatePathToken) isRoot() bool {
	return p.prev == nil
}

func (p *PredicatePathToken) nextToken() (Token, error) {
	return tokenNextToken(p)
}

func (p *PredicatePathToken) prevToken() Token {
	return p.prev
}

func (p *PredicatePathToken) setDefiniteUpdated(definiteUpdated bool) {
	p.definiteUpdated = definiteUpdated
}

func (p *PredicatePathToken) isDefiniteUpdated() bool {
	return p.definiteUpdated
}

func (p *PredicatePathToken) setDefinite(definite bool) {
	p.definite = definite
}

func (p *PredicatePathToken) isDefinite() bool {
	return p.definite
}

func (p *PredicatePathToken) setUpstreamUpdated(upstreamUpdated bool) {
	p.upstreamUpdated = upstreamUpdated
}

func (p *PredicatePathToken) isUpstreamUpdated() bool {
	return p.upstreamUpdated
}

func (p *PredicatePathToken) setUpstreamDefinite(upstreamDefinite bool) {
	p.upstreamDefinite = upstreamDefinite
}

func (p *PredicatePathToken) IsUpstreamDefinite() bool {
	return p.upstreamDefinite
}

func (p *PredicatePathToken) String() string {
	return tokenString(p)
}

func (p *PredicatePathToken) getUpstreamArrayIndex() int {
	return p.upstreamArrayIndex
}

func (p *PredicatePathToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	return tokenInvoke(pathFunction, currentPath, parent, model, ctx)
}

func (p *PredicatePathToken) IsPathDefinite() bool {
	return tokenIsPathDefinite(p)
}

func (p *PredicatePathToken) SetUpstreamArrayIndex(idx int) {
	tokenSetUpstreamArrayIndex(p, idx)
}

func (p *PredicatePathToken) appendTailToken(next Token) Token {
	return tokenAppendTailToken(p, next)
}

func (p *PredicatePathToken) GetTokenCount() (int, error) {
	return tokenGetTokenCount(p)
}

func (p *PredicatePathToken) Evaluate(currentPath string, ref common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	if ctx.JsonProvider().IsMap(model) {
		acceptResult, err := p.accept(model, ctx.RootDocument(), ctx.Configuration(), ctx)
		if err != nil {
			return err
		}
		if acceptResult {
			var op common.PathRef
			if ctx.ForUpdate() {
				op = ref
			} else {
				op = PathRefNoOp
			}
			if p.isLeaf() {
				if err = ctx.AddResult(currentPath, op, model); err != nil {
					return err
				}
			} else {
				next, err := p.nextToken()
				if err != nil {
					return err
				}
				return next.Evaluate(currentPath, op, model, ctx)
			}
		}
	} else if ctx.JsonProvider().IsArray(model) {
		idx := 0
		objects, err := ctx.JsonProvider().ToArray(model)
		if err != nil {
			return err
		}
		for _, idxModel := range objects {
			acceptResult, err := p.accept(idxModel, ctx.RootDocument(), ctx.Configuration(), ctx)
			if err != nil {
				return err
			}
			if acceptResult {
				err = tokenHandleArrayIndex(p, idx, currentPath, model, ctx)
				if err != nil {
					return err
				}
			}
			idx++
		}
	} else {
		if p.IsUpstreamDefinite() {
			return &common.InvalidPathError{Message: fmt.Sprintf("Filter: %s can not be applied to primitives. Current context is: %s", p, model)}
		}
	}
	return nil
}

func (p *PredicatePathToken) accept(obj interface{}, root interface{}, configuration *common.Configuration, evaluationContext *EvaluationContextImpl) (bool, error) {
	ctx := common.CreatePredicateContextImpl(obj, root, configuration, evaluationContext.DocumentEvalCache())

	for _, predicate := range p.predicates {
		pResult, err := predicate.Apply(ctx)
		if err != nil {
			return false, err
		}
		if !pResult {
			return false, nil
		}
		//TODO: err catch
	}
	return true, nil
}

func (p *PredicatePathToken) GetPathFragment() string {
	str := "["
	for i := 0; i < len(p.predicates); i++ {
		if i != 0 {
			str += ","
		}
		str += "?"
	}
	str += "]"
	return str
}

func (p *PredicatePathToken) IsTokenDefinite() bool {
	return false
}

func CreatePredicatePathToken(predicates []common.Predicate) *PredicatePathToken {
	return &PredicatePathToken{defaultToken: &defaultToken{upstreamArrayIndex: -1}, predicates: predicates}
}
