package main

import (
	"github.com/paveldanilin/flow"
	_ "github.com/paveldanilin/flow/component/timer"
	"github.com/paveldanilin/flow/definition"
	_ "github.com/paveldanilin/flow/expr" // <- simple, simple:bool expr
)

var simple = definition.Simple

func main() {
	userFlow := definition.NewBuilder().
		FlowID("abcd").
		Consumer("timer:abcd?interval=5s").
		Log("BZZZZ").
		GetFlow()

	println(definition.Dump(userFlow.Root))

	flowRegistry := flow.NewRegistry(flow.RegistryConfig{ExchangePoolSize: 1000})

	err := flowRegistry.Add(userFlow)
	if err != nil {
		panic(err)
	}

	err = flowRegistry.Start()
	if err != nil {
		panic(err)
	}

	select {}

	flowRegistry.Stop()
}
