package types

import "unsafe"

type Interface struct {
	Type uintptr
	Data unsafe.Pointer
}
