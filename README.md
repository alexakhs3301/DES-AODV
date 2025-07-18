# AODV Protocol Simulator

An Ad-hoc On-Demand Distance Vector (AODV) routing protocol implemented as a Discrete Event Simulation (DES) Automaton in Go. This project incorporates supervisory control theory to formally verify protocol behavior.

## Overview

This project simulates AODV routing behavior using nodes modeled as DES automata. Each node can send Route Requests (RREQ), receive Route Replies (RREP), and transmit data. Two supervisory controllers ensure correct and verifiable protocol behavior.

## Features

- Discrete Event Simulation of AODV protocol  
- Supervisor automata (S1 & S2) for protocol enforcement  
- Formal verification of routing correctness  
- Dual-mode operation: **Simulation** and **Interactive CLI**

## Supervisory Control Model

### Supervisor 1 (S1) — RREQ Control
- **States**:
  - `q1`: Allowed to send RREQ
  - `q2`: Blocked from sending RREQ
- **Transitions**:
  - `RREQ_SEND`: → q2
  - `DATA_SEND`, `RREP_FORWARD` (from other nodes): → q1

### Supervisor 2 (S2) — RREP Acceptance Control
- **States**:
  - `p1`: Cannot accept RREP
  - `p2`: Allowed to accept RREP
- **Transitions**:
  - `DATA_SEND`, `RREP_FORWARD` (from other nodes): → p2
  - `RREP_RECV`: → p1

## Usage

### Build

```bash
go build
```

### Simulation Mode

Run the predefined event sequence:

```bash
./AODV -mode sim
```

### Interactive CLI Mode

Manually interact with the nodes:

```bash
./AODV -mode cli
```

**Available CLI commands:**

- `start` – Start manual simulation  
- `rreq <nodeid>` – Send RREQ from specified node  
- `logs` – Display event logs  
- `status` – Display supervisor status  
- `stop` – Stop manual simulation  
- `exit` – Exit the CLI  

---

## Project Structure

| File              | Description                                 |
|-------------------|---------------------------------------------|
| `main.go`         | Entry point and mode selector                |
| `node.go`         | Node behavior and state transitions          |
| `s1.go`           | Supervisor 1 – RREQ control                  |
| `s2.go`           | Supervisor 2 – RREP control                  |
| `events.go`       | Event types and definitions                  |
| `simulation.go`   | Simulation logic for predefined execution    |
| `cli.go`          | Command Line Interface logic                 |

---

## Testing

Run tests to validate supervisor logic and node behavior:

```bash
go test -v
```

---

## Implementation Highlights

- DES-based simulation of AODV protocol behavior  
- Supervisors enforce correct event sequencing  
- Goroutines and channels simulate asynchronous execution  
- Formal verification via language intersection constraints  

---

## Example Event Trace

```go
{SourceID: 1, Type: EventRREQSend}    // Allowed (enters q2)
{SourceID: 1, Type: EventRREQSend}    // Blocked by Supervisor 1
{SourceID: 2, Type: EventDataSend}    // Unlocks Node 1 (S1 returns to q1)
{SourceID: 1, Type: EventRREQSend}    // Allowed again
{SourceID: 1, Type: EventRREPRecv}    // Blocked by Supervisor 2
{SourceID: 2, Type: EventRREPForward} // Unlocks Node 1 (S2 to p2)
{SourceID: 1, Type: EventRREPRecv}    // Allowed, S2 transitions back to p1
```

This sequence demonstrates the correct enforcement of routing discipline by the supervisors.

---

## Requirements

- Go 1.24 or later

---

## Incoming Features

The next development steps include extending the interactive CLI interface to support:

- **Extension of the CLI events** — Let the user queue events for each node in real time  
- **Live state monitoring** — Continuously display the current state of each node and supervisor  
- **Custom long topologies** — Support more than two nodes and define network structure interactively  
- **Undo / Replay functionality** — Navigate backward or forward through the simulation steps


## References

1. Fragkoulis et al., *Modelling and modular supervisory control for the AODV routing protocol*, AEÜ, 2023  
2. Perkins et al., *Ad hoc On-Demand Distance Vector (AODV) Routing*, RFC 3561, IETF, 2003  
3. Cassandras & Lafortune, *Introduction to Discrete Event Systems*, Springer, 2021  
4. Broch et al., *A performance comparison of multi-hop wireless ad hoc network routing protocols*, MobiCom, 1998  
5. The Go Programming Language. [https://go.dev/doc](https://go.dev/doc)
6. M. K. Marina and S. R. Das, "On-demand multipath distance vector routing in ad hoc networks," in Proceedings of the Ninth International Conference on Network Protocols, 2001.


© 2025 – AODV-DES Golang Simulator | Developed for academic research and formal protocol verification.

---
