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
	fmt.Println("Printing NFA..")

	for _, node := range nfa.nodes {
		fmt.Println(node.toString())
	}

}

// for a+b
func addNFA(a *NFA, b *NFA) *NFA {
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
	var nodes []*Node
	var start, end Node

	nodes = append(nodes, &start)
	nodes = append(nodes, a.nodes...)
	nodes = append(nodes, &end)

	start.Left = a.begin
	a.end.Left = a.begin
	a.end.Right = &end

	a = nil

	return NewNFA(nodes...)
}

func buildNFA(pattern string) *NFA {
	if len(pattern) == 3 && pattern[1] == '+' { // a+b
		return addNFA(NewNFA(&Node{LeftSymbol: "a"}), NewNFA(&Node{LeftSymbol: "b"}))
	}
	if len(pattern) == 2 {
		if pattern[1] == '*' { // a*
			return powerNFA(NewNFA(&Node{LeftSymbol: string(pattern[0])}))
		} else { // ab
			return multiplyNFA(NewNFA(&Node{LeftSymbol: string(pattern[0])}), NewNFA(&Node{LeftSymbol: string(pattern[1])}))
		}
	}

	if len(pattern) > 3 {
		if pattern[0] == '(' && pattern[len(pattern)-1] == ')' {
			pattern = strings.Trim(pattern, "()")
		}
		if strings.Contains(pattern, "+") {
			parts := strings.Split(pattern, "+")
			var nfa *NFA
			for _, subpattern := range parts {
				nfa = buildNFA(subpattern)
				nfa.printNFA()
			}
		}
	}

	return nil
}
