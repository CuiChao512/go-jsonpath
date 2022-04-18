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

func (t *defaultToken) handleObjectProperty(currentPath string, model interface{}, ctx *EvaluationContextImpl, properties []string) error {

	if len(properties) == 1 {
		property := properties[0]
		evalPath := common.UtilsConcat(currentPath, "['", property, "']")
		propertyVal := pathTokenReadObjectProperty(property, model, ctx)
		if propertyVal == common.JsonProviderUndefined {
			// Conditions below heavily depend on current token type (and its logic) and are not "universal",
			// so this code is quite dangerous (I'd rather rewrite it & move to PropertyPathToken and implemented
			// WildcardPathToken as a dynamic multi prop case of PropertyPathToken).
			// Better safe than sorry.
			switch common.UtilsGetPtrElem(t).(type) {
			case PropertyPathToken:
			default:
				return errors.New("only PropertyPathToken is supported")
			}

			if t.isLeaf() {

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
				if !(t.IsUpstreamDefinite() && t.IsTokenDefinite()) &&
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
		if t.isLeaf() {
			idx := "[" + common.UtilsToString(t.upstreamArrayIndex) + "]"
			root, err := ctx.GetRoot()
			if err != nil {
				return err
			}
			if idx == "[-1]" || root.GetTail().prevToken().GetPathFragment() == idx {
				ctx.AddResult(evalPath, ref, propertyVal)
			}
		} else {
			next, _ := t.nextToken()
			err := next.Evaluate(evalPath, ref, propertyVal, ctx)
			if err != nil {
				return err
			}
		}
	} else {
		evalPath := currentPath + "[" + common.UtilsJoin(", ", "'", properties) + "]"

		if !t.isLeaf() {
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
			ctx.JsonProvider().SetProperty(merged, property, propertyVal)
		}
		var pathRef common.PathRef
		if ctx.ForUpdate() {
			pathRef = CreateObjectMultiPropertyPathRef(model, properties)
		} else {
			pathRef = PathRefNoOp
		}
		ctx.AddResult(evalPath, pathRef, merged)
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

func (t *defaultToken) handleArrayIndex(index int, currentPath string, model interface{}, ctx *EvaluationContextImpl) error {
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

	if t.isLeaf() {
		ctx.AddResult(evalPath, pathRef, evalHit)
	} else {
		next, err := t.nextToken()
		if err != nil {
			return err
		}
		return next.Evaluate(evalPath, pathRef, evalHit, ctx)
	}
	return nil
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

func (t *defaultToken) Invoke(pathFunction function.PathFunction, currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	result, err := pathFunction.Invoke(currentPath, parent, model, ctx, nil)
	if err != nil {
		return err
	}
	ctx.AddResult(currentPath, parent, result)
	return nil
}

func (t *defaultToken) nextToken() (Token, error) {
	if t.isLeaf() {
		return nil, &common.IllegalStateException{Message: "Current path token is a leaf"}
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

func (t *defaultToken) Evaluate(currentPath string, parent common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
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

func (r *RootPathToken) Evaluate(currentPath string, ref common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
	if r.isLeaf() {
		var op common.PathRef
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
		return p.handleObjectProperty(currentPath, model, ctx, p.properties)
	}

	if !p.MultiPropertyIterationCase() {
		return errors.New("")
	}

	for _, property := range p.properties {
		err := p.handleObjectProperty(currentPath, model, ctx, []string{property})
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePropertyPathToken(properties []string, stringDelimiter string) *PropertyPathToken {
	return &PropertyPathToken{properties: properties, stringDelimiter: stringDelimiter}
}

//WildCardPathToken

type WildcardPathToken struct {
	*defaultToken
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
			err := w.handleObjectProperty(currentPath, model, ctx, []string{property})
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
			err := w.handleArrayIndex(idx, currentPath, model, ctx)

			if err != nil && common.UtilsSliceContains(ctx.Options(), common.OPTION_REQUIRE_PROPERTIES) {
				return err
			}
		}
	}
	return nil
}

func CreateWildcardPathToken() *WildcardPathToken {
	return &WildcardPathToken{}
}

// ScanPathToken -----

type ScanPredicate interface {
	matches(model interface{}) (bool, error)
}

type defaultScanPredicate struct {
}

func (*defaultScanPredicate) matches(model interface{}) (bool, error) {
	return false, nil
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
	*defaultScanPredicate
}

func (*wildCardPathTokenPredicate) matches(model interface{}) (bool, error) {
	return true, nil
}

type arrayPathTokenPredicate struct {
	*defaultScanPredicate
	ctx *EvaluationContextImpl
}

func (a *arrayPathTokenPredicate) matches(model interface{}) (bool, error) {
	return a.ctx.JsonProvider().IsArray(model), nil
}

type propertyPathTokenPredicate struct {
	*defaultScanPredicate
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
	case *arrayPathToken:
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

// ArrayPathToken

type arrayPathToken struct {
	*defaultToken
}

func (a *arrayPathToken) checkArrayModel(currentPath string, model interface{}, ctx *EvaluationContextImpl) (bool, error) {
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

type ArrayIndexPathToken struct {
	*arrayPathToken
	arrayIndexOperation *ArrayIndexOperation
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
		return a.handleArrayIndex(a.arrayIndexOperation.Indexes()[0], currentPath, model, ctx)
	} else {
		for _, idx := range a.arrayIndexOperation.Indexes() {
			err = a.handleArrayIndex(idx, currentPath, model, ctx)
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
	return &ArrayIndexPathToken{arrayIndexOperation: arrayIndexOperation}
}

// ArraySlicePathToken -----
type ArraySlicePathToken struct {
	*arrayPathToken
	operation *ArraySliceOperation
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
		err := a.handleArrayIndex(i, currentPath, model, ctx)
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
		err := a.handleArrayIndex(i, currentPath, model, ctx)
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
		err := a.handleArrayIndex(i, currentPath, model, ctx)
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
		operation: operation,
	}
}

// PredicatePathToken

type PredicatePathToken struct {
	*defaultToken
	predicates []common.Predicate
}

func (p *PredicatePathToken) evaluate(currentPath string, ref common.PathRef, model interface{}, ctx *EvaluationContextImpl) error {
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
				ctx.AddResult(currentPath, op, model)
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
				err = p.handleArrayIndex(idx, currentPath, model, ctx)
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
	return &PredicatePathToken{predicates: predicates}
}
