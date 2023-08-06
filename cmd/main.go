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

	//ret, err := flowRegistry.Execute(context.TODO(), flow.Params{
	//	FlowID: "abcd",
	//})

	//if err != nil {
	//	panic(err)
	//}

	//fmt.Printf("->%v\n", ret)
	select {}

	flowRegistry.Stop()

	//consumer := direct.NewConsumer(p1)
	//consumer := timer.NewConsumer(30 * time.Second)
	//consumer.AddProcessor(p1)

	//f := flow.New("my-test-flow", consumer)

	//err := f.Start(map[string]any{"a": 1}, nil)
	//if err != nil {
	//	panic(err)
	//}
}
