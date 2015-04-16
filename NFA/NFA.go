package NFA

import (
	"encoding/json"
	"fmt"
	"strings"
    "crypto/md5"
)

type NFA struct {
	begin, end *Node
	nodes      []*Node
}

func (nfa *NFA) Hash() [16]byte {
    var data string
    for _, n := range nfa.nodes {
        var r_name, l_name string
        if n.Right != nil {
            r_name = n.Right.Name
        }
        if n.Left != nil {
            l_name = n.Left.Name
        }
        data += fmt.Sprintf("%s%s%s%s%s", n.Name, n.LeftSymbol, n.RightSymbol, l_name, r_name)
    }

    return md5.Sum([]byte(data))
}

func (nfa *NFA) Length() int {
    return len(nfa.nodes)
}

func (nfa *NFA) PrintJSON() {
    fmt.Println("JSON:")

    type JSONNode struct {
        Name, LeftName, RightName, LeftSymbol, RightSymbol string
    }
    var n_array []JSONNode

    for _, node := range nfa.nodes {
        var left_name, right_name, left_symbol, right_symbol string
        if node.Left != nil {
            left_name = node.Left.Name
            left_symbol = node.LeftSymbol
        }
        if node.Right != nil {
            right_name = node.Right.Name
            right_symbol = node.RightSymbol
        }
        n_array = append(n_array, JSONNode{node.Name, left_name, right_name, left_symbol, right_symbol})
    }

    json_data, _ := json.MarshalIndent(n_array, "", "   ")
    fmt.Println(string(json_data))

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

func (nfa *NFA) PrintNFA() {
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

func BuildNFA(pattern string) *NFA {
//	fmt.Println(pattern)
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

	if len(pattern) == 2 {
		if pattern[1] == '*' {
			if pattern[0] == '(' || pattern[0] == ')' || pattern[0] == '*' || pattern[0] == '+' {
				panic("Error left symbol of pattern " + pattern)
			}
			return powerNFA(BuildNFA(string(pattern[0])))
		}
		return multiplyNFA(BuildNFA(string(pattern[0])), BuildNFA(string(pattern[1])))
	}

	if strings.Count(pattern, "(") == 0 && strings.Count(pattern, ")") == 0 {
		// ..+..
		if index := strings.Index(pattern, "+"); index > 0 {
			return addNFA(BuildNFA(pattern[:index]), BuildNFA(pattern[index+1:]))
		}

		// ...*
		if pattern[len(pattern)-1] == '*' {
			return powerNFA(BuildNFA(pattern[:len(pattern)-1]))
		}

		// abc*db
		if index := strings.Index(pattern, "*"); index > 0 {
			return multiplyNFA(BuildNFA(pattern[:index]), BuildNFA(pattern[index+1:]))
		}

		// abcd
		return multiplyNFA(BuildNFA(string(pattern[0])), BuildNFA(pattern[1:]))
	}

	// a lot of (..)..(..)..(..)
	brackets_level := 0
	left_bracket_index, right_bracket_index := -1, 0
	for i, s := range pattern {
		if s == '(' {

			brackets_level++
			if left_bracket_index == -1 {
				left_bracket_index = i
			}
		}
		if s == ')' {
			brackets_level--
			if brackets_level == 0 {
				right_bracket_index = i
			}
		}

		if brackets_level == 0 && left_bracket_index != -1 {
			break
		}
	}

	if brackets_level != 0 {
		panic("Left and rights bracket doesn't equal!")
	}

	// (..)
	if left_bracket_index == 0 && right_bracket_index == len(pattern)-1 {
		return BuildNFA(pattern[1 : len(pattern)-1])
	}

	// (..)*
	if left_bracket_index == 0 && right_bracket_index == len(pattern)-2 && pattern[len(pattern)-1] == '*' {
		return powerNFA(BuildNFA(pattern[:len(pattern)-1]))
	}

	// X(..) or X+(..)
	if left_bracket_index >= 1 {
		if pattern[left_bracket_index-1] == '+' {
			return addNFA(BuildNFA(pattern[:left_bracket_index-1]), BuildNFA(pattern[left_bracket_index:]))
		} else {
			return multiplyNFA(BuildNFA(pattern[:left_bracket_index]), BuildNFA(pattern[left_bracket_index:]))
		}
	}

	// (..)X or (..)+X or (..)*X or (..)*+X
	if left_bracket_index == 0 {
		if pattern[right_bracket_index+1] == '*' {
			right_bracket_index++
		}

		if pattern[right_bracket_index+1] == '+' {
			return addNFA(BuildNFA(pattern[:right_bracket_index+1]), BuildNFA(pattern[right_bracket_index+2:]))
		} else {
			return multiplyNFA(BuildNFA(pattern[:right_bracket_index+1]), BuildNFA(pattern[right_bracket_index+1:]))
		}
	}

	return nil
}
