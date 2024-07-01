package reflector

import (
	"fmt"
	"reflect"
	"runtime"
)

// Func - .
type Func struct {
	ReflectValue reflect.Value
	Runtime      *runtime.Func
}

// ScanFunc - .
func ScanFunc(fn any) (Func, error) {
	if fn == nil {
		return Func{}, fmt.Errorf("no value found")
	}

	var rv = reflect.ValueOf(fn)

	if rv.Type().Kind() != reflect.Func {
		return Func{}, fmt.Errorf("detected value is not a function, passed: %T", fn)
	}

	ptr := uintptr(rv.UnsafePointer())
	if ptr == 0 {
		return Func{}, fmt.Errorf("passed an empty function pointer")
	}

	return Func{
		ReflectValue: rv,
		Runtime:      runtime.FuncForPC(ptr),
	}, nil
}
