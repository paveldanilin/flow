package definition

import "fmt"

type BodyNode struct {
	BodyExpr Expr
	Node
}

func (n BodyNode) String() string {
	return fmt.Sprintf("SetBody:%s:%s", n.BodyExpr.Lang, n.BodyExpr.Expression)
}

func NewBodyNode(bodyExpr Expr) *BodyNode {
	return &BodyNode{
		BodyExpr: bodyExpr,
		Node:     NewNode(),
	}
}
