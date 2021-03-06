package NFA

import "fmt"

type Node struct {
	Id                   int
	Left, Right             *Node
	LeftSymbol, RightSymbol string
}

// for JSON encode/decode
type JSONNode struct {
    Id, Left_id, Right_id int
    LeftSymbol, RightSymbol string
}

func (n *Node) Name() string {
    return fmt.Sprintf("Node%d", n.Id+1)
}

func (n *Node) toString() string {
	left_symbol, right_symbol := n.LeftSymbol, n.RightSymbol
	if len(left_symbol) == 0 {
		left_symbol = "λ"
	}
	if len(right_symbol) == 0 {
		right_symbol = "λ"
	}

	var left_name, right_name string
	if n.Left != nil {
		left_name = n.Left.Name()
	}
	if n.Right != nil {
		right_name = n.Right.Name()
	}
	return fmt.Sprintf("%s<-%s-%s-%s->%s", left_name, left_symbol, n.Name(), right_symbol, right_name)
}
