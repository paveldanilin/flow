package expr

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/paveldanilin/flow"
)

func init() {
	err := flow.RegisterExprFactory("simple", Simple)
	if err != nil {
		panic(fmt.Sprintf("flow.expr: could not register 'simple' factory: %s", err.Error()))
	}
}

// Simple wraps https://github.com/antonmedv/expr
// Language definition: https://expr.medv.io/docs/Language-Definition
func Simple(expression string, opts ...flow.ExprOption) (flow.Expr, error) {
	pg, err := expr.Compile(expression, expr.Env(simpleEnv{}))
	if err != nil {
		return nil, err
	}

	return &simpleExpr{
		program:  pg,
		BaseExpr: flow.NewBaseExpr(opts...),
	}, nil
}

// MustSimple wraps https://github.com/antonmedv/expr
// Language definition: https://expr.medv.io/docs/Language-Definition
func MustSimple(expression string, opts ...flow.ExprOption) flow.Expr {
	e, err := Simple(expression, opts...)
	if err != nil {
		panic(fmt.Sprintf("expr.simple: %s", err.Error()))
	}
	return e
}

// simpleEnv wraps Exchange with a set of helper functions
// https://expr.medv.io/docs/Getting-Started#configuration
// Example
// `exchange.Prop('a') > 0` equals to `Prop('a') == 1`
// `exchange.In().Header('user_type') == 'admin'` equals to `InHeader('user_type') == 'admin'`
// `exchange.Out().Header('flag') == true` equals to `OutHeader('flag') == true`
type simpleEnv struct {
	Exchange *flow.Exchange `expr:"exchange"`
}

func (env simpleEnv) InHeader(headerName string) any {
	return env.Exchange.In().MustHeader(headerName)
}

func (env simpleEnv) InBody() any {
	return env.Exchange.In().Body()
}

func (env simpleEnv) OutHeader(headerName string) any {
	return env.Exchange.Out().MustHeader(headerName)
}

func (env simpleEnv) OutBody() any {
	return env.Exchange.Out().Body()
}

func (env simpleEnv) Prop(propName string) any {
	return env.Exchange.MustProp(propName)
}

// simpleExpr wraps https://github.com/antonmedv/expr
// Language definition: https://expr.medv.io/docs/Language-Definition
type simpleExpr struct {
	program *vm.Program
	*flow.BaseExpr
}

func (se *simpleExpr) Evaluate(exchange *flow.Exchange) (any, error) {
	ret, err := expr.Run(se.program, simpleEnv{Exchange: exchange})
	if err != nil {
		return nil, err
	}
	return ret, nil
}
