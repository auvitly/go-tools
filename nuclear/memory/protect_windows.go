//go:build windows
// +build windows

package memory

import (
	"syscall"
	"unsafe"
)

// The following are the memory-protection options; you must specify one of the following values when allocating or
// protecting a page in memory. Protection attributes cannot be assigned to a portion of a page;
// they can only be assigned to a whole page.
//
// More: https://learn.microsoft.com/en-us/windows/win32/memory/memory-protection-constants
const (
	_PAGE_EXECUTE           = 0x10
	_PAGE_EXECUTE_READ      = 0x20
	_PAGE_EXECUTE_READWRITE = 0x40
	_PAGE_EXECUTE_WRITECOPY = 0x80
	_PAGE_NOACCESS          = 0x01
	_PAGE_READONLY          = 0x02
	_PAGE_READWRITE         = 0x04
	_PAGE_WRITECOPY         = 0x08
	_PAGE_TARGETS_INVALID   = 0x40000000
	_PAGE_TARGETS_NO_UPDATE = 0x40000000
	_PAGE_GUARD             = 0x100
	_PAGE_NOCACHE           = 0x200
	_PAGE_WRITECOMBINE      = 0x400
)

var (
	modes = map[ProtectMode]uint32{
		// ProtectModeRead:      _PAGE_EXECUTE_READ,
		ProtectModeRead:      _PAGE_READWRITE,
		ProtectModeReadWrite: _PAGE_EXECUTE_READWRITE,
	}
)

var protector = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

func SetProtect(ptr uintptr, size int, mode ProtectMode) {
	protect, ok := modes[mode]
	if !ok {
		panic(ErrUnsupportedProtectMode)
	}

	var oldProtect uint32

	ret, _, _ := protector.Call(
		ptr,
		uintptr(size),
		uintptr(protect),
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if ret == 0 {
		panic(syscall.GetLastError())
	}
}
