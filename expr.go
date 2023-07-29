package flow

import (
	"errors"
	"sync"
)

// Expr represents any type of expression that can be evaluated on the given Exchange object (simple, LUA, xpath, ...etc)
type Expr interface {
	Evaluate(exchange *Exchange) (any, error)
}

type ExprFactoryFunc func(expression string, opts ...any) (Expr, error)

var exprFactoryMap = map[string]ExprFactoryFunc{}
var exprFactoryMu = sync.RWMutex{}

func RegisterExprFactory(kind string, factoryFunc ExprFactoryFunc) error {
	exprFactoryMu.Lock()
	defer exprFactoryMu.Unlock()

	_, exists := exprFactoryMap[kind]
	if exists {
		return errors.New("flow: expr factory already registered")
	}

	exprFactoryMap[kind] = factoryFunc

	return nil
}

func NewExpr(kind, expression string, opts ...any) (Expr, error) {
	exprFactoryMu.RLock()
	defer exprFactoryMu.RUnlock()

	if factoryFunc, exists := exprFactoryMap[kind]; exists {
		return factoryFunc(expression, opts)
	}

	return nil, errors.New("flow: unknown expr kind")
}
