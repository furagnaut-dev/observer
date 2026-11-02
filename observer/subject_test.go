package observer

import (
	"testing"
	"time"
)

type FakeObserver struct {
	count int
}

func (f *FakeObserver) OnEvent(e Event) {
	f.count++
}

func newTestRegistry() (*ObserverRegistry, *FakeObserver) {
	registry := NewObserverRegistry()
	fake := &FakeObserver{}
	return registry, fake
}

func TestSubscribe(t *testing.T) {
	registry, fake := newTestRegistry()
	registry.Subscribe(BookAvailable, fake)

	registry.notifySubscribers(NewEvent(BookAvailable, 1, 1, 1, 1000, time.Now()))

	if fake.count != 1 {
		t.Fatalf("got %d calls, want 1", fake.count)
	}
}

func TestUnsubscribeStopDelivery(t *testing.T) {
	registry, fake := newTestRegistry()
	unsubscribe := registry.Subscribe(BookAvailable, fake)
	registry.notifySubscribers(NewEvent(BookAvailable, 1, 1, 1, 1000, time.Now()))
	if fake.count != 1 {
		t.Fatalf("got %d calls, want 1", fake.count)
	}

	unsubscribe()
	registry.notifySubscribers(NewEvent(BookAvailable, 1, 1, 1, 1000, time.Now()))
	if fake.count != 1 {
		t.Fatalf("got %d calls, want 0", fake.count)
	}
}

func TestBookAvailabilityNegative(t *testing.T) {
	registry, fake := newTestRegistry()
	availability := NewBookAvailability(registry)
	availability.Subscribe(BookAvailable, fake)
	availability.Borrow(1, 1, 1, time.Now())

	if fake.count != 0 {
		t.Fatalf("got %d, want 0", fake.count)
	}
}

func TestBookAvailabilityPositive(t *testing.T) {
	registry, fake := newTestRegistry()
	availability := NewBookAvailability(registry)
	availability.Subscribe(BookAvailable, fake)
	availability.MarkAvailable(1, 1)

	if fake.count != 1 {
		t.Fatalf("got %d, want 1", fake.count)
	}
}

func TestBookAvailabilityBorrowedTwo(t *testing.T) {
	registry, fake := newTestRegistry()
	availability := NewBookAvailability(registry)
	availability.Subscribe(BookBorrowed, fake)
	availability.Borrow(1, 1, 1, time.Now())
	availability.Borrow(2, 2, 2, time.Now())

	if fake.count != 2 {
		t.Fatalf("got %d, want 2", fake.count)
	}
}

func TestStatisticsCountBorrowsPerType(t *testing.T) {
	registry, _ := newTestRegistry()
	availability := NewBookAvailability(registry)
	statistics := NewStatistics()
	availability.Subscribe(BookBorrowed, statistics)
	availability.Borrow(1, 1, 1, time.Now())
	availability.Borrow(2, 2, 2, time.Now())
	borrowedCount := statistics.Count(BookBorrowed)
	availableCount := statistics.Count(BookAvailable)

	if borrowedCount != 2 || availableCount != 0 {
		t.Fatalf("got %d and %d, want 2 and 0:", borrowedCount, availableCount)
	}
}
