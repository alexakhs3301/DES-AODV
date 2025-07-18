package main

import (
	"fmt"
	"sync"
)

type NodeState string

const (
	Idle       NodeState = "Idle"
	Waiting    NodeState = "Waiting"
	Forwarding NodeState = "Forwarding"
)

type Node struct {
	id    int
	state NodeState
	inbox chan Event

	s1 *S1 // Supervisor 1
	s2 *S2 // Supervisor 2

	mu sync.Mutex
}

func NewNode(id int) *Node {
	n := &Node{
		id:    id,
		state: Idle,
	}
	n.s1 = NewSupervisor1(id)
	n.s2 = NewSupervisor2(id)
	return n
}

func (n *Node) SetInbox(ch chan Event) {
	n.inbox = ch
}

func (n *Node) Run() {
	for event := range n.inbox {
		n.mu.Lock()
		switch event.Type {
		case EventRREQSend:
			if n.s1.CanSendRREQ() && event.SourceID == n.ID() {
				fmt.Println("-> Node", n.ID(), "sending RREQ")
				n.SetState(Waiting)
			} else if n.ID() == event.SourceID {
				fmt.Println("-> Node", n.ID(), "blocked from sending RREQ by Supervisor")
			}
		case EventRREPRecv:
			if n.ID() == event.SourceID && n.s2.CanAcceptRREP() {
				fmt.Println("-> Node", n.ID(), "accepting RREP")
				n.SetState(Forwarding)
			} else if n.ID() == event.SourceID {
				fmt.Println("-> Node", n.ID(), "blocked from accepting RREP by Supervisor")
			}
		case EventDataSend:
			if n.ID() == event.SourceID {
				fmt.Println("-> Node", n.ID(), "sending DATA")
				n.SetState(Idle)
			} else if n.ID() != event.SourceID {
				n.SetState(Idle)
			} else {
				fmt.Println("-> Node", n.ID(), "received DATA from another node")
			}
		case EventRREPForward:
			fmt.Println("-> Node", n.ID(), "received RREP_FORWARD event")
		case EventOtherRREQ:
			fmt.Println("-> Node", n.ID(), "received OTHER_RREQ event")
		}
		n.s1.Process(event)
		n.s2.Process(event)

		n.mu.Unlock()
		fmt.Println("Event processed:", event.Type, "from Node", event.SourceID)
		fmt.Println("->Node", n.ID(), "state: ", n.State())
		fmt.Println("----------------------")
	}
}

func (n *Node) ID() int {
	return n.id
}
func (n *Node) State() NodeState {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.state
}
func (n *Node) SetState(state NodeState) {
	n.state = state
}
