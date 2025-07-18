package main

import (
	"flag"
)

func main() {
	mode := flag.String("mode", "sim", "Mode of operation: 'sim' or 'cli'")
	flag.Parse()
	switch *mode {
	case "cli":
		// CLI mode: start the CLI interface
		startCLI()
	case "sim":
		// Simulation mode: run the simulation
		simulate()
	default:
		// Invalid mode: print usage and exit
		flag.Usage()
		println("Invalid mode. Use 'sim' for simulation or 'cli' for command line interface.")
		return
	}
}
