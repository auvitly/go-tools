package memory

import (
	"reflect"
	"unsafe"
)

func As(ptr uintptr, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: ptr,
		Len:  size,
		Cap:  size,
	}))
}

func Scan(ptr uintptr, size int) []byte {
	return append(make([]byte, 0, size), As(ptr, size)...)
}

func Copy(dst uintptr, src uintptr, size int) {
	copy(As(dst, size), As(src, size))
}

func Write(ptr uintptr, data []byte) {
	SetProtect(ptr, len(data), ProtectModeReadWrite)
	defer func() { SetProtect(ptr, len(data), ProtectModeRead) }()

	Copy(ptr, uintptr(unsafe.Pointer(&data[0])), len(data))
}
