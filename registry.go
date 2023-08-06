package flow

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

func (r *Registry) Start() error {
	err := startComponents()
	if err != nil {
		return err
	}

	for _, f := range r.flowMap {
		err := f.Consumer().Start()
		if err != nil {
			fmt.Printf("[%s] failed to start %s", f.FlowID(), err.Error())
		}
	}

	return nil
}

func (r *Registry) Stop() {
	stopComponents()
}

func (r *Registry) Add(flow *definition.Flow) error {
	r.flowMu.Lock()
	defer r.flowMu.Unlock()

	if _, flowExists := r.flowMap[flow.FlowID]; flowExists {
		return fmt.Errorf("flow.registry: flow already exists ID='%s'", flow.FlowID)
	}

	rtFlow, err := Compile(flow)
	if err != nil {
		return err
	}

	rtFlow.registry = r

	r.flowMap[flow.FlowID] = rtFlow

	return nil
}

func (r *Registry) Execute(ctx context.Context, params Params) (any, error) {
	// TODO: queue (for online processing)
	if ctx == nil {
		ctx = context.Background()
	}

	r.flowMu.RLock()
	flow, flowExists := r.flowMap[params.FlowID]
	if !flowExists {
		r.flowMu.RUnlock()
		return nil, fmt.Errorf("flow.registry: flow not found ID='%s'", params.FlowID)
	}
	r.flowMu.RUnlock()

	exchange, shouldRelease := r.GetExchange()
	if shouldRelease {
		defer r.ReleaseExchange(exchange)
	}

	if params.ExchangeProps != nil {
		exchange.SetProps(params.ExchangeProps)
	}
	exchange.SetProp("FLOW_ID", params.FlowID)
	exchange.SetProp("FLOW_RUN_ID", uuid.NewString())

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

func (r *Registry) Send(ctx context.Context, params Params) error {
	r.flowMu.RLock()
	flow, flowExists := r.flowMap[params.FlowID]
	if !flowExists {
		r.flowMu.RUnlock()
		return fmt.Errorf("flow.registry: flow not found ID='%s'", params.FlowID)
	}
	r.flowMu.RUnlock()

	exchange, shouldRelease := r.GetExchange()
	if shouldRelease {
		defer r.ReleaseExchange(exchange)
	}

	if params.ExchangeProps != nil {
		exchange.SetProps(params.ExchangeProps)
	}
	exchange.SetProp("FLOW_ID", params.FlowID)
	exchange.SetProp("FLOW_RUN_ID", uuid.NewString())

	if params.MessageHeaders != nil {
		exchange.In().SetHeaders(params.MessageHeaders)
	}

	exchange.In().SetBody(params.MessageBody)

	flow.Processor().Process(ctx, exchange)

	return nil
}

func (r *Registry) GetExchange() (*Exchange, bool) {
	obj, pooled := r.exchangePool.Get()
	return obj.(*Exchange), pooled
}

func (r *Registry) ReleaseExchange(exchange *Exchange) {
	// TODO: Reset exchange
	exchange.flowRegistry = nil
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
