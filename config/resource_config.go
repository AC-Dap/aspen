package config

import (
	"aspen/router"
	"fmt"
	"log"
	"reflect"
)

type ResourceConfig struct {
	ResourceType string
	Params       []any
}

// We have no way of representing the allowable constructor with go types since we allow arbitrary parameters,
// so we use `any` here and do validation when registering constructors.
type ResourceContructor = any
type ResourceMap = map[string]ResourceContructor

var globalResourceMap ResourceMap = make(ResourceMap)

func RegisterResourceConstructor(resourceType string, constructor ResourceContructor) error {
	// Check if this type alrady exists
	if _, ok := globalResourceMap[resourceType]; ok {
		return fmt.Errorf("\"%s\" resource constructor has already been registered", resourceType)
	}

	// Check that constructor is actually a function
	constructorType := reflect.TypeOf(constructor)
	if constructorType.Kind() != reflect.Func {
		return fmt.Errorf("\"%s\" resource constructor is not a function", resourceType)
	}

	// Check that the return type is router.Resource
	expectedReturnType := reflect.TypeOf((*router.Resource)(nil)).Elem()
	if constructorType.NumOut() != 1 || !constructorType.Out(0).AssignableTo(expectedReturnType) {
		return fmt.Errorf("\"%s\" resource constructor does not return *router.Resource", resourceType)
	}

	log.Printf("Registered \"%s\" resource constructor", resourceType)
	globalResourceMap[resourceType] = constructor
	return nil
}

func (rc ResourceConfig) Parse() (router.Resource, error) {
	constructor, ok := globalResourceMap[rc.ResourceType]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" resource constructor", rc.ResourceType)
	}

	// Verify that rc.Params matches constructor params
	constructorType := reflect.TypeOf(constructor)
	if len(rc.Params) != constructorType.NumIn() {
		return nil, fmt.Errorf("\"%s\" resource constructor expects %d arguments, but got %d",
			rc.ResourceType, constructorType.NumIn(), len(rc.Params))
	}

	args := make([]reflect.Value, constructorType.NumIn())
	for i, param := range rc.Params {
		expectedParamType := constructorType.In(i)
		castParam, err := CastParam(reflect.ValueOf(param), expectedParamType)
		if err != nil {
			return nil, fmt.Errorf("mismatched argument %d for \"%s\" resource constructor: %w", i, rc.ResourceType, err)
		}

		args[i] = castParam
	}

	// Call constructor
	ret := reflect.ValueOf(constructor).Call(args)

	// Cast return value to *router.Resource and return
	// Since we validated constructor, this should not panic
	newResource, ok := ret[0].Interface().(router.Resource)
	if !ok {
		return nil, fmt.Errorf("unable to cast return value from \"%s\" resource constructor, something went wrong",
			rc.ResourceType)
	}

	return newResource, nil
}

/*
 * Tries to cast a into parameter of type b.
 */
func CastParam(param reflect.Value, expectedType reflect.Type) (reflect.Value, error) {
	// First see if we can directly assign
	if param.Type().AssignableTo(expectedType) {
		return param, nil
	}

	// May need to upcast `any` types
	// Work with the concrete values instead of interfaces
	if param.Kind() == reflect.Interface {
		param = param.Elem()
	}

	if param.Kind() == expectedType.Kind() {
		switch param.Kind() {
		case reflect.String:
			return reflect.ValueOf(param.String()), nil

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(param.Int()), nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.ValueOf(param.Uint()), nil

		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(param.Float()), nil

		case reflect.Bool:
			return reflect.ValueOf(param.Bool()), nil

		case reflect.Slice:
			castSlice := reflect.MakeSlice(expectedType, param.Len(), param.Len())
			// Cast and copy over each element
			for i := range param.Len() {
				el := param.Index(i)
				castEl, err := CastParam(el, expectedType.Elem())
				if err != nil {
					return reflect.ValueOf(nil), fmt.Errorf("cannot convert array element of type %s into %s: %w",
						el.Type().String(), expectedType.Elem().String(), err)
				}
				castSlice.Index(i).Set(castEl)
			}
			return castSlice, nil
		}
	}

	return reflect.ValueOf(nil), fmt.Errorf("cannot convert element of kind %s into %s",
		param.Type().String(), expectedType.String())
}
