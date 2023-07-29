package flow

import (
	"errors"
	"github.com/paveldanilin/flow/definition"
)

func Compile(flowDef *definition.Flow) (*Flow, error) {
	startProcessor, err := compile(flowDef.Root.Children())
	if err != nil {
		return nil, err
	}

	return New(flowDef.FlowID, nil, startProcessor), nil
}

func compile(nodes []definition.Node) (Processor, error) {
	var startProcessor Processor
	var curProcessor Processor
	var err error

	for _, nodeDef := range nodes {
		if curProcessor == nil {
			curProcessor, err = compileNode(nodeDef)
			if err != nil {
				return nil, err
			}
			startProcessor = curProcessor
		} else {
			next, err := compileNode(nodeDef)
			if err != nil {
				return nil, err
			}
			curProcessor.SetNext(next)
			curProcessor = next
		}
	}

	return startProcessor, nil
}

func compileNode(node definition.Node) (Processor, error) {
	switch node.(type) {
	case *definition.ConditionNode:
		return createConditionalProcessor(node.(*definition.ConditionNode))

	case *definition.HeaderNode:
		return createHeaderProcessor(node.(*definition.HeaderNode))

	case *definition.BodyNode:
		return createBodyProcessor(node.(*definition.BodyNode))

	case *definition.LogNode:
		return createLogProcessor(node.(*definition.LogNode))

	default:
		return nil, errors.New("flow compiler: unknown node definition")
	}
}

func createConditionalProcessor(condDef *definition.ConditionNode) (*ConditionalProcessor, error) {
	expr, err := NewExpr(condDef.Kind, condDef.Expression)
	if err != nil {
		return nil, err
	}

	thenDst, err := compile(condDef.MustChild(0).Children())
	if err != nil {
		return nil, err
	}

	elseDst, err := compile(condDef.MustChild(1).Children())
	if err != nil {
		return nil, err
	}

	return NewConditionalProcessor(expr, thenDst, elseDst), nil
}

func createHeaderProcessor(headerDef *definition.HeaderNode) (*HeaderProcessor, error) {
	expr, err := NewExpr(headerDef.ExprKind, headerDef.Expression)
	if err != nil {
		return nil, err
	}

	return NewHeaderProcessor(headerDef.HeaderName, expr), nil
}

func createBodyProcessor(bodyDef *definition.BodyNode) (*BodyProcessor, error) {
	expr, err := NewExpr(bodyDef.ExprKind, bodyDef.Expression)
	if err != nil {
		return nil, err
	}

	return NewBodyProcessor(expr), nil
}

func createLogProcessor(logDef *definition.LogNode) (*LogProcessor, error) {
	return NewLogProcessor(logDef.Message), nil
}
