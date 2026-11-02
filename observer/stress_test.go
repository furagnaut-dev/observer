package observer

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolRunsAllJobsBeforeShutdownReturns(t *testing.T) {
	pool := NewWorkerPool(3)
	var done atomic.Int64
	for i := 0; i < 100; i++ {
		pool.Submit(func() { done.Add(1) })
	}

	pool.Shutdown()

	if done.Load() != 100 {
		t.Fatalf("got %d jobs done, want 100", done.Load())
	}
}

func TestStressFourProducersThousandEventsEach(t *testing.T) {
	const producers = 4
	const eventsPerProducer = 1000
	registry, _ := newTestRegistry()
	availability := NewBookAvailability(registry)
	statistics := NewStatistics()
	availability.Subscribe(BookBorrowed, statistics)

	pool := NewWorkerPool(producers)
	for p := 0; p < producers; p++ {
		pool.Submit(func() {
			for i := 0; i < eventsPerProducer; i++ {
				availability.Borrow(p, i, 1, time.Now())
			}
		})
	}
	pool.Shutdown()

	want := producers * eventsPerProducer
	if got := statistics.Count(BookBorrowed); got != want {
		t.Fatalf("got %d events, want %d", got, want)
	}
}

func TestStressSubscribeAndUnsubscribeWhilePublishing(t *testing.T) {
	registry, _ := newTestRegistry()
	availability := NewBookAvailability(registry)
	statistics := NewStatistics()
	availability.Subscribe(BookBorrowed, statistics)

	pool := NewWorkerPool(4)
	for p := 0; p < 3; p++ {
		pool.Submit(func() {
			for i := 0; i < 1000; i++ {
				availability.Borrow(p, i, 1, time.Now())
			}
		})
	}
	pool.Submit(func() {
		for i := 0; i < 1000; i++ {
			unsubscribe := availability.Subscribe(BookBorrowed, &FakeObserverSafe{})
			unsubscribe()
		}
	})
	pool.Shutdown()

	if got := statistics.Count(BookBorrowed); got != 3000 {
		t.Fatalf("got %d events, want 3000", got)
	}
}

type FakeObserverSafe struct {
	count atomic.Int64
}

func (f *FakeObserverSafe) OnEvent(e Event) {
	f.count.Add(1)
}
