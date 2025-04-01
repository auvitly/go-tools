package impl

import (
	"unsafe"
)

type Interface struct {
	Type uintptr
	Data unsafe.Pointer
}
