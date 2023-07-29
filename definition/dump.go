package definition

import (
	"fmt"
	"strings"
)

const (
	newLine      = "\n"
	emptySpace   = "     "
	middleItem   = "├──> "
	continueItem = "│    "
	lastItem     = "└──> "
)

func Dump(node Node) string {
	return printNode(node, []bool{}, 0)
}

func printNode(node Node, spaces []bool, index int) string {
	var result string

	last := node.Parent() == nil
	if !last && index == len(node.Parent().Children())-1 {
		last = true
	}
	spacesChild := append(spaces, last)

	result += printText(fmt.Sprintf("%v", node), spaces, last)

	for i, childNode := range node.Children() {
		last := i == len(node.Children())-1
		spacesChild = append(spaces, last)
		result += printNode(childNode, spacesChild, i)
	}

	return result
}

func printText(text string, spaces []bool, last bool) string {
	var result string
	for _, space := range spaces {
		if space {
			result += emptySpace
		} else {
			result += emptySpace
		}
	}

	indicator := middleItem
	if last {
		indicator = lastItem
	}

	var out string
	lines := strings.Split(text, "\n")
	for i := range lines {
		text := lines[i]
		if i == 0 {
			out += result + indicator + text + newLine
			continue
		}
		if last {
			indicator = emptySpace
		} else {
			indicator = continueItem
		}
		out += result + indicator + text + newLine
	}

	return out
}
