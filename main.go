package main

import (
	"os"
    "github.com/hhh0pE/REtoNFA/NFA"
)

func main() {
	if len(os.Args) > 1 {
		nfa := NFA.BuildNFA(os.Args[1])

        if len(os.Args) >=3 {
            param2 := os.Args[2]
            if param2 == "JSON" {
                nfa.PrintJSON()
                return
            }
            if param2 == "Visual" {
                nfa.PrintNFA()
                return
            }
            nfa.SaveToFile(param2)
        }
	} else {
		panic("You must pass RE as parameter. Nothing passed. Nothing to do. Exit.")
	}

}
