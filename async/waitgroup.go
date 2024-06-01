package async

import (
	"context"
	"errors"
	"github.com/auvitly/go-tools/stderrs"
	"sync"
	"sync/atomic"
)

// WaitGroup adapter over sync.WaitGroup that allows you to complete the wait by context.
type WaitGroup struct {
	mu   sync.Mutex
	done atomic.Value
	once sync.Once
	sync.WaitGroup
}

// WaitContext blocks until the WaitGroup counter is zero or context done.
func (w *WaitGroup) WaitContext(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return stderrs.DeadlineExceeded.EmbedErrors(ctx.Err())
			}

			return stderrs.Canceled.EmbedErrors(ctx.Err())
		case <-w.WaitDone():
			return nil
		}
	}
}

// WaitDone returns a channel that is closed when the wait is complete.
func (w *WaitGroup) WaitDone() <-chan struct{} {
	w.once.Do(w.waitGoroutines)

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

func (w *WaitGroup) waitGoroutines() {
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
}
