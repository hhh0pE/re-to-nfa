package main

import (
	"os"
    "github.com/hhh0pE/REtoNFA/NFA"
)

func main() {
	if len(os.Args) > 1 {
		nfa := NFA.BuildNFA(os.Args[1])
		nfa.PrintNFA()
		nfa.PrintJSON()

        if len(os.Args) >=3 {
            nfa.SaveToFile(os.Args[2])
        }
	} else {
		panic("You must pass RE as parameter. Nothing passed. Nothing to do. Exit.")
	}

}
