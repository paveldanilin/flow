package flow

import "fmt"

type HeaderProcessor struct {
	setExpression Expr
	headerName    string
	BaseProcessor
}

func NewHeaderProcessor(headerName string, setExpression Expr) *HeaderProcessor {
	return &HeaderProcessor{
		setExpression: setExpression,
		headerName:    headerName,
	}
}

func (p *HeaderProcessor) Process(exchange *Exchange) error {
	res, err := p.setExpression.Evaluate(exchange)
	if err != nil {
		exchange.SetError(fmt.Errorf("header[%s]: %w", p.headerName, err))
		return exchange.Error()
	}

	exchange.In().SetHeader(p.headerName, res)

	return p.next(exchange)
}
