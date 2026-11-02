package observer

type Observer interface {
	OnEvent(Event)
}
