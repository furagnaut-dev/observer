package library

const (
	KindBook      = "book"
	KindJournal   = "journal"
	KindEBook     = "ebook"
	KindAudioBook = "audiobook"
)

type Library struct {
	Name     string
	Sections []*Section
}

type Section struct {
	ID      int
	Name    string
	Shelves []*Shelf
}

type Shelf struct {
	Label string
	Items []Node
}

type Book struct {
	ID         int
	Title      string
	PriceCents int
	Pages      int
}

type Journal struct {
	ID         int
	Title      string
	PriceCents int
	Issue      int
}

type EBook struct {
	ID         int
	Title      string
	PriceCents int
	SizeMB     int
}

type AudioBook struct {
	ID         int
	Title      string
	PriceCents int
	Minutes    int
}

func (l *Library) Accept(v Visitor)   { v.VisitLibrary(l) }
func (s *Section) Accept(v Visitor)   { v.VisitSection(s) }
func (s *Shelf) Accept(v Visitor)     { v.VisitShelf(s) }
func (b *Book) Accept(v Visitor)      { v.VisitBook(b) }
func (j *Journal) Accept(v Visitor)   { v.VisitJournal(j) }
func (e *EBook) Accept(v Visitor)     { v.VisitEBook(e) }
func (a *AudioBook) Accept(v Visitor) { v.VisitAudioBook(a) }
