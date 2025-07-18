package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func startCLI() {
	println("Starting CLI")
	reader := bufio.NewReader(os.Stdin)
	// Here I implement the CLI logic, such as reading commands,
	// processing them, and interacting with the supervisors.
	// This could involve creating instances of S1 and S2, processing events,
	// and printing the state of the supervisors based on user input.
	println("CLI started. Type commands to interact with the supervisors.\nType 'help' to see available commands.")
	for {
		var command string
		print("Enter command: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			println("Error reading command:", err.Error())
			continue
		}
		command = strings.TrimSpace(line)
		switch command {
		case "exit":
			println("Exiting CLI.")
			return
		case "help":
			println("Available commands:")
			println("  start - Start the manual simulation")
			println("  status - Display the status of supervisors")
			println("  exit - Exit the CLI")
		case "start":
			fmt.Println("Starting supervisors...")
			fmt.Println("Type 'rreq <nodeid>' to send RREQ")
			var pathlogs []string
			for {
				var input string
				fmt.Print("Enter event: ")
				line, err := reader.ReadString('\n')
				if err != nil {
					println("Error reading event:", err.Error())
					continue
				}
				input = strings.TrimSpace(line)
				parts := strings.Fields(input)
				if len(parts) == 0 {
					continue
				}
				event := parts[0]
				if event == "stop" {
					println("Stopping manual simulation.")
					break
				}
				if parts[0] == "rreq" && len(parts) < 2 {
					println("Please provide both event and nodeid (e.g., rreq 1)")
					continue
				}
				var nodeid string
				if len(parts) != 1 {
					nodeid = parts[1]
					go func() {
						pathlogs = manualSimulation(event, nodeid)
					}()
					continue
				}

				if event == "logs" {
					if len(pathlogs) == 0 {
						println("No logs available.")
					} else {
						println("Path logs:")
						for _, log := range pathlogs {
							println(log)
						}
					}
				}
			}

		default:
			println("Unknown command:", command)
		}
		time.Sleep(500 * time.Millisecond) // Simulate processing delay
	}
}
