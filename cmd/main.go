package main

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/paveldanilin/flow"
	"time"
)

type BoolExpr struct {
	p *vm.Program
}

func (be *BoolExpr) Evaluate(exchange *flow.Exchange) (bool, error) {
	r, err := expr.Run(be.p, exchange)
	if err != nil {
		return false, err
	}
	return r.(bool), nil
}

func createBoolExpr(expression string) flow.ExprBool {
	out, err := expr.Compile(expression, expr.AsBool())
	if err != nil {
		panic(err)
	}
	return &BoolExpr{p: out}
}

type Expr struct {
	p *vm.Program
}

func (be *Expr) Evaluate(exchange *flow.Exchange) (any, error) {
	r, err := expr.Run(be.p, exchange)
	if err != nil {
		return false, err
	}
	return r, nil
}

func createExpr(expression string) flow.Expr {
	out, err := expr.Compile(expression)
	if err != nil {
		panic(err)
	}
	return &Expr{p: out}
}

func main() {
	p1 := flow.NewLogProcessor("START")

	p2 := flow.NewConditionalProcessor(createBoolExpr("Prop('a')>1"), flow.NewLogProcessor("THEN"), flow.NewLogProcessor("ELSE"))
	p2.Then().SetNext(flow.NewLogProcessor("AA"))
	p2.Else().SetNext(flow.NewLogProcessor("BB")).SetNext(flow.NewDelayProcessor(5 * time.Second)).SetNext(flow.NewLogProcessor("Z"))
	p2.SetNext(flow.NewLogProcessor("FINISH")).
		SetNext(flow.NewDelayProcessor(1 * time.Second)).
		SetNext(flow.NewHeaderProcessor("zZz", createExpr("123456"))).
		SetNext(flow.NewBodyProcessor(createExpr("Prop('a') + 4"))).
		SetNext(flow.NewLogProcessor("N"))

	p1.SetNext(p2)

	exchangePool := flow.NewObjectPool(2)
	exchangePool.Init(func() interface{} {
		return flow.NewExchange(nil)
	})

	start(p1, exchangePool)
}

func start(p flow.Processor, exchangePool *flow.ObjectPool) {
	ex := exchangePool.Get().(*flow.Exchange)
	ex.SetProp("a", 1)

	p.Process(ex)

	exchangePool.Put(ex)
}
