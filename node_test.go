package main

import (
	"sync"
	"testing"
	"time"
)

func TestNodeEventSequence(t *testing.T) {
	node := NewNode(1)
	inbox := make(chan Event, 3)
	node.SetInbox(inbox)

	events := []Event{
		{SourceID: 1, Type: EventRREQSend}, // Should send RREQ, state: Waiting
		{SourceID: 2, Type: EventDataSend}, // Should unlock S1, state: Idle
		{SourceID: 1, Type: EventRREQSend}, // Should send RREQ again, state: Waiting
		{SourceID: 2, Type: EventDataSend}, // Should unlock S1 again, state: Idle
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		node.Run()
		wg.Done()
	}()

	for _, e := range events {
		inbox <- e
	}
	close(inbox)
	wg.Wait()

	// Wait for processing
	// (In real tests, use sync primitives or refactor for synchronous processing)

	// Check final state
	if node.State() != Waiting {
		t.Errorf("Expected node state Waiting, got %v", node.State())
	}
	if !node.s1.CanSendRREQ() {
		t.Errorf("Expected S1 to allow RREQ after DATA_SEND from other node")
	}
}

func TestNodeBlockRREP(t *testing.T) {
	node := NewNode(1)
	inbox := make(chan Event, 2)
	node.SetInbox(inbox)

	events := []Event{
		{SourceID: 1, Type: EventRREPRecv},    // Should block RREP at start
		{SourceID: 2, Type: EventRREPForward}, // Should unlock S2
		{SourceID: 1, Type: EventRREPRecv},    // Should block RREP again after receiving
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		node.Run()
		wg.Done()
	}()

	for _, e := range events {
		inbox <- e
	}
	wg.Wait()
	close(inbox)

	// Wait for processing
	// (In real tests, use sync primitives or refactor for synchronous processing)

	if node.s2.CanAcceptRREP() {
		t.Errorf("Expected S2 to block RREPs after receiving RREP")
	}
}

func TestNodeUnlockRREQ(t *testing.T) {
	node := NewNode(1)
	inbox := make(chan Event, 2)
	node.SetInbox(inbox)

	events := []Event{
		{SourceID: 1, Type: EventRREQSend},    // Should send RREQ, state: Waiting
		{SourceID: 2, Type: EventRREPForward}, // Should unlock S1
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		node.Run()
		wg.Done()
	}()

	for _, e := range events {
		inbox <- e
	}
	close(inbox)
	wg.Wait()

	// Wait for processing
	// (In real tests, use sync primitives or refactor for synchronous processing)

	if !node.s1.CanSendRREQ() {
		t.Errorf("Expected S1 to allow RREQ after receiving RREP_FORWARD")
	}
}

func TestRequestfromtwonodes(t *testing.T) {
	node1 := NewNode(1)
	node2 := NewNode(2)

	inbox1 := make(chan Event, 3)
	inbox2 := make(chan Event, 3)

	node1.SetInbox(inbox1)
	node2.SetInbox(inbox2)

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

	// Simulate events from both nodes
	events := []Event{
		{SourceID: 1, Type: EventRREQSend}, // Node 1 sends RREQ
		{SourceID: 2, Type: EventRREQSend}, // Node 2 sends RREQ
	}

	for _, e := range events {
		inbox1 <- e
		inbox2 <- e
	}

	close(inbox1)
	close(inbox2)
	wg.Wait()

	if !node1.s1.CanSendRREQ() || !node2.s1.CanSendRREQ() {
		t.Error("Expected both nodes to allow RREQ after processing")
	}
}

func TestRREQMutualExclusion(t *testing.T) {
	node1 := NewNode(1)
	node2 := NewNode(2)

	inbox1 := make(chan Event, 3)
	inbox2 := make(chan Event, 3)
	node1.SetInbox(inbox1)
	node2.SetInbox(inbox2)

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

	// Step 1: Node 1 sends RREQ
	node1.inbox <- Event{SourceID: 1, Type: EventRREQSend}
	time.Sleep(100 * time.Millisecond)

	if node1.s1.State != q2 {
		t.Errorf("Expected node1 supervisor1 to be in q2 after RREQ, got %v", node1.s1.State)
	}

	// Step 2: Node 2 sends RREQ
	node2.inbox <- Event{SourceID: 2, Type: EventRREQSend}
	time.Sleep(100 * time.Millisecond)

	if node2.s1.State != q2 {
		t.Errorf("Expected node2 supervisor1 to be in q2 after RREQ, got %v", node2.s1.State)
	}

	// Step 3: Node 1 tries RREQ again (should be blocked)
	canSend := node1.s1.CanSendRREQ()
	if canSend {
		t.Error("Node 1 should be blocked from sending second RREQ (still in q2)")
	}

	// Step 4: Node 2 sends DATA_SEND to unlock node1
	node1.inbox <- Event{SourceID: 2, Type: EventDataSend}
	node2.inbox <- Event{SourceID: 2, Type: EventDataSend}
	time.Sleep(100 * time.Millisecond)

	// Step 5: Node 1 should now be allowed
	canSendAfter := node1.s1.CanSendRREQ()
	if !canSendAfter {
		t.Error("Node 1 should be allowed to send RREQ after Node 2's DATA_SEND")
	}
	canSendAfterNode2 := node2.s1.CanSendRREQ()
	if !canSendAfterNode2 {
		t.Error("Node 2 should be allowed to send RREQ after Node 1's DATA_SEND")
	}
	close(node1.inbox)
	close(node2.inbox)
	wg.Wait()
}
