package memory

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"unsafe"
)

type PatchFrame[T any] struct {
	old     T
	new     T
	unpatch func()
	once    sync.Once
}

func NewPatch[T any](o, n T) *PatchFrame[T] {
	return &PatchFrame[T]{
		old: o,
		new: n,
	}
}

func (p *PatchFrame[T]) WithUnpatch(fn func()) {
	if p.unpatch != nil {
		return
	}

	p.unpatch = fn
}

func (p *PatchFrame[T]) Unpatch() {
	if p.unpatch == nil {
		return
	}

	p.once.Do(p.unpatch)
}

func (p *PatchFrame[T]) OldImpl() T {
	return p.old
}

func (p *PatchFrame[T]) NewImpl() T {
	return p.new
}

func Patch(header, footer, eof uintptr, newHeader, newFooter []byte) {
	var headerUint32 = *(**uint32)(unsafe.Pointer(&header))

	SetProtect(header, int(eof-header), ProtectModeReadWrite)
	defer SetProtect(header, int(eof-header), ProtectModeRead)

	atomic.StoreUint32(headerUint32, 0xFEEB)
	defer atomic.StoreUint32(headerUint32, binary.LittleEndian.Uint32(newHeader))

	Copy(header+4, uintptr(unsafe.Pointer(&newHeader[0])), len(newHeader)-4)
	Copy(footer, uintptr(unsafe.Pointer(&newFooter[0])), len(newFooter))
}
