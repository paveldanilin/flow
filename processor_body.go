package flow

import (
	"context"
	"fmt"
)

type BodyProcessor struct {
	bodyExpression Expr
	// TODO: add message type (in/out). default: in
	BaseProcessor
}

func NewBodyProcessor(bodyExpression Expr) *BodyProcessor {
	return &BodyProcessor{
		bodyExpression: bodyExpression,
	}
}

func (p *BodyProcessor) Process(ctx context.Context, exchange *Exchange) error {
	ret, err := p.bodyExpression.Evaluate(exchange)
	if err != nil {
		exchange.SetError(fmt.Errorf("body: %w", err))
		return exchange.Error()
	}

	exchange.In().SetBody(ret)

	return p.next(ctx, exchange)
}
