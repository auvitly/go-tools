package memory

import (
	"github.com/auvitly/go-tools/nuclear/impls"
	"unsafe"
)

func As(ptr uintptr, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&impls.Slice{
		Data: ptr,
		Len:  size,
		Cap:  size,
	}))
}

func Scan(ptr uintptr, size int) []byte {
	var (
		buf  = make([]byte, 0, size)
		data = As(ptr, size)
	)

	return append(buf, data...)
}

func Copy(dst uintptr, src uintptr, size int) {
	copy(As(dst, size), As(src, size))
}

func Write(ptr uintptr, data []byte) {
	SetProtect(ptr, len(data), ProtectModeReadWrite)
	defer func() { SetProtect(ptr, len(data), ProtectModeRead) }()

	Copy(ptr, uintptr(unsafe.Pointer(&data[0])), len(data))
}
