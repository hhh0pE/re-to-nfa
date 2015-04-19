package main

import "./NFA"
import "testing"

func TestPatterns(t *testing.T) {
	patterns := []string{}
	patterns = append(patterns, "a+b")
	patterns = append(patterns, "(a+b)")
	patterns = append(patterns, "(a+b)*")
	patterns = append(patterns, "(a+b)+(c+d)")
	patterns = append(patterns, "a+c+d+e+f+g")
	patterns = append(patterns, "(a+b)*+c")
	patterns = append(patterns, "(a+b)*+(a+bb)")
	patterns = append(patterns, "(a+b)*+((a+bb))")
	patterns = append(patterns, "((a+b*)*abb)*")
	patterns = append(patterns, "(aba((aaaaa*ab*(ab)*)*+(abb(a+b+a+b)*)*aaab*)+a*b*(ab)*ab)*ababa")

	for i, pattern := range patterns {
		nfa := NFA.BuildNFA(pattern)
		if nfa == nil {
			t.Errorf("Test %d/%d failed for pattern %s", i, len(patterns), pattern)
		}
	}
}

func TestNodesCount(t *testing.T) {
    patterns := make(map[string]int)
    patterns["a+b"] = 6
    patterns["bb"] = 5
    patterns["(a+b)*"] = 8
    patterns["(a+b)*((a+bb))"] = 17

    var i int
    i = 0
    for pattern, count_must_be := range patterns {
        i++
        nfa := NFA.BuildNFA(pattern)
        if nfa.Length() != count_must_be {
            t.Errorf("Test count %d/%d failed for pattern %s. Must be %d but recieved %d", i+1, len(patterns), pattern, count_must_be, nfa.Length())
        }
    }
}

func TestWithEtalons(t *testing.T) {

    // a
    node_a := NFA.Node{}
    node_a2 := NFA.Node{}

    node_a.Left = &node_a2
    node_a.LeftSymbol = "a"

    nfa_a := NFA.NewNFA(&node_a, &node_a2)

    // ab
    var node_ab1, node_ab2, node_ab3, node_ab4, node_ab5 NFA.Node
    node_ab1.Left = &node_ab2
    node_ab2.Left = &node_ab3
    node_ab2.LeftSymbol = "a"
    node_ab3.Left = &node_ab4
    node_ab3.LeftSymbol = "b"
    node_ab4.Left = &node_ab5

    nfa_ab := NFA.NewNFA(&node_ab1, &node_ab2, &node_ab3, &node_ab4, &node_ab5)

    nfa1 := NFA.BuildNFA("a")
    nfa2 := NFA.BuildNFA("ab")

    if nfa_a.Hash() != nfa1.Hash() {
        t.Errorf("NFA of pattern \"%s\" is not equal to etalon.", string('a'))
    }
    if nfa_ab.Hash() != nfa2.Hash() {
        t.Errorf("NFA of pattern \"%s\" is not equal to etalon.\nNFA1 Hash:%x\nEtalon Hash:%x", "ab", nfa_ab.Hash(), nfa1.Hash())
    }
}
