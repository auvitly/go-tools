package memory

import (
	"github.com/auvitly/go-tools/nuclear/impls"
	"unsafe"
)

// As - performs a representation of a memory fragment as a slice of bytes
func As(ptr uintptr, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&impls.Slice{
		Data: ptr,
		Len:  size,
		Cap:  size,
	}))
}

// Clone - copy of a memory fragment as a slice of bytes.
func Clone(ptr uintptr, size int) []byte {
	var (
		buf  = make([]byte, 0, size)
		data = As(ptr, size)
	)

	return append(buf, data...)
}

// Copy - copies data from src to dst by number of bytes.
func Copy(dst uintptr, src uintptr, size int) {
	copy(As(dst, size), As(src, size))
}

// Write - overwrites memory by ptr with the following data bytes.
func Write(ptr uintptr, data []byte) {
	SetProtect(ptr, len(data), ProtectModeReadWrite)
	defer func() { SetProtect(ptr, len(data), ProtectModeRead) }()

	Copy(ptr, uintptr(unsafe.Pointer(&data[0])), len(data))
}
