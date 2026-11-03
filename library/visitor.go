package library

type Node interface {
	Accept(v Visitor)
}

type Visitor interface {
	VisitLibrary(l *Library)
	VisitSection(s *Section)
	VisitShelf(s *Shelf)
	VisitBook(b *Book)
	VisitJournal(j *Journal)
	VisitEBook(e *EBook)
	VisitAudioBook(a *AudioBook)
}

func acceptAll[T Node](v Visitor, nodes []T) {
	for _, n := range nodes {
		n.Accept(v)
	}
}
