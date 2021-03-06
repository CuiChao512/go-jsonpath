package path

import (
	"github.com/CuiChao512/go-jsonpath/jsonpath/common"
	"strings"
)

type noOpPathRef struct {
	parent interface{}
}

func (*noOpPathRef) GetAccessor() interface{} {
	return nil
}
func (*noOpPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	return nil
}

func (*noOpPathRef) Convert(mapFunction common.MapFunction, configuration *common.Configuration) error {
	return nil
}
func (*noOpPathRef) Delete(configuration *common.Configuration) error {
	return nil
}
func (*noOpPathRef) Add(newVal interface{}, configuration *common.Configuration) error {
	return nil
}
func (*noOpPathRef) Put(key string, newVal interface{}, configuration *common.Configuration) error {
	return nil
}
func (*noOpPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	return nil
}

func renameInMap(targetMap interface{}, oldKeyName string, newKeyName string, config *common.Configuration) error {
	if config.JsonProvider().IsMap(targetMap) {
		if config.JsonProvider().GetMapValue(targetMap, oldKeyName) == common.JsonProviderUndefined {
			return &common.PathNotFoundError{Message: "No results for Key " + oldKeyName + " found in map!"}
		}
		err := config.JsonProvider().SetProperty(&targetMap, newKeyName, config.JsonProvider().GetMapValue(targetMap, oldKeyName))
		if err != nil {
			return err
		}
		err = config.JsonProvider().RemoveProperty(&targetMap, oldKeyName)
		if err != nil {
			return err
		}
	} else {
		return &common.InvalidModificationError{Message: "Can only rename properties in a map"}
	}
	return nil
}

func (r *noOpPathRef) CompareTo(o common.PathRef) int {
	return strings.Compare(common.UtilsToString(r.GetAccessor()), common.UtilsToString(o.GetAccessor())) * -1
}

func isTargetInvalid(target interface{}) bool {
	return target == common.JsonProviderUndefined || target == nil
}

var PathRefNoOp common.PathRef = &noOpPathRef{}

func CreateObjectPropertyPathRef(obj interface{}, property string) common.PathRef {
	om := &objectPropertyPathRef{}
	om.parent = obj
	om.property = property
	return om
}

func CreateObjectMultiPropertyPathRef(obj interface{}, properties []string) common.PathRef {
	om := &objectMultiPropertyPathRef{}
	om.parent = obj
	om.properties = properties
	return om
}

func CreateArrayIndexPathRef(array interface{}, index int) common.PathRef {
	a := &arrayIndexPathRef{}
	a.parent = array
	a.index = index
	return a
}

func CreateRootPathRef(root interface{}) common.PathRef {
	d := &noOpPathRef{parent: root}
	return d
}

// rootPathRef -----------
type rootPathRef struct {
	parent interface{}
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
		length, err := config.JsonProvider().Length(r.parent)
		if err != nil {
			return err
		}
		return config.JsonProvider().SetArrayIndex(&r.parent, length, newVal)
	} else {
		return &common.InvalidModificationError{Message: "Invalid add operation. $ is not an array"}
	}
}

func (r *rootPathRef) Put(key string, newVal interface{}, configuration *common.Configuration) error {
	if configuration.JsonProvider().IsMap(r.parent) {
		return configuration.JsonProvider().SetProperty(&r.parent, key, newVal)
	} else {
		return &common.InvalidModificationError{Message: "Invalid put operation. $ is not a map"}
	}
}

func (r *rootPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	target := r.parent
	if isTargetInvalid(target) {
		return nil
	}
	err := renameInMap(target, oldKeyName, newKeyName, configuration)
	return err
}

func (r *rootPathRef) CompareTo(o common.PathRef) int {
	return strings.Compare(common.UtilsToString(r.GetAccessor()), common.UtilsToString(o.GetAccessor())) * -1
}

// arrayIndexPathRef
type arrayIndexPathRef struct {
	parent interface{}
	index  int
}

func (r *arrayIndexPathRef) GetAccessor() interface{} {
	return r.index
}

func (r *arrayIndexPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	return configuration.JsonProvider().SetArrayIndex(&r.parent, r.index, newVal)
}

func (r *arrayIndexPathRef) Convert(mapFunction common.MapFunction, configuration *common.Configuration) error {
	currentValue, err := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	if err != nil {
		return err
	}
	return configuration.JsonProvider().SetArrayIndex(&r.parent, r.index, mapFunction.Map(currentValue, configuration))
}

func (r *arrayIndexPathRef) Delete(configuration *common.Configuration) error {
	return configuration.JsonProvider().RemoveProperty(&r.parent, r.index)
}

func (r *arrayIndexPathRef) Add(value interface{}, configuration *common.Configuration) error {
	target, err := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	if err != nil {
		return err
	}
	if isTargetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsArray(target) {
		return configuration.JsonProvider().SetProperty(&target, nil, value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add to an array"}
	}
}

func (r *arrayIndexPathRef) Put(key string, value interface{}, configuration *common.Configuration) error {
	target, err := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	if err != nil {
		return err
	}
	if isTargetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsMap(target) {
		return configuration.JsonProvider().SetProperty(&target, key, value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add properties to a map"}
	}
}

func (r *arrayIndexPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	target, err := configuration.JsonProvider().GetArrayIndex(r.parent, r.index)
	if err != nil {
		return err
	}
	if isTargetInvalid(target) {
		return nil
	}
	return renameInMap(target, oldKeyName, newKeyName, configuration)
}

func (r *arrayIndexPathRef) CompareTo(o common.PathRef) int {
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
	parent   interface{}
	property string
}

func (r *objectPropertyPathRef) GetAccessor() interface{} {
	return r.property
}

func (r *objectPropertyPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	return configuration.JsonProvider().SetProperty(&r.parent, r.property, newVal)
}

func (r *objectPropertyPathRef) Convert(mapFunction common.MapFunction, configuration *common.Configuration) error {
	currentValue := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	return configuration.JsonProvider().SetProperty(&r.parent, r.property, mapFunction.Map(currentValue, configuration))
}

func (r *objectPropertyPathRef) Delete(configuration *common.Configuration) error {
	return configuration.JsonProvider().RemoveProperty(r.parent, r.property)
}

func (r *objectPropertyPathRef) Add(value interface{}, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	if isTargetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsArray(target) {
		length, err := configuration.JsonProvider().Length(target)
		if err != nil {
			return err
		}
		return configuration.JsonProvider().SetArrayIndex(&target, length, value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add to an array"}
	}
}

func (r *objectPropertyPathRef) Put(keyStr string, value interface{}, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	if isTargetInvalid(target) {
		return nil
	}
	if configuration.JsonProvider().IsMap(target) {
		return configuration.JsonProvider().SetProperty(&target, keyStr, value)
	} else {
		return &common.InvalidModificationError{Message: "Can only add properties to a map"}
	}
}

func (r *objectPropertyPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *common.Configuration) error {
	target := configuration.JsonProvider().GetMapValue(r.parent, r.property)
	if isTargetInvalid(target) {
		return nil
	}
	return renameInMap(target, oldKeyName, newKeyName, configuration)
}

func (r *objectPropertyPathRef) CompareTo(o common.PathRef) int {
	return strings.Compare(common.UtilsToString(r.GetAccessor()), common.UtilsToString(o.GetAccessor())) * -1
}

// objectMultiPropertyPathRef
type objectMultiPropertyPathRef struct {
	parent     interface{}
	properties []string
}

func (r *objectMultiPropertyPathRef) GetAccessor() interface{} {
	return common.UtilsJoin("&&", "", r.properties)
}

func (r *objectMultiPropertyPathRef) Set(newVal interface{}, configuration *common.Configuration) error {
	for _, property := range r.properties {
		err := configuration.JsonProvider().SetProperty(&r.parent, property, newVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *objectMultiPropertyPathRef) Convert(mapFunction common.MapFunction, config *common.Configuration) error {
	for _, property := range r.properties {
		currentValue := config.JsonProvider().GetMapValue(r.parent, property)
		if currentValue != common.JsonProviderUndefined {
			err := config.JsonProvider().SetProperty(&r.parent, property, mapFunction.Map(currentValue, config))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *objectMultiPropertyPathRef) Delete(configuration *common.Configuration) error {
	for _, property := range r.properties {
		err := configuration.JsonProvider().RemoveProperty(r.parent, property)
		if err != nil {
			return err
		}
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

func (r *objectMultiPropertyPathRef) CompareTo(o common.PathRef) int {
	return strings.Compare(common.UtilsToString(r.GetAccessor()), common.UtilsToString(o.GetAccessor())) * -1
}
