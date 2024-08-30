package main

import (
	"flag"
)

func main() {
	cliMode := flag.Bool("cli", false, "Run in CLI mode")

	flag.Parse()

	if(*cliMode) {
		runCLI()
	} else {
		handleRequests()
	}
}