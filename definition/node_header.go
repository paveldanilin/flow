package definition

import "fmt"

type HeaderNode struct {
	ExprKind   string
	Expression string
	HeaderName string
	Node
}

func (n HeaderNode) String() string {
	return fmt.Sprintf("SetHeader:[%s]=%s:%s", n.HeaderName, n.ExprKind, n.Expression)
}

func NewHeaderNode(exprKind, expression, headerName string) *HeaderNode {
	return &HeaderNode{
		ExprKind:   exprKind,
		Expression: expression,
		HeaderName: headerName,
		Node:       NewNode(),
	}
}
