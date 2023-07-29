package definition

import "fmt"

type ConditionNode struct {
	Kind       string
	Expression string
	Node
}

func (n ConditionNode) String() string {
	return fmt.Sprintf("Condition:%s:%s", n.Kind, n.Expression)
}

func NewConditionalNode(kind, expression string) *ConditionNode {
	return &ConditionNode{
		Kind:       kind,
		Expression: expression,
		Node:       NewNode(),
	}
}
