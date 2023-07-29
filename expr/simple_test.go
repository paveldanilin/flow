package expr

import (
	"github.com/paveldanilin/flow"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleBool_Prop(t *testing.T) {
	expression := MustSimpleBool("Prop('username') == 'test'")

	exchange := flow.NewExchange()
	exchange.SetProp("username", "test")

	ret, err := expression.Evaluate(exchange)

	assert.NoError(t, err)
	assert.Equal(t, true, ret)
}

func TestSimple_InHeader(t *testing.T) {
	expression := MustSimple("InHeader('is_admin') == true")

	exchange := flow.NewExchange()
	exchange.In().SetHeader("is_admin", true)

	ret, err := expression.Evaluate(exchange)

	assert.NoError(t, err)
	assert.Equal(t, true, ret)
}

func TestSimple_Evaluate(t *testing.T) {
	expression := MustSimple("exchange.Prop('a') > 100")

	exchange := flow.NewExchange()
	exchange.SetProps(map[string]any{"a": 101})

	ret, err := expression.Evaluate(exchange)

	assert.NoError(t, err)
	assert.Equal(t, true, ret)
}
