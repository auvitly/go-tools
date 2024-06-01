package async_test

import (
	"context"
	"github.com/auvitly/go-tools/async"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestAsyncWaitGroup_Wait(t *testing.T) {
	var (
		d      int
		result = 10
		mu     sync.Mutex
		use    = func() {
			mu.Lock()
			defer mu.Unlock()
			d++
		}
		wg async.WaitGroup
	)

	for i := 0; i < result; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(time.Second)

			use()
		}()
	}

	wg.Wait()

	require.Equal(t, d, result)
}

func TestSyncWaitGroup_WaitContext(t *testing.T) {
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
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
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

	require.Equal(t, d, result)
}

func TestAsyncWaitGroup_WaitContextDone(t *testing.T) {
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
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	)

	t.Cleanup(func() {
		cancel()
	})

	for i := 0; i < result; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(2 * time.Second)

			use()
		}()
	}

	wg.WaitContext(ctx)

	require.Equal(t, d, 0)
}

func TestWaitGroup_WaitContextEdge(t *testing.T) {
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
		ctx, cancel = context.WithTimeout(context.Background(), time.Microsecond)
	)

	t.Cleanup(func() {
		cancel()
	})

	for i := 0; i < result; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(time.Microsecond)

			use()
		}()
	}

	wg.WaitContext(ctx)

	require.GreaterOrEqual(t, d, 0)
}

func benchmarkAsyncWaitGroupWait(b *testing.B, localWork int) {
	var wg async.WaitGroup
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			wg.Wait()
			for i := 0; i < localWork; i++ {
				foo *= 2
				foo /= 2
			}
		}
		_ = foo
	})
}

func BenchmarkAsyncWaitGroupWait(b *testing.B) {
	benchmarkAsyncWaitGroupWait(b, 0)
}

func BenchmarkAsyncWaitGroupWaitWork(b *testing.B) {
	benchmarkAsyncWaitGroupWait(b, 100)
}

func BenchmarkAsyncWaitGroupActuallyWait(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg async.WaitGroup
			wg.Add(1)
			go func() {
				wg.Done()
			}()
			wg.Wait()
		}
	})
}

func benchmarkSyncWaitGroupWait(b *testing.B, localWork int) {
	var wg sync.WaitGroup
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			wg.Wait()
			for i := 0; i < localWork; i++ {
				foo *= 2
				foo /= 2
			}
		}
		_ = foo
	})
}

func BenchmarkSyncWaitGroupWait(b *testing.B) {
	benchmarkSyncWaitGroupWait(b, 0)
}

func BenchmarkSyncWaitGroupWaitWork(b *testing.B) {
	benchmarkSyncWaitGroupWait(b, 100)
}

func BenchmarkSyncWaitGroupActuallyWait(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				wg.Done()
			}()
			wg.Wait()
		}
	})
}
