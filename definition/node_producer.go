package definition

import "fmt"

type ProducerNode struct {
	URI string
	Node
}

func (n ProducerNode) String() string {
	return fmt.Sprintf("Producer: %s", n.URI)
}

func NewProducerNode(uri string) *ProducerNode {
	return &ProducerNode{
		URI:  uri,
		Node: NewNode(),
	}
}
