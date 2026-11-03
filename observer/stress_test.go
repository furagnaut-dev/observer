package observer

import (
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	registry, _ := newTestRegistry()
	statistics := NewStatistics()
	availability := NewBookAvailability(registry)
	availability.Subscribe(BookBorrowed, statistics)
	workers := 4
	eventsPerWorker := 1000
	pool := NewWorkerPool(workers)
	for p := 0; p < workers; p++ {
		pool.Submit(func() {
			for i := 0; i < eventsPerWorker; i++ {
				availability.Borrow(p, i, 1, time.Now())
			}
		})
	}
	pool.Shutdown()
	want := workers * eventsPerWorker
	if got := statistics.Count(BookBorrowed); got != want {
		t.Fatalf("got %d events, want %d", got, want)
	}
}
