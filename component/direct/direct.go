package direct

import (
	"github.com/google/uuid"
	"github.com/paveldanilin/flow"
)

type Component struct {
}

func (c *Component) GetConsumer(params map[string]any) flow.Consumer {
	return nil
}

func (c *Component) GetProducer(params map[string]any) flow.Producer {
	return nil
}

type Consumer struct {
	processor flow.Processor
}

func NewConsumer(processor flow.Processor) *Consumer {
	return &Consumer{
		processor: processor,
	}
}

func (c *Consumer) Start(opts flow.ConsumerOpts) error {
	exchange := flow.GetExchange()
	exchange.SetProps(opts.ExchangeProperties)
	exchange.SetProp("DIRECT_RUN_ID", uuid.NewString())
	exchange.In().SetHeaders(opts.MessageHeaders)
	exchange.In().SetBody(opts.MessageBody)
	err := c.processor.Process(exchange)
	flow.ReleaseExchange(exchange)
	return err
}

func (c *Consumer) Stop() {

}
