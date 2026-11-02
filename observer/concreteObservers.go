package observer

import (
	"sync"
	"time"
)

type Statistics struct {
	mu     sync.Mutex
	counts map[EventType]int
}

func NewStatistics() *Statistics {
	return &Statistics{
		counts: map[EventType]int{},
	}
}

func (s *Statistics) OnEvent(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[e.Type()]++
}

func (s *Statistics) Count(eventType EventType) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counts[eventType]
}

type PriceDropTracker struct {
	mu      sync.Mutex
	pricing map[int]int
}

func NewPriceDropTracker() *PriceDropTracker {
	return &PriceDropTracker{
		pricing: map[int]int{},
	}
}

func (p *PriceDropTracker) OnEvent(e Event) {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch e.Type() {
	case PriceDropped:
		p.pricing[e.SectionID()] = e.Price()
	case SaleEnded:
		delete(p.pricing, e.SectionID())
	}
}

func (p *PriceDropTracker) Price(sectionID int) (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	price, ok := p.pricing[sectionID]
	return price, ok
}

type OverdueAlert struct {
	mu      sync.Mutex
	overdue []Event
}

func NewOverdueAlert() *OverdueAlert {
	return &OverdueAlert{
		overdue: make([]Event, 0),
	}
}

func (o *OverdueAlert) OnEvent(e Event) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.overdue = append(o.overdue, e)
}

func (o *OverdueAlert) Overdue() []Event {
	o.mu.Lock()
	defer o.mu.Unlock()
	result := make([]Event, 0, len(o.overdue))
	result = append(result, o.overdue...)
	return result
}

type LibrarianPanel struct {
	mu    sync.Mutex
	feed  []Event
	limit int
}

func NewLibrarianPanel(limit int) *LibrarianPanel {
	return &LibrarianPanel{
		feed:  make([]Event, 0),
		limit: limit,
	}
}

func (l *LibrarianPanel) OnEvent(e Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.feed = append(l.feed, e)
	if len(l.feed) > l.limit {
		l.feed = l.feed[len(l.feed)-l.limit:]
	}
}

func (l *LibrarianPanel) Feed() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	result := make([]Event, 0, len(l.feed))
	result = append(result, l.feed...)
	return result
}

type SectionAggregator struct {
	*ObserverRegistry
	mu        sync.Mutex
	borrows   map[int]int
	threshold int
}

func NewSectionAggregator(registry *ObserverRegistry, threshold int) *SectionAggregator {
	return &SectionAggregator{
		ObserverRegistry: registry,
		borrows:          make(map[int]int),
		threshold:        threshold,
	}
}

func (a *SectionAggregator) OnEvent(e Event) {
	a.mu.Lock()
	a.borrows[e.SectionID()]++
	count := a.borrows[e.SectionID()]
	a.mu.Unlock()
	if count == a.threshold {
		a.notifySubscribers(NewEvent(SectionBusy, e.SectionID(), 0, 0, count, time.Time{}))
	}
}
