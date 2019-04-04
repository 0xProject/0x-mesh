package core

type EventType uint

const (
	Stored EventType = iota
	Evicted
)

type Event struct {
	Type    EventType
	Message *Message
}
