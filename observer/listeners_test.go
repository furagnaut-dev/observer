package observer

import (
	"testing"
	"time"
)

type selfRemovingObserver struct {
	count       int
	unsubscribe func()
}

func (s *selfRemovingObserver) OnEvent(e Event) {
	s.count++
	s.unsubscribe()
}

func TestUnsubscribeInsideHandler(t *testing.T) {
	registry, _ := newTestRegistry()
	observer := &selfRemovingObserver{}
	observer.unsubscribe = registry.Subscribe(BookAvailable, observer)

	registry.notifySubscribers(NewEvent(BookAvailable, 1, 1, 0, 0, time.Time{}))
	registry.notifySubscribers(NewEvent(BookAvailable, 1, 1, 0, 0, time.Time{}))

	if observer.count != 1 {
		t.Fatalf("got %d calls, want 1", observer.count)
	}
}

func TestOverdueAlertCollectsOverdue(t *testing.T) {
	registry, _ := newTestRegistry()
	availability := NewBookAvailability(registry)
	alert := NewOverdueAlert()
	availability.Subscribe(BookOverdue, alert)
	availability.MarkOverdue(1, 1, 1, time.Now())
	availability.MarkOverdue(1, 2, 2, time.Now())

	if len(alert.Overdue()) != 2 {
		t.Fatalf("got %d, want 2", len(alert.Overdue()))
	}
}

func TestPriceDropTrackerSaleEnded(t *testing.T) {
	registry, _ := newTestRegistry()
	sales := NewSectionSales(registry)
	tracker := NewPriceDropTracker()
	sales.Subscribe(PriceDropped, tracker)
	sales.Subscribe(SaleEnded, tracker)

	sales.PriceDrop(1, 500)
	price, onSale := tracker.Price(1)
	if price != 500 || !onSale {
		t.Fatalf("got %d, %v, want 500, true", price, onSale)
	}

	sales.SalesEnd(1)
	price, onSale = tracker.Price(1)
	if price != 0 || onSale {
		t.Fatalf("got %d, %v, want 0, false", price, onSale)
	}
}

func TestLibrarianPanelKeepsNewest(t *testing.T) {
	registry, _ := newTestRegistry()
	availability := NewBookAvailability(registry)
	panel := NewLibrarianPanel(2)
	availability.Subscribe(BookAvailable, panel)
	availability.MarkAvailable(1, 1)
	availability.MarkAvailable(1, 2)
	availability.MarkAvailable(1, 3)

	feed := panel.Feed()
	if len(feed) != 2 || feed[0].BookID() != 2 || feed[1].BookID() != 3 {
		t.Fatalf("got %d events, want the 2 newest (books 2 and 3)", len(feed))
	}
}

func TestSectionAggregatorPublishesBusyOnce(t *testing.T) {
	sourceRegistry, fake := newTestRegistry()
	availability := NewBookAvailability(sourceRegistry)
	aggregator := NewSectionAggregator(NewObserverRegistry(), 3)
	availability.Subscribe(BookBorrowed, aggregator)
	aggregator.Subscribe(SectionBusy, fake)

	for book := 1; book <= 4; book++ {
		availability.Borrow(7, book, 1, time.Now())
	}

	if fake.count != 1 {
		t.Fatalf("got %d SectionBusy events, want 1", fake.count)
	}
}

func TestSectionAggregatorCountsSectionsSeparately(t *testing.T) {
	sourceRegistry, fake := newTestRegistry()
	availability := NewBookAvailability(sourceRegistry)
	aggregator := NewSectionAggregator(NewObserverRegistry(), 2)
	availability.Subscribe(BookBorrowed, aggregator)
	aggregator.Subscribe(SectionBusy, fake)

	availability.Borrow(1, 1, 1, time.Now())
	availability.Borrow(2, 2, 1, time.Now())

	if fake.count != 0 {
		t.Fatalf("got %d, want 0: one borrow per section is below the threshold", fake.count)
	}
}
