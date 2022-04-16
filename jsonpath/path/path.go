package path

import (
	"cuichao.com/go-jsonpath/jsonpath/common"
	"cuichao.com/go-jsonpath/jsonpath/evaluationContext"
	"strings"
)

type Path interface {
	Evaluate(document interface{}, rootDocument interface{}, configuration *common.Configuration) (evaluationContext.EvaluationContext, error)
	EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *common.Configuration, forUpdate bool) evaluationContext.EvaluationContext
	String() string
	IsDefinite() bool
	IsFunctionPath() bool
	IsRootPath() bool
}

type Ref interface {
	GetAccessor() interface{}
	Set(newVal interface{}, configuration *common.Configuration) error
	Convert(mapFunction common.MapFunction, configuration *common.Configuration) error
	Delete(configuration *common.Configuration) error
	Add(newVal interface{}, configuration *common.Configuration) error
	Put(key string, newVal interface{}, configuration *common.Configuration) error
	RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error
	CompareTo(o Ref) int
}

type defaultRef struct {
	parent interface{}
}

func (*defaultRef) GetAccessor() interface{} {
	return nil
}
func (*defaultRef) Set(newVal interface{}, configuration *common.Configuration) error {
	return nil
}

func (*defaultRef) Convert(mapFunction common.MapFunction, configuration *common.Configuration) error {
	return nil
}
func (*defaultRef) Delete(configuration *common.Configuration) error {
	return nil
}
func (*defaultRef) Add(newVal interface{}, configuration *common.Configuration) error {
	return nil
}
func (*defaultRef) Put(key string, newVal interface{}, configuration *common.Configuration) error {
	return nil
}
func (*defaultRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	return nil
}

func (r *defaultRef) renameInMap(targetMap interface{}, oldKeyName string, newKeyName string, config *common.Configuration) error {
	if config.JsonProvider().IsMap(targetMap) {
		if config.JsonProvider().GetMapValue(targetMap, oldKeyName) == common.JsonProviderUndefined {
			return &common.PathNotFoundError{Message: "No results for Key " + oldKeyName + " found in map!"}
		}
		config.JsonProvider().SetProperty(targetMap, newKeyName, config.JsonProvider().GetMapValue(targetMap, oldKeyName))
		config.JsonProvider().RemoveProperty(targetMap, oldKeyName)
	} else {
		return &common.InvalidModificationError{Message: "Can only rename properties in a map"}
	}
	return nil
}

func (r *defaultRef) targetInvalid(target interface{}) bool {
	return target == common.JsonProviderUndefined || target == nil
}

func (r *defaultRef) CompareTo(o Ref) int {
	return strings.Compare(common.UtilsToString(r.GetAccessor()), common.UtilsToString(o.GetAccessor())) * -1
}

var PathRefNoOp Ref = &defaultRef{}

func CreateObjectPropertyPathRef(obj interface{}, property string) Ref {
	om := &objectPropertyPathRef{}
	om.parent = obj
	om.property = property
	return om
}

func CreateObjectMultiPropertyPathRef(obj interface{}, properties []string) Ref {
	om := &objectMultiPropertyPathRef{}
	om.parent = obj
	om.properties = properties
	return om
}

func CreateArrayIndexPathRef(array interface{}, index int) Ref {
	a := &arrayIndexPathRef{}
	a.parent = array
	a.index = index
	return a
}

func CreateRootPathRef(root interface{}) Ref {
	d := &defaultRef{parent: root}
	return d
}

// rootPathRef -----------
type rootPathRef struct {
	*defaultRef
}

func (*rootPathRef) GetAccessor() interface{} {
	return "$"
}
func (*rootPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	return &common.InvalidModificationError{Message: "Invalid set operation"}
}

func (*rootPathRef) Convert(mapFunction common.MapFunction, configuration *common.Configuration) error {
	return &common.InvalidModificationError{Message: "Invalid map operation"}
}

func (*rootPathRef) Delete(configuration *common.Configuration) error {
	return &common.InvalidModificationError{Message: "Invalid delete operation"}
}

func (r *rootPathRef) Add(newVal interface{}, config *common.Configuration) error {
	if config.JsonProvider().IsArray(r.parent) {
		config.JsonProvider().SetArrayIndex(r.parent, config.JsonProvider().Length(r.parent), newVal)
		return nil
	} else {
		return &common.InvalidModificationError{Message: "Invalid add operation. $ is not an array"}
	}
}

func (r *rootPathRef) Put(key string, newVal interface{}, configuration *common.Configuration) error {
	if configuration.JsonProvider().IsMap(r.parent) {
		configuration.JsonProvider().SetProperty(r.parent, key, newVal)
		return nil
	} else {
		return &common.InvalidModificationError{Message: "Invalid put operation. $ is not a map"}
	}
}

func (r *rootPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	target := r.parent
	if r.targetInvalid(target) {
		return nil
	}
	err := r.renameInMap(target, oldKeyName, newKeyName, configuration)
	return err
}

// arrayIndexPathRef
type arrayIndexPathRef struct {
	*defaultRef
	index int
}

func (r *arrayIndexPathRef) GetAccessor() interface{} {
	return r.index
}

func (r *arrayIndexPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	configuration.JsonProvider().SetArrayIndex(r.parent, r.index, newVal)
	return nil
}

func (r *arrayIndexPathRef) Convert(mapFunction common.MapFunction, configuration *common.Configuration) error {
	currentValue := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	configuration.JsonProvider().SetArrayIndex(r.parent, r.index, mapFunction.Map(currentValue, configuration))
	return nil
}

func (r *arrayIndexPathRef) Delete(configuration *common.Configuration) error {
	configuration.JsonProvider().RemoveProperty(r.parent, r.index)
	return nil
}

func (r *arrayIndexPathRef) Add(value interface{}, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	if r.targetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsArray(target) {
		configuration.JsonProvider().SetProperty(target, nil, value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add to an array"}
	}
	return nil
}

func (r *arrayIndexPathRef) Put(key string, value interface{}, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	if r.targetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsMap(target) {
		configuration.JsonProvider().SetProperty(target, key, value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add properties to a map"}
	}
	return nil
}

func (r *arrayIndexPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	if r.targetInvalid(target) {
		return nil
	}
	return r.renameInMap(target, oldKeyName, newKeyName, configuration)
}

func (r *arrayIndexPathRef) CompareTo(o Ref) int {
	switch o.(type) {
	case *arrayIndexPathRef:
		pf, _ := o.(*arrayIndexPathRef)
		return pf.index - r.index
	default:
		return r.CompareTo(o)
	}
}

// objectPropertyPathRef
type objectPropertyPathRef struct {
	*defaultRef
	property string
}

func (r *objectPropertyPathRef) GetAccessor() interface{} {
	return r.property
}

func (r *objectPropertyPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	configuration.JsonProvider().SetProperty(r.parent, r.property, newVal)
	return nil
}

func (r *objectPropertyPathRef) Convert(mapFunction common.MapFunction, configuration *common.Configuration) error {
	currentValue := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	configuration.JsonProvider().SetProperty(r.parent, r.property, mapFunction.Map(currentValue, configuration))
	return nil
}

func (r *objectPropertyPathRef) Delete(configuration *common.Configuration) error {
	configuration.JsonProvider().RemoveProperty(r.parent, r.property)
	return nil
}

func (r *objectPropertyPathRef) Add(value interface{}, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	if r.targetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsArray(target) {
		configuration.JsonProvider().SetArrayIndex(target, configuration.JsonProvider().Length(target), value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add to an array"}
	}
	return nil
}

func (r *objectPropertyPathRef) Put(keyStr string, value interface{}, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	if r.targetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsMap(target) {
		configuration.JsonProvider().SetProperty(target, keyStr, value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add properties to a map"}
	}
	return nil
}

func (r *objectPropertyPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	if r.targetInvalid(target) {
		return nil
	}
	return r.renameInMap(target, oldKeyName, newKeyName, configuration)
}

// objectMultiPropertyPathRef
type objectMultiPropertyPathRef struct {
	*defaultRef
	properties []string
}

func (r *objectMultiPropertyPathRef) GetAccessor() interface{} {
	return common.UtilsJoin("&&", "", r.properties)
}

func (r *objectMultiPropertyPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	for _, property := range r.properties {
		configuration.JsonProvider().SetProperty(r.parent, property, newVal)
	}
	return nil
}

func (r *objectMultiPropertyPathRef) Convert(mapFunction common.MapFunction, config *common.Configuration) error {
	for _, property := range r.properties {
		currentValue := config.JsonProvider().GetMapValue(r.parent, property)
		if currentValue != common.JsonProviderUndefined {
			config.JsonProvider().SetProperty(r.parent, property, mapFunction.Map(currentValue, config))
		}
	}
	return nil
}

func (r *objectMultiPropertyPathRef) Delete(configuration *common.Configuration) error {
	for _, property := range r.properties {
		configuration.JsonProvider().RemoveProperty(r.parent, property)
	}
	return nil
}

func (*objectMultiPropertyPathRef) Add(newVal interface{}, configuration *common.Configuration) error {
	return &common.InvalidModificationError{Message: "Add can not be performed to multiple properties"}
}

func (*objectMultiPropertyPathRef) Put(key string, newVal interface{}, configuration *common.Configuration) error {
	return &common.InvalidModificationError{Message: "Put can not be performed to multiple properties"}
}

func (*objectMultiPropertyPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	return &common.InvalidModificationError{Message: "Rename can not be performed to multiple properties"}
}
