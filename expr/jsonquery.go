package expr

import (
	"errors"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xpath"
	"github.com/paveldanilin/flow"
	"strings"
)

func init() {
	err := flow.RegisterExprFactory("jsonquery", JsonQuery)
	if err != nil {
		panic(fmt.Sprintf("flow.expr: could not register 'jsonquery' factory: %s", err.Error()))
	}
}

func JsonQuery(query string, opts ...flow.ExprOption) (flow.Expr, error) {
	var retType string
	if strings.HasSuffix(query, "#int") {
		retType = "int"
		query = strings.TrimSuffix(query, "#int")
	}

	selector, err := xpath.Compile(query)
	if err != nil {
		return nil, err
	}

	return &jsonQuery{
		querySelector: selector,
		retType:       retType,
		BaseExpr:      flow.NewBaseExpr(opts...),
	}, nil
}

func MustJsonQuery(query string, opts ...flow.ExprOption) flow.Expr {
	e, err := JsonQuery(query, opts...)
	if err != nil {
		panic(fmt.Sprintf("expr.jsonquery: %s", err.Error()))
	}
	return e
}

// jsonQuery wraps https://github.com/antchfx/jsonquery
type jsonQuery struct {
	querySelector *xpath.Expr
	retType       string
	*flow.BaseExpr
}

func (jq *jsonQuery) Evaluate(exchange *flow.Exchange) (any, error) {
	val, err := getValueForEvaluation(exchange, jq.BaseExpr)
	if err != nil {
		return nil, fmt.Errorf("expr.jsonquery: %w", err)
	}

	if strVal, isString := val.(string); isString {
		return jq.eval(strVal)
	}

	return nil, errors.New("expr.jsonquery: value must be JSON string")
}

func (jq *jsonQuery) eval(v string) (any, error) {
	if v == "" {
		return nil, errors.New("expr.jsonquery: value must be non empty JSON string")
	}

	doc, err := jsonquery.Parse(strings.NewReader(v))
	if err != nil {
		return nil, err
	}

	n := jsonquery.QuerySelector(doc, jq.querySelector)

	if jq.retType == "" {
		return n.Value(), nil
	}

	ret := n.Value()

	switch ret.(type) {
	case float64:
		if jq.retType == "int" {
			return int(ret.(float64)), nil
		}
	}

	return n.Value(), nil
}
