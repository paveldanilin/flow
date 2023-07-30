package flow

import (
	"context"
	"fmt"
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

func (r *Registry) Add(flowDef *definition.Flow) error {
	r.flowMu.Lock()
	defer r.flowMu.Unlock()

	if _, flowExists := r.flowMap[flowDef.FlowID]; flowExists {
		return fmt.Errorf("flow.registry: flow with ID '%s' already registered", flowDef.FlowID)
	}

	rtFlow, err := Compile(flowDef)
	if err != nil {
		return err
	}

	r.flowMap[flowDef.FlowID] = rtFlow

	return nil
}

func (r *Registry) Execute(ctx context.Context, params Params) (any, error) {
	// TODO: queue (for online processing)

	r.flowMu.RLock()
	flow, flowExists := r.flowMap[params.FlowID]
	if !flowExists {
		r.flowMu.RUnlock()
		return nil, fmt.Errorf("flow.registry: flow with ID '%s' not found", params.FlowID)
	}
	r.flowMu.RUnlock()

	exchange, shouldRelease := r.getExchange()
	if shouldRelease {
		defer r.releaseExchange(exchange)
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

func (r *Registry) Send(params Params) error {
	// TODO: queue (for event processing)
	return nil
}

func (r *Registry) getExchange() (*Exchange, bool) {
	obj, pooled := r.exchangePool.Get()
	return obj.(*Exchange), pooled
}

func (r *Registry) releaseExchange(exchange *Exchange) {
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
	r.exchangePool.Put(exchange)
}
