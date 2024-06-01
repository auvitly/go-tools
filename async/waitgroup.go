package async

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

// WaitGroup adapter over sync.WaitGroup that allows you to complete the wait by context.
type WaitGroup struct {
	sync.WaitGroup
}

// _WaitGroup - internal implementation.
type _WaitGroup struct {
	_     struct{}
	state atomic.Uint64
	_     uint32
}

// WaitContext blocks until the WaitGroup counter is zero or context done.
func (w *WaitGroup) WaitContext(ctx context.Context) {
	wg := (*_WaitGroup)(unsafe.Pointer(w))

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if wg.state.Load() == 0 {
				return
			}
		}
	}
}
