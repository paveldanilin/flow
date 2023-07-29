package timer

import (
	"github.com/paveldanilin/flow"
	"time"
)

type Consumer struct {
	id         string
	interval   time.Duration
	processors []flow.Processor
	done       chan bool
}

func NewConsumer(id string, interval time.Duration) *Consumer {
	return &Consumer{
		id:         id,
		interval:   interval,
		processors: []flow.Processor{},
		done:       make(chan bool),
	}
}

func (c *Consumer) AddProcessor(processor flow.Processor) {
	c.processors = append(c.processors, processor)
}

func (c *Consumer) Start(opts flow.ConsumerOpts) error {
	ticker := time.NewTicker(c.interval)
	var tickCounter int64

	//go func() {
	for {
		select {
		case <-c.done:
			break
		case t := <-ticker.C:
			for _, processor := range c.processors {
				p := processor
				go func() {
					exchange := flow.GetExchange()
					exchange.SetProps(opts.ExchangeProperties)
					exchange.SetProp("TIMER_FIRED_TIME", t)
					tickCounter++
					exchange.SetProp("TIMER_TICK_ID", tickCounter)
					p.Process(exchange)
					flow.ReleaseExchange(exchange)
				}()
			}
		}
	}
	//}()

	return nil
}

func (c *Consumer) Stop() {
	c.done <- true
}
