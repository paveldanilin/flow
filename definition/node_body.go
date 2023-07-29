package definition

import "fmt"

type BodyNode struct {
	ExprKind   string
	Expression string
	Node
}

func (n BodyNode) String() string {
	return fmt.Sprintf("SetBody:%s:%s", n.ExprKind, n.Expression)
}

func NewBodyNode(exprKind, expression string) *BodyNode {
	return &BodyNode{
		ExprKind:   exprKind,
		Expression: expression,
		Node:       NewNode(),
	}
}
