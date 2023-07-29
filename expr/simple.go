package expr

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/paveldanilin/flow"
)

func init() {
	// TODO: add error check
	flow.RegisterExprFactory("simple", func(expression string, opts ...any) (flow.Expr, error) {
		if len(opts) > 0 {
			if asBool, ifBool := opts[0].(bool); ifBool && asBool {
				return SimpleBool(expression)
			}
		}
		return Simple(expression)
	})
	flow.RegisterExprFactory("simple:bool", func(expression string, opts ...any) (flow.Expr, error) {
		return SimpleBool(expression)
	})
}

// simpleEnv wraps Exchange with a set of helper functions
// https://expr.medv.io/docs/Getting-Started#configuration
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

// simpleExpr wraps a vm.Program
// Language definition: https://expr.medv.io/docs/Language-Definition
// https://github.com/antonmedv/expr
type simpleExpr struct {
	program *vm.Program
}

func (se *simpleExpr) Evaluate(exchange *flow.Exchange) (any, error) {
	r, err := expr.Run(se.program, simpleEnv{Exchange: exchange})
	if err != nil {
		return false, err
	}
	return r, nil
}

// Simple wraps a github.com/antonmedv/expr/vm/Program
// Language definition: https://expr.medv.io/docs/Language-Definition
// https://github.com/antonmedv/expr
func Simple(expression string) (flow.Expr, error) {
	out, err := expr.Compile(expression, expr.Env(simpleEnv{}))
	if err != nil {
		return nil, err
	}

	return &simpleExpr{program: out}, nil
}

func MustSimple(expression string) flow.Expr {
	e, err := Simple(expression)
	if err != nil {
		panic(fmt.Sprintf("simple: %s", err.Error()))
	}
	return e
}

func SimpleBool(expression string) (flow.Expr, error) {
	out, err := expr.Compile(expression, expr.Env(simpleEnv{}), expr.AsBool())
	if err != nil {
		return nil, err
	}

	return &simpleExpr{program: out}, nil
}

func MustSimpleBool(expression string) flow.Expr {
	e, err := SimpleBool(expression)
	if err != nil {
		panic(fmt.Sprintf("simple: %s", err.Error()))
	}
	return e
}
