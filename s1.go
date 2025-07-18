package main

type StateS1 int

const (
	q1 StateS1 = iota
	q2         //not allowed
)

type S1 struct {
	NodeID int
	State  StateS1
}

func NewSupervisor1(nodeID int) *S1 {
	return &S1{
		NodeID: nodeID,
		State:  q1, // Initial state is q1
	}
}

// Process an event and update state
func (s *S1) Process(event Event) {
	switch s.State {
	case q1:
		if event.Type == EventRREQSend && event.SourceID == s.NodeID {
			s.State = q2
		} else if event.Type == EventRREQSend && event.SourceID != s.NodeID {
			s.State = q2
		}
	case q2:
		if event.Type == EventDataSend || event.Type == EventRREPForward {
			s.State = q1
		}
	}
}

func (s *S1) CanSendRREQ() bool {
	return s.State == q1
}
