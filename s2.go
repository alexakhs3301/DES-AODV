package main

type States string

const (
	p1 States = "p1" // locked
	p2 States = "p2" // unlocked
)

type S2 struct {
	NodeID int
	State  States
}

func NewSupervisor2(nodeID int) *S2 {
	return &S2{
		NodeID: nodeID,
		State:  p1, // Initial state is p1
	}
}

func (s *S2) CanAcceptRREP() bool {
	return s.State == p2
}

func (s *S2) Process(event Event) {
	switch s.State {
	case p1:
		if event.Type == EventDataSend || event.Type == EventRREPForward {
			if event.SourceID != s.NodeID {
				s.State = p2
			}
		}
	case p2:
		if event.Type == EventRREPRecv && event.SourceID == s.NodeID {
			s.State = p1
		}
	}
}
