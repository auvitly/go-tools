//go:build darwin
// +build darwin

package memory

import "fmt"

const (
	_VMProtectRead uint32 = 1 << iota
	_VMProtectWrite
	_VMProtectExecute
	_VMProtectNoChange
	_VMProtectCopy
)

var (
	// ErrMachVMProtect - .
	ErrMachVMProtect = "memory protection change failed"
	modes            = map[ProtectMode]uint32{
		ProtectModeRead:      _VMProtectRead | _VMProtectExecute | _VMProtectCopy,
		ProtectModeReadWrite: _VMProtectRead | _VMProtectExecute | _VMProtectCopy | _VMProtectWrite,
	}
)

//go:noescape
func taskSelfTrap() (ret uint32)

//go:noescape
func vmProtect(targetTask uint32, address uintptr, size int, setMaximum, newProtection uint32) (ret uint32)

func SetProtect(ptr uintptr, size int, mode ProtectMode) {
	protect, ok := modes[mode]
	if !ok {
		panic(ErrUnsupportedProtectMode)
	}

	if ret := vmProtect(taskSelfTrap(), ptr, size, 0, protect); ret != 0 {
		panic(fmt.Sprintf("vmProtect: ret = %v", ret))
	}
}
