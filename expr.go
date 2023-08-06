package flow

import (
	"fmt"
	"sync"
)

// Expr represents any type of expression that can be evaluated on the given Exchange object (simple, LUA, xpath, ...etc)
type Expr interface {
	Evaluate(exchange *Exchange) (any, error)
}

type ExprFactoryFunc func(expression string, opts ...any) (Expr, error)

var (
	exprFactoryMap = map[string]ExprFactoryFunc{}
	exprFactoryMu  = sync.RWMutex{}
)

func RegisterExprFactory(lang string, factoryFunc ExprFactoryFunc) error {
	exprFactoryMu.Lock()
	defer exprFactoryMu.Unlock()

	_, factoryExists := exprFactoryMap[lang]
	if factoryExists {
		return fmt.Errorf("flow: expr factory already registered for lang='%s'", lang)
	}

	exprFactoryMap[lang] = factoryFunc

	return nil
}

func NewExpr(lang, expression string, opts ...any) (Expr, error) {
	exprFactoryMu.RLock()
	defer exprFactoryMu.RUnlock()

	if factoryFunc, exists := exprFactoryMap[lang]; exists {
		return factoryFunc(expression, opts)
	}

	return nil, fmt.Errorf("flow: unknown expr lang='%s'", lang)
}
