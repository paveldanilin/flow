package flow

import (
	"context"
	"errors"
	"github.com/paveldanilin/flow/definition"
	"sync"
)

type RegistryConfig struct {
	ExchangePoolSize int
}

type Registry struct {
	flowMap      map[string]*Flow
	flowMu       sync.RWMutex
	exchangePool *ObjectPool
}

func NewRegistry(cfg RegistryConfig) *Registry {
	ctx := &Registry{
		flowMap:      map[string]*Flow{},
		exchangePool: NewObjectPool(cfg.ExchangePoolSize),
	}

	ctx.exchangePool.Init(func() interface{} {
		return NewExchange()
	})

	return ctx
}

func (c *Registry) Add(flowDef *definition.Flow) error {
	c.flowMu.Lock()
	defer c.flowMu.Unlock()

	if _, flowExists := c.flowMap[flowDef.FlowID]; flowExists {
		return errors.New("flow: already registered")
	}

	rtFlow, err := Compile(flowDef)
	if err != nil {
		return err
	}

	c.flowMap[flowDef.FlowID] = rtFlow

	return nil
}

func (c *Registry) Execute(ctx context.Context, params Params) (any, error) {
	// TODO: queue (for online processing)

	c.flowMu.RLock()
	flow, flowExists := c.flowMap[params.FlowID]
	if !flowExists {
		c.flowMu.RUnlock()
		return nil, errors.New("flow: not found")
	}
	c.flowMu.RUnlock()

	exchange, shouldRelease := c.getExchange()
	if shouldRelease {
		defer c.releaseExchange(exchange)
	}

	if params.ExchangeProps != nil {
		exchange.SetProps(params.ExchangeProps)
	}
	if params.MessageHeaders != nil {
		exchange.In().SetHeaders(params.MessageHeaders)
	}
	exchange.In().SetBody(params.MessageBody)

	err := flow.Processor().Process(ctx, exchange)
	if err != nil {
		return nil, err
	}

	var ret any
	if exchange.out == nil {
		ret = exchange.in.body
	} else {
		ret = exchange.out.body
	}

	return ret, nil
}

func (c *Registry) Send(params Params) error {
	// TODO: queue (for event processing)
	return nil
}

func (c *Registry) getExchange() (*Exchange, bool) {
	obj, pooled := c.exchangePool.Get()
	return obj.(*Exchange), pooled
}

func (c *Registry) releaseExchange(exchange *Exchange) {
	// TODO: Reset exchange
	exchange.flowContext = nil
	exchange.props = map[string]any{}
	exchange.err = nil
	exchange.in.headers = map[string]any{}
	exchange.in.body = nil
	if exchange.out != nil {
		exchange.out.headers = map[string]any{}
		exchange.out.body = nil
	}
	c.exchangePool.Put(exchange)
}
