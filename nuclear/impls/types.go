package impls

import "unsafe"

type Interface struct {
	Type uintptr
	Data *unsafe.Pointer
}

type Slice struct {
	Data uintptr
	Len  int
	Cap  int
}
