package main

import "testing"

func TestSupervisor1(t *testing.T) {
	s1 := NewSupervisor1(1)

	// Αρχικά: πρέπει να επιτρέπει RREQ_SEND
	if !s1.CanSendRREQ() {
		t.Error("Expected to allow RREQ at start")
	}

	// Στέλνει RREQ ο ίδιος -> πρέπει να μπλοκάρει από εδώ και πέρα
	s1.Process(Event{SourceID: 1, Type: EventRREQSend})
	if s1.CanSendRREQ() {
		t.Error("Expected RREQ to be blocked after sending")
	}

	// Άλλος κόμβος στέλνει DATA -> unlock
	s1.Process(Event{SourceID: 2, Type: EventDataSend})
	if !s1.CanSendRREQ() {
		t.Error("Expected RREQ to be allowed after other node sent DATA")
	}
}

func TestSupervisor2(t *testing.T) {
	s2 := NewSupervisor2(1)

	// Αρχικά: δεν επιτρέπεται RREP_RECV
	if s2.CanAcceptRREP() {
		t.Error("Expected RREP to be blocked at start")
	}

	// Άλλος κόμβος στέλνει RREP_FORWARD -> ξεκλειδώνει
	s2.Process(Event{SourceID: 2, Type: EventRREPForward})
	if !s2.CanAcceptRREP() {
		t.Error("Expected RREP to be allowed after other node forwarded RREP")
	}

	// Ο κόμβος μου λαμβάνει RREP -> κλειδώνει ξανά
	s2.Process(Event{SourceID: 1, Type: EventRREPRecv})
	if s2.CanAcceptRREP() {
		t.Error("Expected RREP to be blocked again after accepting")
	}
}
