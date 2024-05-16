package main

import "fmt"

func runCLI() {
	exit := false
	for !exit {
		var input string
		println("1) Generate Dockerfile")
		println("q) Exit")
		println("Enter a command (1,q):")
		_, _ = fmt.Scanln(&input)
		switch input {
		case "q":
			exit = true
		case "1":
		default:
			println("Unknown command")
		}
	}
}