package main

import (
	"os"
)

func main() {
	nfa := buildNFA(os.Args[1])
	nfa.printNFA()
}
