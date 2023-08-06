package definition

import "fmt"

type ConditionalNode struct {
	Condition Expr
	Node
}

func (n ConditionalNode) String() string {
	return fmt.Sprintf("Condition:(%s:%s)", n.Condition.Lang, n.Condition.Expression)
}

func NewConditionalNode(condition Expr) *ConditionalNode {
	return &ConditionalNode{
		Condition: condition,
		Node:      NewNode(),
	}
}
