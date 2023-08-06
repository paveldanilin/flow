package flow

type Params struct {
	FlowID         string
	ExchangeProps  map[string]any
	MessageHeaders map[string]any
	MessageBody    any
}

type Setter interface {
	SetFlow(flow *Flow)
}

type Flow struct {
	id        string
	consumer  Consumer
	processor Processor
	registry  *Registry
}

func New(id string, consumer Consumer, processor Processor) *Flow {
	return &Flow{
		id:        id,
		consumer:  consumer,
		processor: processor,
	}
}

func (f *Flow) FlowID() string {
	return f.id
}

func (f *Flow) Consumer() Consumer {
	return f.consumer
}

func (f *Flow) Processor() Processor {
	return f.processor
}

func (f *Flow) Registry() *Registry {
	return f.registry
}
