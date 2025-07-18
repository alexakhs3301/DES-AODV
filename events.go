package main

// EventType represents a type for events in the DES system.
type EventType string

const (
	EventRREQSend    EventType = "RREQ_SEND"
	EventRREPRecv    EventType = "RREP_RECV"
	EventDataSend    EventType = "DATA_SEND"
	EventRREPForward EventType = "RREP_FORWARD"
	EventOtherRREQ   EventType = "OTHER_RREQ"
)

type Event struct {
	SourceID int
	Type     EventType
}
