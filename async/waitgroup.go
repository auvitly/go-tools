package async

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	_ch = make(chan struct{})
)

func init() {
	close(_ch)
}

// WaitGroup adapter over sync.WaitGroup that allows you to complete the wait by context.
type WaitGroup struct {
	mu   sync.Mutex
	done atomic.Value
	sync.WaitGroup
	goroutine bool
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

// WaitDone returns a channel that is closed when the wait is complete.
func (w *WaitGroup) WaitDone() <-chan struct{} {
	if !w.goroutine {
		w.waitGoroutine()
	}

	d := w.done.Load()
	if d != nil {
		return d.(chan struct{})
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	d = w.done.Load()
	if d == nil {
		d = make(chan struct{})
		w.done.Store(d)
	}

	return d.(chan struct{})
}

func (w *WaitGroup) waitGoroutine() {
	w.mu.Lock()
	defer w.mu.Unlock()

	go func() {
		w.Wait()

		w.mu.Lock()
		defer w.mu.Unlock()

		d := w.done.Load()
		if d == nil {
			w.done.Store(_ch)
		} else {
			close(d.(chan struct{}))
		}
	}()

	w.goroutine = true
}
