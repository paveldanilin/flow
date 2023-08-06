package timer

import (
	"context"
	"github.com/paveldanilin/flow"
	"time"
)

type Consumer struct {
	flow  *flow.Flow
	timer *Timer
}

func NewConsumer(timer *Timer) *Consumer {
	return &Consumer{
		timer: timer,
	}
}

func (c *Consumer) SetFlow(flow *flow.Flow) {
	c.flow = flow
}

func (c *Consumer) Start() error {
	c.timer.AddTask(c.flow.FlowID(), func(timerID string, firedAt time.Time, counter int64) {
		c.flow.Registry().Send(context.Background(), flow.Params{
			FlowID: c.flow.FlowID(),
			ExchangeProps: map[string]any{
				"TIMER_ID":      timerID,
				"TIMER_TIME":    firedAt,
				"TIMER_COUNTER": counter,
			},
		})
	})
	return nil
}

func (c *Consumer) Stop() {
	c.timer.RemoveTask(c.flow.FlowID())
}
