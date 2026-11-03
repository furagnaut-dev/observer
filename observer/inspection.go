package observer

import (
	"time"

	"observer/library"
)

type InspectionService struct {
	*ObserverRegistry
	tree *library.Library
}

func NewInspectionService(registry *ObserverRegistry, tree *library.Library) *InspectionService {
	return &InspectionService{ObserverRegistry: registry, tree: tree}
}

func (s *InspectionService) OnEvent(e Event) {
	finder := library.NewSectionFinder(e.SectionID())
	s.tree.Accept(finder)
	section, found := finder.Found()
	if !found {
		return
	}
	value := library.NewInventoryValue()
	section.Accept(value)
	s.notifySubscribers(NewEvent(InspectionDone, e.SectionID(), e.BookID(), e.ReaderID(), value.Total(), time.Time{}))
}
