package flow

type Processor interface {
	Process(exchange *Exchange) error
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
func (bp *BaseProcessor) To(targetProcessor Processor, exchange *Exchange) error {
	if targetProcessor == nil {
		return nil
	}

	err := targetProcessor.Process(exchange)
	if err != nil {
		return err
	}
	return nil
}

func (bp *BaseProcessor) next(exchange *Exchange) error {
	return bp.To(bp.Next(), exchange)
}
