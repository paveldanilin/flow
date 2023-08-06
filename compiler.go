package flow

import (
	"fmt"
	"github.com/paveldanilin/flow/definition"
)

func Compile(flow *definition.Flow) (*Flow, error) {
	consumer, err := getConsumer(flow.ConsumerURI)
	if err != nil {
		return nil, err
	}

	processor, err := compileNodes(flow.Root.Children())
	if err != nil {
		return nil, err
	}

	f := New(flow.FlowID, consumer, processor)

	if setter, isSetter := consumer.(Setter); isSetter {
		setter.SetFlow(f)
	}

	return f, nil
}

func compileNodes(nodes []definition.Node) (Processor, error) {
	var startProcessor Processor
	var curProcessor Processor
	var err error

	for _, node := range nodes {
		if curProcessor == nil {
			curProcessor, err = createProcessor(node)
			if err != nil {
				return nil, err
			}
			startProcessor = curProcessor
		} else {
			next, err := createProcessor(node)
			if err != nil {
				return nil, err
			}
			curProcessor.SetNext(next)
			curProcessor = next
		}
	}

	return startProcessor, nil
}

func createProcessor(node definition.Node) (Processor, error) {
	switch node.(type) {
	case *definition.ConditionalNode:
		return createConditionalProcessor(node.(*definition.ConditionalNode))

	case *definition.HeaderNode:
		return createHeaderProcessor(node.(*definition.HeaderNode))

	case *definition.BodyNode:
		return createBodyProcessor(node.(*definition.BodyNode))

	case *definition.LogNode:
		return createLogProcessor(node.(*definition.LogNode))

	case *definition.ProducerNode:
		return createProducerProcessor(node.(*definition.ProducerNode))

	default:
		return nil, fmt.Errorf("flow.compiler: unknown node '%T'", node)
	}
}

func createConditionalProcessor(node *definition.ConditionalNode) (*ConditionalProcessor, error) {
	expr, err := NewExpr(node.Condition.Lang, node.Condition.Expression)
	if err != nil {
		return nil, err
	}

	thenDst, err := compileNodes(node.MustChild(0).Children())
	if err != nil {
		return nil, err
	}

	elseDst, err := compileNodes(node.MustChild(1).Children())
	if err != nil {
		return nil, err
	}

	return NewConditionalProcessor(expr, thenDst, elseDst), nil
}

func createHeaderProcessor(node *definition.HeaderNode) (*HeaderProcessor, error) {
	expr, err := NewExpr(node.HeaderExpr.Lang, node.HeaderExpr.Expression)
	if err != nil {
		return nil, err
	}

	return NewHeaderProcessor(node.HeaderName, expr), nil
}

func createBodyProcessor(node *definition.BodyNode) (*BodyProcessor, error) {
	expr, err := NewExpr(node.BodyExpr.Lang, node.BodyExpr.Expression)
	if err != nil {
		return nil, err
	}

	return NewBodyProcessor(expr), nil
}

func createLogProcessor(node *definition.LogNode) (*LogProcessor, error) {
	return NewLogProcessor(node.Message), nil
}

func createProducerProcessor(node *definition.ProducerNode) (*ProducerProcessor, error) {
	producer, err := getProducer(node.URI)
	if err != nil {
		return nil, err
	}
	return NewProducerProcessor(producer), nil
}
