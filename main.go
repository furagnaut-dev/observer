package main

import (
	"fmt"
	"os"
	"time"

	"observer/library"
	"observer/observer"
)

const (
	workers    = 4
	producers  = 4
	eventsEach = 1000
	busyAfter  = 500
)

func buildTree() *library.Library {
	return &library.Library{Name: "Hamburger", Sections: []*library.Section{
		{ID: 1, Name: "Philosophy", Shelves: []*library.Shelf{
			{Label: "A", Items: []library.Node{
				&library.Book{ID: 1, Title: "Republic", PriceCents: 2000, Pages: 300},
				&library.EBook{ID: 2, Title: "Ethics", PriceCents: 500, SizeMB: 4},
			}},
		}},
		{ID: 2, Name: "Science", Shelves: []*library.Shelf{
			{Label: "B", Items: []library.Node{
				&library.Journal{ID: 3, Title: "Nature", PriceCents: 1500, Issue: 12},
				&library.AudioBook{ID: 4, Title: "Cosmos", PriceCents: 2500, Minutes: 600},
			}},
		}},
	}}
}

func main() {
	tree := buildTree()

	availability := observer.NewBookAvailability(observer.NewObserverRegistry())
	sales := observer.NewSectionSales(observer.NewObserverRegistry())
	aggregator := observer.NewSectionAggregator(observer.NewObserverRegistry(), busyAfter)
	inspection := observer.NewInspectionService(observer.NewObserverRegistry(), tree)

	statistics := observer.NewStatistics()
	alert := observer.NewOverdueAlert()
	panel := observer.NewLibrarianPanel(10)
	tracker := observer.NewPriceDropTracker()

	availability.Subscribe(observer.BookBorrowed, statistics)
	availability.Subscribe(observer.BookBorrowed, aggregator)
	availability.Subscribe(observer.BookOverdue, alert)
	availability.Subscribe(observer.BookOverdue, inspection)
	aggregator.Subscribe(observer.SectionBusy, panel)
	inspection.Subscribe(observer.InspectionDone, panel)
	sales.Subscribe(observer.PriceDropped, tracker)
	sales.Subscribe(observer.SaleEnded, tracker)

	failed := 0
	expect := func(label string, got, want any) {
		mark := "ok  "
		if got != want {
			mark = "FAIL"
			failed++
		}
		fmt.Printf("  [%s] %s: %v (want %v)\n", mark, label, got, want)
	}

	fmt.Printf("%d producers x %d borrows on %d workers\n", producers, eventsEach, workers)
	pool := observer.NewWorkerPool(workers)
	for p := 0; p < producers; p++ {
		section := p%2 + 1
		pool.Submit(func() {
			for i := 0; i < eventsEach; i++ {
				availability.Borrow(section, i, p, time.Now())
			}
		})
	}
	pool.Shutdown()
	expect("borrows counted", statistics.Count(observer.BookBorrowed), producers*eventsEach)

	fmt.Println("section 1 goes on sale:")
	sales.PriceDrop(1, 1500)
	price, _ := tracker.Price(1)
	expect("tracked sale price", price, 1500)

	fmt.Println("a book in section 1 is overdue:")
	availability.MarkOverdue(1, 1, 7, time.Now().Add(-48*time.Hour))
	feed := panel.Feed()
	expect("inspected value (section 1 only, the library holds 6500)", feed[len(feed)-1].Price(), 2500)
	for _, e := range feed {
		fmt.Printf("    panel: %-15s section=%d value=%d\n", e.Type(), e.SectionID(), e.Price())
	}

	fmt.Println("visitors:")
	value := library.NewInventoryValue()
	tree.Accept(value)
	expect("inventory value, whole library", value.Total(), 6500)
	export := library.NewJSONExport()
	tree.Sections[1].Accept(export)
	fmt.Println("  section 2 as JSON:", export.String())

	if failed > 0 {
		fmt.Printf("%d check(s) FAILED\n", failed)
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
