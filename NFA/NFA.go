package NFA

import (
	"encoding/json"
	"fmt"
	"strings"
    "crypto/md5"
    "bytes"
    "os"
    "io/ioutil"
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
            r_name = n.Right.Name()
        }
        if n.Left != nil {
            l_name = n.Left.Name()
        }
        data += fmt.Sprintf("%s%s%s%s%s", n.Name, n.LeftSymbol, n.RightSymbol, l_name, r_name)
    }

    return md5.Sum([]byte(data))
}

func (nfa *NFA) Nodes() []*Node {
    return nfa.nodes
}

func (nfa *NFA) Length() int {
    return len(nfa.nodes)
}

func (nfa *NFA) SaveToFile(filename string) {
    file, err := os.Create(filename)
    if err != nil {
        fmt.Println("Error when creating file: "+err.Error())
        return
    }

    _, err = file.Write(nfa.JSON())
    if err != nil {
        fmt.Println("Error when writing to file: "+err.Error())
    }
    fmt.Println("NFA successfully saved to the file.")
    file.Close()
}

func (nfa *NFA) JSON() []byte {
    var n_array []JSONNode

    for _, node := range nfa.nodes {
        var left_symbol, right_symbol string
        var left_id, right_id int
        if node.Left != nil {
            left_id = node.Left.Id
            left_symbol = node.LeftSymbol
        }
        if node.Right != nil {
            right_id = node.Right.Id
            right_symbol = node.RightSymbol
        }
        n_array = append(n_array, JSONNode{node.Id, left_id, right_id, left_symbol, right_symbol})
    }

    json_data, _ := json.Marshal(n_array)
    return json_data
}

func (nfa *NFA) PrintJSON() {
    var pretty_json bytes.Buffer
    json.Indent(&pretty_json, nfa.JSON(), "", "     ")
    fmt.Println(pretty_json.String())
}


func NewFromFile(filename string) *NFA {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        panic("Error when opening file: "+err.Error())
    }
    return NewFromJSON(file)
}

func NewFromJSON(input_data []byte) *NFA {
    var n_array []JSONNode
    err := json.Unmarshal(input_data, &n_array)
    if err != nil {
        panic("Error when parsing NFA file. "+err.Error())
    }

    var nodes []*Node
    for _, node := range n_array {
        new_node := Node{LeftSymbol: node.LeftSymbol, RightSymbol: node.RightSymbol}
        nodes = append(nodes, &new_node)
    }

    for i, node := range n_array {
        if node.Left_id > 0 {
            nodes[i].Left = nodes[node.Left_id]
        }
        if node.Right_id > 0 {
            nodes[i].Right = nodes[node.Right_id]
        }
    }

    return NewNFA(nodes...)
}





func (nfa *NFA) addNode(new_node *Node) {
    new_node.Id = len(nfa.nodes)
	nfa.nodes = append(nfa.nodes, new_node)
}

func (nfa *NFA) addNodes(new_nodes []*Node) {
    var start_i = len(nfa.nodes)
    for i, nnode := range new_nodes {
        nnode.Id = start_i+i
    }
	nfa.nodes = append(nfa.nodes, new_nodes...)
}

func NewNFA(nodes ...*Node) *NFA {
	var nfa NFA
	nfa.addNodes(nodes)
	for i, v := range nodes {
		v.Id = i
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
            return multiplyNFA(BuildNFA(pattern[:len(pattern)-2]), BuildNFA(pattern[len(pattern)-2:]))
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