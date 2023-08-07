package expr

import (
	"fmt"
	"github.com/paveldanilin/flow"
)

func getValueForEvaluation(exchange *flow.Exchange, baseExpr *flow.BaseExpr) (any, error) {
	var msg *flow.Message
	if baseExpr.MessageType() == flow.MessageIn {
		msg = exchange.In()
	} else {
		msg = exchange.Out()
	}

	var val any
	if baseExpr.HeaderName() == "" {
		val = msg.Body()
	} else {
		v, exists := msg.Header(baseExpr.HeaderName())
		if !exists {
			return nil, fmt.Errorf("[%s] header not found", baseExpr.HeaderName())
		}
		val = v
	}

	return val, nil
}
