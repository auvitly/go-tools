package function

import (
	"encoding/binary"
	"github.com/auvitly/go-tools/nuclear/memory"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Patch[T any] struct {
	Old     T
	New     T
	unpatch func()
	once    sync.Once
}

func NewPatch[T any](o, n T) *Patch[T] {
	return &Patch[T]{
		Old: o,
		New: n,
	}
}

func (p *Patch[T]) WithUnpatch(fn func()) *Patch[T] {
	if p.unpatch != nil {
		return p
	}

	p.unpatch = fn

	return p
}

func (p *Patch[T]) Unpatch() {
	if p.unpatch == nil {
		return
	}

	p.once.Do(p.unpatch)
}

func doPatch(header, footer, eof uintptr, newHeader, newFooter []byte) {
	var headerUint32 = *(**uint32)(unsafe.Pointer(&header))

	memory.SetProtect(header, int(eof-header), memory.ProtectModeReadWrite)
	defer memory.SetProtect(header, int(eof-header), memory.ProtectModeRead)

	atomic.StoreUint32(headerUint32, 0xFEEB)
	defer atomic.StoreUint32(headerUint32, binary.LittleEndian.Uint32(newHeader))

	memory.Copy(header+4, uintptr(unsafe.Pointer(&newHeader[4])), len(newHeader)-4)
	memory.Copy(footer, uintptr(unsafe.Pointer(&newFooter[0])), len(newFooter))
}
