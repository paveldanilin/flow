package flow

import "context"

type ConditionalProcessor struct {
	condition     Expr
	thenProcessor Processor
	elseProcessor Processor
	BaseProcessor
}

func NewConditionalProcessor(condition Expr, thenProcessor Processor, elseProcessor Processor) *ConditionalProcessor {
	return &ConditionalProcessor{
		condition:     condition,
		thenProcessor: thenProcessor,
		elseProcessor: elseProcessor,
	}
}

func (p *ConditionalProcessor) Process(ctx context.Context, exchange *Exchange) error {
	nextProcessor, err := p.nextProcessor(exchange)
	if err != nil {
		exchange.SetError(err)
		return err
	}

	err = p.To(ctx, nextProcessor, exchange)
	if err != nil {
		exchange.SetError(err)
		return err
	}

	return p.next(ctx, exchange)
}

func (p *ConditionalProcessor) Then() Processor {
	return p.thenProcessor
}

func (p *ConditionalProcessor) Else() Processor {
	return p.elseProcessor
}

func (p *ConditionalProcessor) nextProcessor(exchange *Exchange) (Processor, error) {
	ret, err := p.condition.Evaluate(exchange)
	if err != nil {
		return nil, err
	}

	if AsBool(ret) {
		return p.thenProcessor, nil
	}

	return p.elseProcessor, nil
}
