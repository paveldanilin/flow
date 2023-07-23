package flow

import "fmt"

type LogProcessor struct {
	message string
	BaseProcessor
}

func NewLogProcessor(message string) *LogProcessor {
	return &LogProcessor{
		message: message,
	}
}

func (p *LogProcessor) Process(exchange *Exchange) error {
	println("-------------------------------------------------")
	fmt.Printf("-- EXCHANGE[%s]\n", exchange.ExchangeID())
	fmt.Printf("-- MSG: %s\n", p.message)
	fmt.Printf("-- PROPS: %v\n", exchange.Props())
	if exchange.in == nil {
		fmt.Printf("-- IN: NIL\n")
	} else {
		fmt.Printf("-- IN.HEADERS: %v\n", exchange.In().Headers())
		fmt.Printf("-- IN.BODY: %v\n", exchange.In().Body())
	}

	println("-------------------------------------------------")

	return p.next(exchange)
}
