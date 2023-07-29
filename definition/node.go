package definition

import (
	"fmt"
)

// Node a flow node definition
type Node interface {
	Parent() Node
	Children() []Node
	LastChild() Node
	FirstChild() Node
	MustChild(index int) Node
	Add(child Node) Node
	SetParent(parent Node)
	HasChildren() bool
	NoChildren() bool
}

type flowNode struct {
	parent   Node
	children []Node
}

func NewNode() Node {
	return &flowNode{
		parent:   nil,
		children: []Node{},
	}
}

func (n *flowNode) Parent() Node {
	return n.parent
}

func (n *flowNode) Children() []Node {
	return n.children
}

func (n *flowNode) LastChild() Node {
	s := len(n.children)
	if s == 0 {
		return nil
	}
	return n.children[s-1]
}

func (n *flowNode) FirstChild() Node {
	s := len(n.children)
	if s == 0 {
		return nil
	}
	return n.children[0]
}

func (n *flowNode) MustChild(index int) Node {
	if index > len(n.children) {
		panic(fmt.Sprintf("flow: a child node not found at index [%d]", index))
	}
	return n.children[index]
}

func (n *flowNode) HasChildren() bool {
	return len(n.children) > 0
}

func (n *flowNode) NoChildren() bool {
	return !n.HasChildren()
}

func (n *flowNode) Add(child Node) Node {
	n.children = append(n.children, child)
	child.SetParent(n)
	return n
}

func (n *flowNode) SetParent(parent Node) {
	n.parent = parent
}
