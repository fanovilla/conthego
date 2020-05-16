package conthego

import (
	"errors"
	"log"
	"reflect"
)

//https://github.com/kubernetes/kops/blob/master/util/pkg/reflectutils/walk.go
func InvokeMethod(target interface{}, name string, args ...interface{}) ([]reflect.Value, error) {
	v := reflect.ValueOf(target)

	method, found := v.Type().MethodByName(name)
	if !found {
		return nil, errors.New("method not found:" + name)
	}

	var argValues []reflect.Value
	for _, a := range args {
		argValues = append(argValues, reflect.ValueOf(a))
	}
	log.Printf("Calling method %s on %T", method.Name, target)
	m := v.MethodByName(method.Name)
	rv := m.Call(argValues)
	return rv, nil
}

func formatAtom(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint()
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		return v.String()
	case reflect.Slice:
		return v.Slice(0, 1)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}
