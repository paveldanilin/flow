package main

import (
	"context"
	"fmt"
	"github.com/paveldanilin/flow"
	"github.com/paveldanilin/flow/definition"
	_ "github.com/paveldanilin/flow/expr" // <- simple, simple:bool expr
)

func main() {
	userFlow := definition.NewBuilder().
		FlowID("abcd").
		// TODO: expression notation - <kind>:<expression> (simple:1,simple:InHeader('a')>1)
		SetHeader("a", "simple", "1").
		SetHeader("b", "simple", "10").
		SetBody("simple", "InHeader('a') + InHeader('b')").
		Condition("InBody() == 2").
		Then().Log("OK!").End().
		Else().Log("NOK!").End().
		End().
		GetFlow()

	println(definition.Dump(userFlow.Root))

	flowRegistry := flow.NewRegistry(flow.RegistryConfig{ExchangePoolSize: 1000})

	flowRegistry.Add(userFlow)

	ret, err := flowRegistry.Execute(context.TODO(), flow.Params{
		FlowID: "abcd",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("->%v\n", ret)

	//consumer := direct.NewConsumer(p1)
	//consumer := timer.NewConsumer(30 * time.Second)
	//consumer.AddProcessor(p1)

	//f := flow.New("my-test-flow", consumer)

	//err := f.Start(map[string]any{"a": 1}, nil)
	//if err != nil {
	//	panic(err)
	//}
}
