package flow

import "fmt"

type BodyProcessor struct {
	setExpression Expr
	BaseProcessor
}

func NewBodyProcessor(setExpression Expr) *BodyProcessor {
	return &BodyProcessor{
		setExpression: setExpression,
	}
}

func (p *BodyProcessor) Process(exchange *Exchange) error {
	res, err := p.setExpression.Evaluate(exchange)
	if err != nil {
		exchange.SetError(fmt.Errorf("body: %w", err))
		return exchange.Error()
	}

	exchange.In().SetBody(res)

	return p.next(exchange)
}
