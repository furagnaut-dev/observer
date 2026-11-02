package library

import (
	"fmt"
	"strconv"
	"strings"
)

type InventoryValue struct {
	total int
}

func NewInventoryValue() *InventoryValue {
	return &InventoryValue{}
}

func (i *InventoryValue) Total() int { return i.total }

func (i *InventoryValue) VisitLibrary(l *Library) { acceptAll(i, l.Sections) }
func (i *InventoryValue) VisitSection(s *Section) { acceptAll(i, s.Shelves) }
func (i *InventoryValue) VisitShelf(s *Shelf)     { acceptAll(i, s.Items) }
func (i *InventoryValue) VisitBook(b *Book)       { i.total += b.PriceCents }
func (i *InventoryValue) VisitJournal(j *Journal) { i.total += j.PriceCents }
func (i *InventoryValue) VisitEBook(e *EBook)     { i.total += e.PriceCents }
func (i *InventoryValue) VisitAudioBook(a *AudioBook) {
	i.total += a.PriceCents
}

type ItemCensus struct {
	kinds   map[string]int
	shelves int
}

func NewItemCensus() *ItemCensus {
	return &ItemCensus{kinds: make(map[string]int)}
}

func (c *ItemCensus) Count(kind string) int { return c.kinds[kind] }
func (c *ItemCensus) Shelves() int          { return c.shelves }

func (c *ItemCensus) VisitLibrary(l *Library) { acceptAll(c, l.Sections) }
func (c *ItemCensus) VisitSection(s *Section) { acceptAll(c, s.Shelves) }
func (c *ItemCensus) VisitShelf(s *Shelf) {
	c.shelves++
	acceptAll(c, s.Items)
}
func (c *ItemCensus) VisitBook(b *Book)           { c.kinds[KindBook]++ }
func (c *ItemCensus) VisitJournal(j *Journal)     { c.kinds[KindJournal]++ }
func (c *ItemCensus) VisitEBook(e *EBook)         { c.kinds[KindEBook]++ }
func (c *ItemCensus) VisitAudioBook(a *AudioBook) { c.kinds[KindAudioBook]++ }

type CheapItems struct {
	maxPriceCents int
	found         []Node
}

func NewCheapItems(maxPriceCents int) *CheapItems {
	return &CheapItems{maxPriceCents: maxPriceCents, found: make([]Node, 0)}
}

func (c *CheapItems) Found() []Node {
	result := make([]Node, 0, len(c.found))
	return append(result, c.found...)
}

func (c *CheapItems) consider(priceCents int, item Node) {
	if priceCents <= c.maxPriceCents {
		c.found = append(c.found, item)
	}
}

func (c *CheapItems) VisitLibrary(l *Library)     { acceptAll(c, l.Sections) }
func (c *CheapItems) VisitSection(s *Section)     { acceptAll(c, s.Shelves) }
func (c *CheapItems) VisitShelf(s *Shelf)         { acceptAll(c, s.Items) }
func (c *CheapItems) VisitBook(b *Book)           { c.consider(b.PriceCents, b) }
func (c *CheapItems) VisitJournal(j *Journal)     { c.consider(j.PriceCents, j) }
func (c *CheapItems) VisitEBook(e *EBook)         { c.consider(e.PriceCents, e) }
func (c *CheapItems) VisitAudioBook(a *AudioBook) { c.consider(a.PriceCents, a) }

type JSONExport struct {
	out strings.Builder
}

func NewJSONExport() *JSONExport {
	return &JSONExport{}
}

func (j *JSONExport) String() string { return j.out.String() }

func exportAll[T Node](j *JSONExport, nodes []T) {
	for i, n := range nodes {
		if i > 0 {
			j.out.WriteByte(',')
		}
		n.Accept(j)
	}
}

func (j *JSONExport) text(s string) {
	j.out.WriteByte('"')
	for _, r := range s {
		switch {
		case r == '"':
			j.out.WriteString(`\"`)
		case r == '\\':
			j.out.WriteString(`\\`)
		case r == '\n':
			j.out.WriteString(`\n`)
		case r == '\r':
			j.out.WriteString(`\r`)
		case r == '\t':
			j.out.WriteString(`\t`)
		case r < 0x20:
			fmt.Fprintf(&j.out, `\u%04x`, r)
		default:
			j.out.WriteRune(r)
		}
	}
	j.out.WriteByte('"')
}

func (j *JSONExport) item(kind string, id int, title string, priceCents int, extraKey string, extra int) {
	j.out.WriteString(`{"type":"` + kind + `","id":` + strconv.Itoa(id) + `,"title":`)
	j.text(title)
	j.out.WriteString(`,"priceCents":` + strconv.Itoa(priceCents))
	j.out.WriteString(`,"` + extraKey + `":` + strconv.Itoa(extra) + `}`)
}

func (j *JSONExport) VisitLibrary(l *Library) {
	j.out.WriteString(`{"type":"library","name":`)
	j.text(l.Name)
	j.out.WriteString(`,"sections":[`)
	exportAll(j, l.Sections)
	j.out.WriteString(`]}`)
}

func (j *JSONExport) VisitSection(s *Section) {
	j.out.WriteString(`{"type":"section","id":` + strconv.Itoa(s.ID) + `,"name":`)
	j.text(s.Name)
	j.out.WriteString(`,"shelves":[`)
	exportAll(j, s.Shelves)
	j.out.WriteString(`]}`)
}

func (j *JSONExport) VisitShelf(s *Shelf) {
	j.out.WriteString(`{"type":"shelf","label":`)
	j.text(s.Label)
	j.out.WriteString(`,"items":[`)
	exportAll(j, s.Items)
	j.out.WriteString(`]}`)
}

func (j *JSONExport) VisitBook(b *Book) {
	j.item(KindBook, b.ID, b.Title, b.PriceCents, "pages", b.Pages)
}

func (j *JSONExport) VisitJournal(jr *Journal) {
	j.item(KindJournal, jr.ID, jr.Title, jr.PriceCents, "issue", jr.Issue)
}

func (j *JSONExport) VisitEBook(e *EBook) {
	j.item(KindEBook, e.ID, e.Title, e.PriceCents, "sizeMB", e.SizeMB)
}

func (j *JSONExport) VisitAudioBook(a *AudioBook) {
	j.item(KindAudioBook, a.ID, a.Title, a.PriceCents, "minutes", a.Minutes)
}
