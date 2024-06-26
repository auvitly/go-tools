//go:build !windows && !darwin
// +build !windows,!darwin

package memory

import (
	"fmt"
	"syscall"
)

var (
	pageSize = syscall.Getpagesize()
	modes    = map[ProtectMode]int{
		ProtectModeRead:      syscall.PROT_READ | syscall.PROT_EXEC,
		ProtectModeReadWrite: syscall.PROT_READ | syscall.PROT_EXEC | syscall.PROT_WRITE,
	}
)

func pageFrom(ptr uintptr) uintptr {
	return ptr & ^uintptr(pageSize-1)
}

func SetProtect(ptr uintptr, size int, mode ProtectMode) {
	protect, ok := modes[mode]
	if !ok {
		panic(ErrUnsupportedProtectMode)
	}

	var page = ptr & ^uintptr(pageSize-1)

	for i := page; i < ptr+uintptr(size); i += uintptr(pageSize) {
		if err := syscall.Mprotect(As(i, pageSize), protect); err != nil {
			panic(fmt.Sprintf("syscall.Mprotect: %v", err))
		}
	}
}
