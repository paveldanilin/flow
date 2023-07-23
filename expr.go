package flow

type Expr interface {
	Evaluate(exchange *Exchange) (any, error)
}

type ExprBool interface {
	Evaluate(exchange *Exchange) (bool, error)
}
