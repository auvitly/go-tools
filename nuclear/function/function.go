package function

import (
	"fmt"
	"github.com/auvitly/go-tools/nuclear/impls"
	"reflect"
	"runtime"
	"unsafe"
)

type Info uintptr

func New(fn any) Info {
	if reflect.ValueOf(fn).Type().Kind() != reflect.Func {
		panic(fmt.Sprintf("%T", fn))
	}

	if fn == nil {
		return 0
	}

	if ptr := (*impls.Interface)(unsafe.Pointer(&fn)).Data; ptr != nil {
		return Info(*ptr)
	}

	return 0
}

func (i Info) Name() string {
	if i == 0 {
		return ""
	}

	return runtime.FuncForPC(uintptr(i)).Name()
}
