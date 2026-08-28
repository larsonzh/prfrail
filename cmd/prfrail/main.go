// Command prfrail is the ProofRail CLI entry point.
//
// This is a scaffold stub. The design baseline is defined in
// docs/RFC-proofrail-unattended-ai-engineering-product.md; full commands
// (init/validate/run/report/serve) land in S0/S1 milestones.
package main

import (
	"fmt"
	"os"
)

const version = "0.1.0-dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("prfrail %s\n", version)
	case "init":
		// S1: wizard -> generate chain-file + workspace config.
		fmt.Println("prfrail init: not implemented yet (see RFC S1)")
	case "run":
		// S1: execute task chain with TUI progress.
		fmt.Println("prfrail run: not implemented yet (see RFC S1)")
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: prfrail <%s>\n", "version|init|run")
}
