package flow

import "context"

type Processor interface {
	Process(ctx context.Context, exchange *Exchange) error
	Next() Processor
	SetNext(nextProcessor Processor) Processor
}

type BaseProcessor struct {
	nextProcessor Processor
}

func (bp *BaseProcessor) Next() Processor {
	return bp.nextProcessor
}

func (bp *BaseProcessor) SetNext(nextProcessor Processor) Processor {
	bp.nextProcessor = nextProcessor
	return bp.nextProcessor
}

// To sends a processing flow to the target processor.
func (bp *BaseProcessor) To(ctx context.Context, targetProcessor Processor, exchange *Exchange) error {
	if targetProcessor == nil {
		return nil
	}

	err := targetProcessor.Process(ctx, exchange)
	if err != nil {
		return err
	}
	return nil
}

func (bp *BaseProcessor) next(ctx context.Context, exchange *Exchange) error {
	return bp.To(ctx, bp.Next(), exchange)
}
