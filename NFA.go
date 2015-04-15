package main

import (
	"fmt"
	"strings"
)

type NFA struct {
	begin, end *Node
	nodes      []*Node
}

func (nfa *NFA) addNode(new_node *Node) {
	nfa.nodes = append(nfa.nodes, new_node)
}

func (nfa *NFA) addNodes(new_nodes []*Node) {
	nfa.nodes = append(nfa.nodes, new_nodes...)
}

func NewNFA(nodes ...*Node) *NFA {
	var nfa NFA
	nfa.addNodes(nodes)
	for i, v := range nodes {
		v.Name = fmt.Sprintf("Node%d", i+1)
	}
	nfa.begin = nodes[0]
	nfa.end = nodes[len(nodes)-1]

	return &nfa
}

func (nfa *NFA) printNFA() {
	if nfa != nil {
		fmt.Println("NFA built success. Printing NFA..")
		for _, node := range nfa.nodes {
			fmt.Println(node.toString())
		}
	} else {
		fmt.Println("Error when building NFA")
	}
}

// for a+b
func addNFA(a *NFA, b *NFA) *NFA {
	if a == nil || b == nil {
		return nil
	}
	var nodes []*Node

	var start, end Node

	nodes = append(nodes, &start)
	nodes = append(nodes, a.nodes...)
	nodes = append(nodes, b.nodes...)
	nodes = append(nodes, &end)

	start.Left = a.begin
	start.Right = b.begin

	a.end.Left = &end
	b.end.Left = &end

	a = nil
	b = nil

	return NewNFA(nodes...)
}

// for ab (a*b)
func multiplyNFA(a *NFA, b *NFA) *NFA {
	if a == nil || b == nil {
		return nil
	}
	var nodes []*Node
	var end Node

	nodes = append(nodes, a.nodes...)
	nodes = append(nodes, b.nodes...)
	nodes = append(nodes, &end)

	a.end.Left = b.begin
	b.end.Left = &end

	a = nil
	b = nil

	return NewNFA(nodes...)
}

// for a*
func powerNFA(a *NFA) *NFA {
	if a == nil {
		return nil
	}
	var nodes []*Node
	var start, end Node

	nodes = append(nodes, &start)
	nodes = append(nodes, a.nodes...)
	nodes = append(nodes, &end)

	start.Left = a.begin
	start.Right = &end
	a.end.Left = a.begin
	a.end.Right = &end

	a = nil

	return NewNFA(nodes...)
}

func buildNFA(pattern string) *NFA {
	//fmt.Println(pattern)
	if len(pattern) == 1 {
		if strings.Contains(pattern, "+*()") {
			return nil
		}

		// ->Node1-a->Node2
		var node1, node2 Node
		node1.Left = &node2
		node1.LeftSymbol = string(pattern)

		return NewNFA(&node1, &node2)
	}

	//    // has only one ( and )
	//    if strings.Count(pattern, "(") == 1 && strings.Count(pattern, ")") == 1 {
	//        // (...) => ,,,
	//        if pattern[0]=='(' && pattern[len(pattern)-1]==')' {
	//            pattern = strings.Trim(pattern, "()")
	//        }
	//        // (,,.)*
	//        if pattern[len(pattern)-1] == '*' {
	//            return powerNFA(buildNFA(pattern[:len(pattern)-1]))
	//        }
	//    }
	//    // (..)..(..)..
	//    if strings.Count(pattern, "()">0) {
	//
	//    }

	// if ( or ) don't exists
	if strings.Count(pattern, "(") == 0 && strings.Count(pattern, ")") == 0 {
		// ..+..
		if index := strings.Index(pattern, "+"); index > 0 {
			return addNFA(buildNFA(pattern[:index]), buildNFA(pattern[index+1:]))
		}
		// ...*
		if pattern[len(pattern)-1] == '*' {
			return powerNFA(buildNFA(pattern[:len(pattern)-1]))
		}
	}

	// if (..)
	// change to ..
    
	if strings.Count(pattern, "(") == 1 && strings.Count(pattern, ")") == 1 && pattern[0] == '(' && pattern[len(pattern)-1] == ')' {
		pattern = strings.Trim(pattern, "()")
	}

	// a lot of (..)..(..)..(..)
	left_bracket_count, right_bracket_count := 0, 0
	for i, s := range pattern {
		if s == '(' {
			left_bracket_count++
		}
		if s == ')' {
			right_bracket_count++
		}

		if s == '+' && left_bracket_count == right_bracket_count {
			return addNFA(buildNFA(pattern[:i]), buildNFA(pattern[i+1:]))
		}
	}

	if left_bracket_count != right_bracket_count {
		panic("Left and rights bracket doesn't equal!")
	}

	if pattern[0] == '(' && pattern[len(pattern)-1] == ')' {
		pattern = strings.Trim(pattern, "()")
	}

	if strings.Count(pattern, "+") == 0 && strings.Count(pattern, "*") == 0 {
		return multiplyNFA(buildNFA(pattern[:1]), buildNFA(pattern[1:]))
	}

	if pattern[len(pattern)-1] == '*' {
		return powerNFA(buildNFA(pattern[:len(pattern)-1]))
	}

	fmt.Println("!!" + pattern)

	return nil
}
