package path

import (
	"cuichao.com/go-jsonpath/jsonpath"
	"strings"
)

type Path interface {
	Evaluate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration) (jsonpath.EvaluationContext, error)
	EvaluateForUpdate(document interface{}, rootDocument interface{}, configuration *jsonpath.Configuration, forUpdate bool) jsonpath.EvaluationContext
	String() string
	IsDefinite() bool
	IsFunctionPath() bool
	IsRootPath() bool
}

type Ref interface {
	GetAccessor() interface{}
	Set(newVal interface{}, configuration *jsonpath.Configuration) error
	Convert(mapFunction jsonpath.MapFunction, configuration *jsonpath.Configuration) error
	Delete(configuration *jsonpath.Configuration) error
	Add(newVal interface{}, configuration *jsonpath.Configuration) error
	Put(key string, newVal interface{}, configuration *jsonpath.Configuration) error
	RenameKey(oldKeyName string, newKeyName string, configuration *jsonpath.Configuration) error
	CompareTo(o Ref) int
}

type defaultRef struct {
	parent interface{}
}

func (*defaultRef) GetAccessor() interface{} {
	return nil
}
func (*defaultRef) Set(newVal interface{}, configuration *jsonpath.Configuration) error {
	return nil
}

func (*defaultRef) Convert(mapFunction jsonpath.MapFunction, configuration *jsonpath.Configuration) error {
	return nil
}
func (*defaultRef) Delete(configuration *jsonpath.Configuration) error {
	return nil
}
func (*defaultRef) Add(newVal interface{}, configuration *jsonpath.Configuration) error {
	return nil
}
func (*defaultRef) Put(key string, newVal interface{}, configuration *jsonpath.Configuration) error {
	return nil
}
func (*defaultRef) RenameKey(oldKeyName string, newKeyName string, configuration *jsonpath.Configuration) error {
	return nil
}

func (r *defaultRef) renameInMap(targetMap interface{}, oldKeyName string, newKeyName string, configuration *jsonpath.Configuration) error {
	if configuration.JsonProvider().IsMap(targetMap) {
		if configuration.JsonProvider().GetMapValue(targetMap, oldKeyName) == jsonpath.JsonProviderUndefined {
			return &jsonpath.PathNotFoundError{Message: "No results for Key " + oldKeyName + " found in map!"}
		}
		configuration.JsonProvider().SetProperty(targetMap, newKeyName, configuration.JsonProvider().GetMapValue(targetMap, oldKeyName))
		configuration.JsonProvider().RemoveProperty(targetMap, oldKeyName)
	} else {
		return &jsonpath.InvalidModificationError{Message: "Can only rename properties in a map"}
	}
	return nil
}

func (r *defaultRef) targetInvalid(target interface{}) bool {
	return target == jsonpath.JsonProviderUndefined || target == nil
}

func (r *defaultRef) CompareTo(o Ref) int {
	return strings.Compare(jsonpath.UtilsToString(r.GetAccessor()), jsonpath.UtilsToString(o.GetAccessor())) * -1
}

var PATH_REF_NO_OP Ref = &defaultRef{}

func CreateObjectPropertyPathRef(obj interface{}, property string) Ref {
	return &objectPropertyPathRef(obj, property)
}

func CreateObjectMultiPropertyPathRef(obj interface{}, properties []string) Ref {
	return &objectMultiPropertyPathRef(obj, property)
}

func CreateArrayIndexPathRef(array interface{}, index int) Ref {
	return &objectArrayIndexPathRef(obj, property)
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
func (*rootPathRef) Set(newVal interface{}, configuration *jsonpath.Configuration) error {
	return &jsonpath.InvalidModificationError{Message: "Invalid set operation"}
}

func (*rootPathRef) Convert(mapFunction jsonpath.MapFunction, configuration *jsonpath.Configuration) error {
	return &jsonpath.InvalidModificationError{Message: "Invalid map operation"}
}

func (*rootPathRef) Delete(configuration *jsonpath.Configuration) error {
	return &jsonpath.InvalidModificationError{Message: "Invalid delete operation"}
}

func (r *rootPathRef) Add(newVal interface{}, configuration *jsonpath.Configuration) error {
	if configuration.JsonProvider().IsArray(r.parent) {
		configuration.JsonProvider().SetArrayIndex(r.parent, configuration.JsonProvider().Length(r.parent), newVal)
		return nil
	} else {
		return &jsonpath.InvalidModificationError{Message: "Invalid add operation. $ is not an array"}
	}
}

func (r *rootPathRef) Put(key string, newVal interface{}, configuration *jsonpath.Configuration) error {
	if configuration.JsonProvider().IsMap(r.parent) {
		configuration.JsonProvider().SetProperty(r.parent, key, newVal)
		return nil
	} else {
		return &jsonpath.InvalidModificationError{Message: "Invalid put operation. $ is not a map"}
	}
}

func (r *rootPathRef) RenameKey(oldKeyName string, newKeyName string, configuration *jsonpath.Configuration) error {
	target := r.parent
	if r.targetInvalid(target) {
		return nil
	}
	err := r.renameInMap(target, oldKeyName, newKeyName, configuration)
	return err
}
