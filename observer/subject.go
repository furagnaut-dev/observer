package observer

import (
	"sync"
	"time"
)

type Registry interface {
	Subscribe(EventType, Observer) (unsubscribe func())
	notifySubscribers(Event)
}

type subscription struct {
	id     int
	reader Observer
}

type ObserverRegistry struct {
	mu      sync.RWMutex
	Readers map[EventType][]subscription
	nextID  int
}

func NewObserverRegistry() *ObserverRegistry {
	return &ObserverRegistry{
		Readers: make(map[EventType][]subscription),
	}
}

func (o *ObserverRegistry) Subscribe(eventType EventType, observer Observer) (unsubscribe func()) {
	o.mu.Lock()
	defer o.mu.Unlock()
	id := o.nextID
	o.nextID++
	o.Readers[eventType] = append(o.Readers[eventType], subscription{id: id, reader: observer})
	return func() {
		o.mu.Lock()
		kept := []subscription{}
		for _, s := range o.Readers[eventType] {
			if s.id != id {
				kept = append(kept, s)
			}
		}
		o.Readers[eventType] = kept
		defer o.mu.Unlock()
	}
}

func (o *ObserverRegistry) notifySubscribers(event Event) {
	o.mu.RLock()
	snap := o.Readers[event.Type()]
	o.mu.RUnlock()

	for _, s := range snap {
		s.reader.OnEvent(event)
	}
}

type BookAvailability struct {
	*ObserverRegistry
}

func NewBookAvailability(registry *ObserverRegistry) *BookAvailability {
	return &BookAvailability{ObserverRegistry: registry}
}

func (b *BookAvailability) MarkAvailable(sectionID int, bookID int) {
	b.notifySubscribers(NewEvent(BookAvailable, sectionID, bookID,
		0, 0, time.Time{}))
}

func (b *BookAvailability) Borrow(sectionID int, bookID int, readerID int, due time.Time) {
	b.notifySubscribers(NewEvent(BookBorrowed, sectionID, bookID, readerID,
		0, due))
}

func (b *BookAvailability) MarkOverdue(sectionID int, bookID int, readerID int, due time.Time) {
	b.notifySubscribers(NewEvent(BookOverdue, sectionID, bookID, readerID,
		0, due))
}

type SectionSales struct {
	*ObserverRegistry
}

func NewSectionSales(registry *ObserverRegistry) *SectionSales {
	return &SectionSales{ObserverRegistry: registry}
}

func (s *SectionSales) PriceDrop(sectionID int, newPrice int) {
	s.notifySubscribers(NewEvent(PriceDropped, sectionID, 0, 0, newPrice, time.Time{}))
}

func (s *SectionSales) SalesEnd(sectionID int) {
	s.notifySubscribers(NewEvent(SaleEnded, sectionID, 0, 0, 0, time.Time{}))
}
