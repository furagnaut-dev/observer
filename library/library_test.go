package library

import "testing"

func sampleLibrary() *Library {
	return &Library{
		Name: "City",
		Sections: []*Section{
			{ID: 1, Name: "Philosophy", Shelves: []*Shelf{
				{Label: "A", Items: []Node{
					&Book{ID: 1, Title: "Republic", PriceCents: 2000, Pages: 300},
					&EBook{ID: 2, Title: "Ethics", PriceCents: 500, SizeMB: 4},
				}},
				{Label: "B", Items: []Node{
					&Journal{ID: 3, Title: "Mind", PriceCents: 1500, Issue: 12},
				}},
			}},
			{ID: 2, Name: "Science", Shelves: []*Shelf{
				{Label: "C", Items: []Node{
					&AudioBook{ID: 4, Title: "Cosmos", PriceCents: 2500, Minutes: 600},
					&Book{ID: 5, Title: "Origin", PriceCents: 3000, Pages: 500},
				}},
			}},
		},
	}
}

func itemID(n Node) int {
	switch item := n.(type) {
	case *Book:
		return item.ID
	case *Journal:
		return item.ID
	case *EBook:
		return item.ID
	case *AudioBook:
		return item.ID
	}
	return -1
}

func TestInventoryValueWholeTree(t *testing.T) {
	value := NewInventoryValue()
	sampleLibrary().Accept(value)

	if value.Total() != 9500 {
		t.Fatalf("got %d, want 9500", value.Total())
	}
}

func TestInventoryValueSubtreeOnly(t *testing.T) {
	value := NewInventoryValue()
	sampleLibrary().Sections[0].Accept(value)

	if value.Total() != 4000 {
		t.Fatalf("got %d, want 4000", value.Total())
	}
}

func TestInventoryValueSingleShelf(t *testing.T) {
	value := NewInventoryValue()
	sampleLibrary().Sections[0].Shelves[1].Accept(value)

	if value.Total() != 1500 {
		t.Fatalf("got %d, want 1500", value.Total())
	}
}

func TestInventoryValueEmptyLibrary(t *testing.T) {
	value := NewInventoryValue()
	(&Library{Name: "Empty"}).Accept(value)

	if value.Total() != 0 {
		t.Fatalf("got %d, want 0", value.Total())
	}
}

func TestItemCensusCountsEveryKind(t *testing.T) {
	census := NewItemCensus()
	sampleLibrary().Accept(census)

	if census.Count(KindBook) != 2 || census.Count(KindEBook) != 1 ||
		census.Count(KindJournal) != 1 || census.Count(KindAudioBook) != 1 {
		t.Fatalf("wrong counts: books %d, ebooks %d, journals %d, audiobooks %d",
			census.Count(KindBook), census.Count(KindEBook),
			census.Count(KindJournal), census.Count(KindAudioBook))
	}
	if census.Shelves() != 3 {
		t.Fatalf("got %d shelves, want 3", census.Shelves())
	}
}

func TestCheapItemsReturnsMatchesInTreeOrder(t *testing.T) {
	cheap := NewCheapItems(2000)
	sampleLibrary().Accept(cheap)

	found := cheap.Found()
	if len(found) != 3 || itemID(found[0]) != 1 || itemID(found[1]) != 2 || itemID(found[2]) != 3 {
		t.Fatalf("got %d items, want ids 1, 2, 3 in order", len(found))
	}
}

func TestCheapItemsNothingMatches(t *testing.T) {
	cheap := NewCheapItems(100)
	sampleLibrary().Accept(cheap)

	if len(cheap.Found()) != 0 {
		t.Fatalf("got %d items, want 0", len(cheap.Found()))
	}
}

func TestJSONExportNesting(t *testing.T) {
	tree := &Library{Name: "L", Sections: []*Section{
		{ID: 1, Name: "S", Shelves: []*Shelf{
			{Label: "A", Items: []Node{
				&Book{ID: 1, Title: "B", PriceCents: 100, Pages: 10},
				&EBook{ID: 2, Title: "E", PriceCents: 50, SizeMB: 3},
			}},
		}},
	}}
	want := `{"type":"library","name":"L","sections":[{"type":"section","id":1,"name":"S","shelves":[{"type":"shelf","label":"A","items":[{"type":"book","id":1,"title":"B","priceCents":100,"pages":10},{"type":"ebook","id":2,"title":"E","priceCents":50,"sizeMB":3}]}]}]}`

	export := NewJSONExport()
	tree.Accept(export)

	if export.String() != want {
		t.Fatalf("got\n%s\nwant\n%s", export.String(), want)
	}
}

func TestJSONExportEscapesText(t *testing.T) {
	export := NewJSONExport()
	(&Book{ID: 1, Title: "a\"b\\c\nd", PriceCents: 1, Pages: 2}).Accept(export)
	want := `{"type":"book","id":1,"title":"a\"b\\c\nd","priceCents":1,"pages":2}`

	if export.String() != want {
		t.Fatalf("got\n%s\nwant\n%s", export.String(), want)
	}
}

func TestJSONExportEmptyLibrary(t *testing.T) {
	export := NewJSONExport()
	(&Library{Name: "E"}).Accept(export)
	want := `{"type":"library","name":"E","sections":[]}`

	if export.String() != want {
		t.Fatalf("got %s, want %s", export.String(), want)
	}
}
