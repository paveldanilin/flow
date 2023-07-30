package flow

import (
	"context"
	"fmt"
)

type HeaderProcessor struct {
	setExpression Expr
	headerName    string
	// TODO: add message type (in/out). default: in
	BaseProcessor
}

func NewHeaderProcessor(headerName string, setExpression Expr) *HeaderProcessor {
	return &HeaderProcessor{
		setExpression: setExpression,
		headerName:    headerName,
	}
}

func (p *HeaderProcessor) Process(ctx context.Context, exchange *Exchange) error {
	ret, err := p.setExpression.Evaluate(exchange)
	if err != nil {
		exchange.SetError(fmt.Errorf("header[%s]: %w", p.headerName, err))
		return exchange.Error()
	}

	exchange.In().SetHeader(p.headerName, ret)

	return p.next(ctx, exchange)
}
