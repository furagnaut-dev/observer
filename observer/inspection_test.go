package observer

import (
	"testing"
	"time"

	"observer/library"
)

func inspectionTree() *library.Library {
	return &library.Library{Name: "L", Sections: []*library.Section{
		{ID: 1, Name: "S1", Shelves: []*library.Shelf{
			{Label: "A", Items: []library.Node{
				&library.Book{ID: 1, Title: "B", PriceCents: 700, Pages: 1},
				&library.EBook{ID: 2, Title: "E", PriceCents: 300, SizeMB: 1},
			}},
		}},
		{ID: 2, Name: "S2", Shelves: []*library.Shelf{
			{Label: "B", Items: []library.Node{
				&library.Journal{ID: 3, Title: "J", PriceCents: 5000, Issue: 1},
			}},
		}},
	}}
}

func TestInspectionServiceReportsOnlyOverdueSection(t *testing.T) {
	registry, _ := newTestRegistry()
	availability := NewBookAvailability(registry)
	inspection := NewInspectionService(NewObserverRegistry(), inspectionTree())
	panel := NewLibrarianPanel(5)
	availability.Subscribe(BookOverdue, inspection)
	inspection.Subscribe(InspectionDone, panel)

	availability.MarkOverdue(1, 1, 7, time.Now())

	feed := panel.Feed()
	if len(feed) != 1 || feed[0].Type() != InspectionDone || feed[0].SectionID() != 1 || feed[0].Price() != 1000 {
		t.Fatalf("got %d events, want one INSPECTION_DONE for section 1 with value 1000", len(feed))
	}
}

func TestInspectionServiceIgnoresUnknownSection(t *testing.T) {
	registry, _ := newTestRegistry()
	availability := NewBookAvailability(registry)
	inspection := NewInspectionService(NewObserverRegistry(), inspectionTree())
	panel := NewLibrarianPanel(5)
	availability.Subscribe(BookOverdue, inspection)
	inspection.Subscribe(InspectionDone, panel)

	availability.MarkOverdue(99, 1, 7, time.Now())

	if len(panel.Feed()) != 0 {
		t.Fatalf("got %d events, want 0", len(panel.Feed()))
	}
}
