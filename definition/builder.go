package definition

import (
	"fmt"
)

// Builder a flow definition builder
type Builder struct {
	flowID      string
	consumerURI string
	parent      *Builder
	root        Node
	stack       []Node
}

func NewBuilder() *Builder {
	return newBuilder(nil, nil)
}

func newBuilder(parent *Builder, root Node) *Builder {
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

func (b *Builder) AddNode(node Node) *Builder {
	b.current().Add(node)
	return b
}

func (b *Builder) FlowID(flowID string) *Builder {
	b.flowID = flowID
	return b
}

func (b *Builder) Consumer(uri string) *Builder {
	b.consumerURI = uri
	return b
}

func (b *Builder) To(uri string) *Builder {
	b.current().Add(NewProducerNode(uri))
	return b
}

func (b *Builder) Log(message string) *Builder {
	b.current().Add(NewLogNode(message))
	return b
}

func (b *Builder) SetHeader(headerName string, headerValue Expr) *Builder {
	b.current().Add(NewHeaderNode(headerName, headerValue))
	return b
}

func (b *Builder) SetBody(bodyValue Expr) *Builder {
	b.current().Add(NewBodyNode(bodyValue))
	return b
}

func (b *Builder) Condition(condition Expr) *Builder {
	b.beginTree(NewConditionalNode(condition))

	conditionBuilder := newBuilder(b, nil)
	conditionBuilder.AddNode(NewNode()) // Then arc
	conditionBuilder.AddNode(NewNode()) // Else arc

	return conditionBuilder
}

func (b *Builder) Then() *Builder {
	_, isConditionalNode := b.current().(*ConditionalNode)
	if isConditionalNode {
		return newBuilder(b, b.current().MustChild(0))
	}
	panic(fmt.Sprintf("flow.definition.builder: cannot begin 'Then' flow, since the current node is not a conditional"))
}

func (b *Builder) Else() *Builder {
	_, isConditionalNode := b.current().(*ConditionalNode)
	if isConditionalNode {
		return newBuilder(b, b.current().MustChild(1))
	}
	panic(fmt.Sprintf("flow.definition.builder: cannot begin the 'Else' flow, since the current node is not a conditional"))
}

func (b *Builder) End() *Builder {
	if b.parent == nil {
		b.endTree()
		return b
	}

	b.parent.endTree()

	return b.parent
}

func (b *Builder) GetFlow() *Flow {
	return &Flow{
		FlowID:      b.flowID,
		ConsumerURI: b.consumerURI,
		Root:        b.root,
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
