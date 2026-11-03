package observer

import "time"

type EventType string

const (
	BookAvailable  EventType = "BOOK_AVAILABLE"
	BookBorrowed   EventType = "BOOK_BORROWED"
	BookOverdue    EventType = "BOOK_OVERDUE"
	PriceDropped   EventType = "PRICE_DROPPED"
	SaleEnded      EventType = "SALE_ENDED"
	SectionBusy    EventType = "SECTION_BUSY"
	InspectionDone EventType = "INSPECTION_DONE"
)

type Event struct {
	eventType  EventType
	sectionID  int
	bookID     int
	readerID   int
	price      int
	due        time.Time
	occurredAt time.Time
}

func NewEvent(eventType EventType, sectionID, bookID, readerID, price int, due time.Time) Event {
	return Event{
		eventType:  eventType,
		sectionID:  sectionID,
		bookID:     bookID,
		readerID:   readerID,
		price:      price,
		due:        due,
		occurredAt: time.Now(),
	}
}

func (e Event) Type() EventType       { return e.eventType }
func (e Event) SectionID() int        { return e.sectionID }
func (e Event) BookID() int           { return e.bookID }
func (e Event) ReaderID() int         { return e.readerID }
func (e Event) Price() int            { return e.price }
func (e Event) Due() time.Time        { return e.due }
func (e Event) OccurredAt() time.Time { return e.occurredAt }
