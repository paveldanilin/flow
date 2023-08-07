package flow

import (
	"fmt"
	"sync"
)

// Expr represents any type of expression that can be evaluated on the given Exchange object (simple, LUA, xpath, ...etc)
type Expr interface {
	Evaluate(exchange *Exchange) (any, error)
}

type BaseExpr struct {
	messageType MessageType
	headerName  string
}

func NewBaseExpr(opts ...ExprOption) *BaseExpr {
	be := &BaseExpr{messageType: MessageIn, headerName: ""}
	for _, opt := range opts {
		opt(be)
	}
	return be
}

func (be BaseExpr) MessageType() MessageType {
	return be.messageType
}

func (be BaseExpr) HeaderName() string {
	return be.headerName
}

// ExprOption represents an option for expression.
// If HeaderName option is empty, the expression will be evaluated on the message body.
type ExprOption func(*BaseExpr)

func InMessage() ExprOption {
	return func(expr *BaseExpr) {
		expr.messageType = MessageIn
	}
}

func OutMessage() ExprOption {
	return func(expr *BaseExpr) {
		expr.messageType = MessageOut
	}
}

func HeaderName(headerName string) ExprOption {
	return func(expr *BaseExpr) {
		expr.headerName = headerName
	}
}

type ExprFactoryFunc func(expression string, opts ...ExprOption) (Expr, error)

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

func newExpr(lang, expression string, opts ...ExprOption) (Expr, error) {
	exprFactoryMu.RLock()
	defer exprFactoryMu.RUnlock()

	if factoryFunc, exists := exprFactoryMap[lang]; exists {
		return factoryFunc(expression, opts...)
	}

	return nil, fmt.Errorf("flow: unknown expr lang '%s'", lang)
}
