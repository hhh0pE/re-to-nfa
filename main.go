package main

import (
	"os"
)

func main() {
	if len(os.Args) > 1 {
		nfa := buildNFA(os.Args[1])
		nfa.printNFA()
	} else {
		panic("You must pass RE as parameter. Nothing passed. Nothing to do. Exit.")
	}

}
