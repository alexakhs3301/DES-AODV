package main

import (
	"fmt"
	"strconv"
	"sync"
)

func simulate() {
	node1 := NewNode(1)
	node2 := NewNode(2)

	node1Inbox := make(chan Event)
	node2Inbox := make(chan Event)

	node1.SetInbox(node1Inbox)
	node2.SetInbox(node2Inbox)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		node1.Run()
		wg.Done()
	}()
	go func() {
		node2.Run()
		wg.Done()
	}()

	// events
	events := []Event{
		{SourceID: 1, Type: EventRREQSend},    // Allowed: Node 1 sends RREQ (S1: q1→q2)
		{SourceID: 1, Type: EventRREQSend},    // Blocked: Node 1 tries again (still q2)
		{SourceID: 1, Type: EventRREQSend},    // Allowed: Node 1 sends RREQ (S1: q1→q2)
		{SourceID: 1, Type: EventRREQSend},    // Allowed: Node 1 sends RREQ (S1: q1→q2)
		{SourceID: 1, Type: EventRREPRecv},    // Blocked: Node 1 tries to accept RREP (S2: p1)
		{SourceID: 2, Type: EventRREPForward}, // Unlock: Node 2 forwards RREP (Node 1 S2: p1→p2)
		{SourceID: 1, Type: EventRREPRecv},    // Allowed: Node 1 accepts RREP (S2: p2→p1)
		{SourceID: 2, Type: EventRREQSend},    // Allowed: Node 2 sends RREQ (S1: q1→q2)
		{SourceID: 2, Type: EventRREQSend},    // Blocked: Node 2 tries again (still q2)
		{SourceID: 1, Type: EventDataSend},    // Unlock: Node 1 sends DATA (Node 2 S1: q2→q1)
		{SourceID: 2, Type: EventRREQSend},    // Allowed: Node 2 sends RREQ (S1: q1→q2)
		{SourceID: 1, Type: EventRREPForward}, // Unlock: Node 1 forwards RREP (Node 2 S1: q2→q1)
		{SourceID: 2, Type: EventRREQSend},    // Allowed: Node 2 sends RREQ (S1: q1→q2)
		{SourceID: 2, Type: EventRREPRecv},    // Blocked: Node 2 tries to accept RREP (S2: p1)
		{SourceID: 1, Type: EventRREPForward}, // Unlock: Node 1 forwards RREP (Node 2 S2: p1→p2)
		{SourceID: 2, Type: EventRREPRecv},    // Allowed: Node 2 accepts RREP (S2: p2→p1)
	}

	//events := []Event{
	//	// Step 1: Node 1 tries to accept RREP (should be blocked - S2 starts in p1)
	//	{SourceID: 1, Type: EventRREPRecv},
	//
	//	// Step 2: Node 2 forwards RREP, unlocking Node 1's S2
	//	{SourceID: 2, Type: EventRREPForward},
	//
	//	// Step 3: Node 1 tries to accept RREP again (should be allowed now)
	//	{SourceID: 1, Type: EventRREPRecv},
	//
	//	// Step 4: Node 1 tries to accept RREP once more (should be blocked again)
	//	{SourceID: 1, Type: EventRREPRecv},
	//}

	for _, event := range events {
		node1Inbox <- event
		node2Inbox <- event
	}
	close(node1Inbox)
	close(node2Inbox)
	wg.Wait()

}

func manualSimulation(event string, nodeid string) []string {
	nodeID, _ := strconv.Atoi(nodeid)
	node1 := NewNode(nodeID)
	otherNodeID := 2
	if node1.ID() == otherNodeID {
		otherNodeID = 1
	}
	node2 := NewNode(otherNodeID)
	node1Inbox := make(chan Event)
	node2Inbox := make(chan Event)

	node1.SetInbox(node1Inbox)
	node2.SetInbox(node2Inbox)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		node1.Run()
		wg.Done()
	}()
	go func() {
		node2.Run()
		wg.Done()
	}()

	var pathlogs []string

	switch event {
	case "rreq":
		if node1.s1.CanSendRREQ() {
			//rreq send
			node1Inbox <- Event{SourceID: node1.ID(), Type: EventRREQSend}
			node2Inbox <- Event{SourceID: node1.ID(), Type: EventRREQSend}
			log := fmt.Sprintln("Node", node1.ID(), "sent RREQ")
			pathlogs = append(pathlogs, log)

			//rreq receive
			node1Inbox <- Event{SourceID: node1.ID(), Type: EventOtherRREQ}
			node2Inbox <- Event{SourceID: node1.ID(), Type: EventOtherRREQ}
			log = fmt.Sprintln("Node", node2.ID(), "received RREQ from Node", node1.ID())
			pathlogs = append(pathlogs, log)

			//rrep forward
			node1Inbox <- Event{SourceID: node2.ID(), Type: EventRREPForward}
			node2Inbox <- Event{SourceID: node2.ID(), Type: EventRREPForward}
			log = fmt.Sprintln("Node", node2.ID(), "forwarded RREP to Node", node1.ID())
			pathlogs = append(pathlogs, log)
			//rrep receive
			node1Inbox <- Event{SourceID: node1.ID(), Type: EventRREPRecv}
			node2Inbox <- Event{SourceID: node1.ID(), Type: EventRREPRecv}
			log = fmt.Sprintln("Node", node1.ID(), "received RREP from Node", node2.ID())
			pathlogs = append(pathlogs, log)

			// data send
			node1Inbox <- Event{SourceID: node1.ID(), Type: EventDataSend}
			node2Inbox <- Event{SourceID: node1.ID(), Type: EventDataSend}
			log = fmt.Sprintln("Node", node1.ID(), "sent DATA to Node", node2.ID())
			pathlogs = append(pathlogs, log)
			//
		} else {
			println("Node", node1.ID(), "blocked from sending RREQ by Supervisor")
		}
	}
	close(node1Inbox)
	close(node2Inbox)
	wg.Wait()

	return pathlogs
}
