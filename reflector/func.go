package reflector

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Func - .
type Func struct {
	Name    string
	Pkg     string
	rv      reflect.Value
	runtime *runtime.Func
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

	rfn := runtime.FuncForPC(ptr)

	return Func{
		rv:      rv,
		runtime: rfn,
		Pkg:     pkgName(rfn),
		Name:    name(rfn),
	}, nil
}

func pkgName(fn *runtime.Func) string {
	var full = fn.Name()

	var (
		paths = strings.Split(full, "/")
		desc  = strings.Split(paths[len(paths)-1], ".")
	)

	return strings.Join(append(paths[:len(paths)-1], desc[0]), "/")
}

func name(fn *runtime.Func) string {
	var full = fn.Name()

	var (
		paths = strings.Split(full, "/")
		desc  = strings.Split(paths[len(paths)-1], ".")
	)

	return strings.Join(desc[1:], ".")
}
