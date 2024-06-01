package main

import (
	"github.com/auvitly/go-tools/async"
	"log/slog"
	"time"
)

func goroutine(wg *async.WaitGroup) {
	defer wg.Done()

	time.Sleep(5 * time.Second)
}

func main() {
	var (
		wg     async.WaitGroup
		ticker = time.NewTicker(time.Second)
	)

	wg.Add(1)

	go goroutine(&wg)

	for {
		select {
		case <-wg.WaitDone():
			slog.Info("all goroutines done")

			return
		case <-ticker.C:
			slog.Info("1 second has passed")
		}
	}
}
