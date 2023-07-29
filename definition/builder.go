package definition

import (
	"fmt"
)

// Builder a flow definition builder
type Builder struct {
	flowID string
	parent *Builder
	root   Node
	stack  []Node
}

func NewBuilder() *Builder {
	return newWithRootBuilder(nil, nil)
}

func newWithRootBuilder(parent *Builder, root Node) *Builder {
	b := &Builder{
		parent: parent,
		root:   nil,
		stack:  []Node{},
	}

	if root == nil {
		if parent == nil {
			root = NewNode()
		} else {
			root = parent.current()
		}
	}

	b.pushNode(root)
	b.root = root

	return b
}

func (b *Builder) FlowID(flowID string) *Builder {
	b.flowID = flowID
	return b
}

func (b *Builder) Log(message string) *Builder {
	b.current().Add(NewLogNode(message))
	return b
}

func (b *Builder) SetHeader(headerName, exprKind, expression string) *Builder {
	b.current().Add(NewHeaderNode(exprKind, expression, headerName))
	return b
}

func (b *Builder) SetBody(exprKind, expression string) *Builder {
	b.current().Add(NewBodyNode(exprKind, expression))
	return b
}

func (b *Builder) Condition(expression string) *Builder {
	b.beginTree(NewConditionalNode("simple:bool", expression))

	conditionBuilder := newWithRootBuilder(b, nil)
	conditionBuilder.ChildNode(NewNode()) // Then arc
	conditionBuilder.ChildNode(NewNode()) // Else arc

	return conditionBuilder
}

func (b *Builder) Then() *Builder {
	_, isConditionalNode := b.current().(*ConditionNode)
	if isConditionalNode {
		return newWithRootBuilder(b, b.current().MustChild(0))
	}
	panic(fmt.Sprintf("definition: cannot start the 'Then' flow, since the current current is not a conditional current"))
}

func (b *Builder) Else() *Builder {
	_, isConditionalNode := b.current().(*ConditionNode)
	if isConditionalNode {
		return newWithRootBuilder(b, b.current().MustChild(1))
	}
	panic(fmt.Sprintf("definition: cannot start the 'Else' flow, since the current current is not a conditional current"))
}

func (b *Builder) End() *Builder {
	if b.parent == nil {
		b.endTree()
		return b
	}

	b.parent.endTree()

	return b.parent
}

func (b *Builder) ChildNode(node Node) *Builder {
	b.current().Add(node)
	return b
}

func (b *Builder) GetFlow() *Flow {
	return &Flow{
		FlowID: b.flowID,
		Root:   b.root,
	}
}

func (b *Builder) current() Node {
	return b.stack[len(b.stack)-1]
}

func (b *Builder) beginTree(node Node) *Builder {
	b.current().Add(node)
	b.pushNode(node)
	return b
}

func (b *Builder) endTree() *Builder {
	b.popNode()
	return b
}

func (b *Builder) pushNode(node Node) {
	b.stack = append(b.stack, node)
}

func (b *Builder) popNode() {
	if len(b.stack) > 1 {
		b.stack = b.stack[:len(b.stack)-1]
	}
}
