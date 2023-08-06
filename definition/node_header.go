package definition

import "fmt"

type HeaderNode struct {
	HeaderName string
	HeaderExpr Expr
	Node
}

func (n HeaderNode) String() string {
	return fmt.Sprintf("SetHeader:[%s]=%s:%s", n.HeaderName, n.HeaderExpr.Lang, n.HeaderExpr.Expression)
}

func NewHeaderNode(headerName string, headerExpr Expr) *HeaderNode {
	return &HeaderNode{
		HeaderExpr: headerExpr,
		HeaderName: headerName,
		Node:       NewNode(),
	}
}
