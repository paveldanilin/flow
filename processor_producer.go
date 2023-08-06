package flow

import (
	"context"
)

type ProducerProcessor struct {
	Producer
	BaseProcessor
}

func NewProducerProcessor(producer Producer) *ProducerProcessor {
	return &ProducerProcessor{Producer: producer}
}

func (p *ProducerProcessor) Process(ctx context.Context, exchange *Exchange) error {
	err := p.Producer.Process(ctx, exchange)
	if err != nil {
		exchange.SetError(err)
		return err
	}
	return p.next(ctx, exchange)
}
