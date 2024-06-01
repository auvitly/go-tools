package async_test

import (
	"context"
	"github.com/auvitly/go-tools/async"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestWaitGroup_WaitContext(t *testing.T) {
	var (
		d      int
		result = 1000
		mu     sync.Mutex
		use    = func() {
			mu.Lock()
			defer mu.Unlock()
			d++
		}
		wg  async.WaitGroup
		ctx = context.Background()
	)

	for i := 0; i < result; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(time.Second)

			use()
		}()
	}

	wg.WaitContext(ctx)

	require.Equal(t, d, result)
}

func TestWaitGroup_WaitContext_Edge(t *testing.T) {
	var (
		d      int
		result = 1000
		mu     sync.Mutex
		use    = func() {
			mu.Lock()
			defer mu.Unlock()
			d++
		}
		wg          async.WaitGroup
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	)

	t.Cleanup(func() {
		cancel()
	})

	for i := 0; i < result; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(time.Second)

			use()
		}()
	}

	wg.WaitContext(ctx)

	require.NotEqual(t, d, result)
}

func TestWaitGroup_WaitContext_Done(t *testing.T) {
	var (
		d      int
		result = 10
		mu     sync.Mutex
		use    = func() {
			mu.Lock()
			defer mu.Unlock()
			d++
		}
		wg          async.WaitGroup
		ctx, cancel = context.WithTimeout(context.Background(), 0)
	)

	t.Cleanup(func() {
		cancel()
	})

	for i := 0; i < result; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(time.Second)

			use()
		}()
	}

	wg.WaitContext(ctx)

	require.Equal(t, d, 0)
}

func TestWaitGroup_WaitCh(t *testing.T) {
	var (
		d      int
		result = 1000
		use    = func() {}
		wg     async.WaitGroup
	)

	for i := 0; i < result; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(time.Second)

			use()
		}()
	}

	<-wg.WaitDone()

	require.NotEqual(t, d, result)
}
