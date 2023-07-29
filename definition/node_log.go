package definition

import "fmt"

type LogNode struct {
	Message string
	Node
}

func (n LogNode) String() string {
	return fmt.Sprintf("Log:%s", n.Message)
}

func NewLogNode(message string) *LogNode {
	return &LogNode{
		Message: message,
		Node:    NewNode(),
	}
}
