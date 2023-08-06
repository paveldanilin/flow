package timer

import (
	"github.com/paveldanilin/flow"
	"net/url"
	"sync"
	"time"
)

func init() {
	flow.RegisterComponent("timer", &Component{
		timersMap: map[string]*Timer{},
	})
}

type Component struct {
	timersMap map[string]*Timer
	mu        sync.Mutex
}

// timer:my-timer?interval=5s
func (c *Component) GetConsumer(componentURI string) (flow.Consumer, error) {
	u, err := url.Parse(componentURI)
	if err != nil {
		return nil, err
	}

	interval, err := time.ParseDuration(u.Query().Get("interval"))
	if err != nil {
		return nil, err
	}

	timerId := u.Opaque
	timer := c.createOrGetTimer(timerId, interval)

	return NewConsumer(timer), nil
}

func (c *Component) GetProducer(componentURI string) (flow.Producer, error) {
	panic("timer: producer is not implemented")
}

func (c *Component) Start() error {
	for _, timer := range c.timersMap {
		timer.Start()
	}
	return nil
}

func (c *Component) Stop() {
	for _, timer := range c.timersMap {
		timer.Stop()
	}
}

func (c *Component) createOrGetTimer(timerID string, interval time.Duration) *Timer {
	c.mu.Lock()
	timer, exists := c.timersMap[timerID]
	if exists {
		c.mu.Unlock()
		return timer
	}
	timer = NewTimer(timerID, interval)
	c.timersMap[timerID] = timer
	c.mu.Unlock()
	return timer
}
