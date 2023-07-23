package flow

import (
	"fmt"
	"github.com/google/uuid"
)

type Exchange struct {
	id    string
	err   error
	in    *Message
	out   *Message
	props map[string]any
}

func NewExchange(in *Message) *Exchange {
	return &Exchange{
		id:    uuid.NewString(),
		err:   nil,
		in:    in,
		out:   nil,
		props: make(map[string]any),
	}
}

func (e *Exchange) ExchangeID() string {
	return e.id
}

func (e *Exchange) In() *Message {
	if e.in == nil {
		e.in = NewMessage()
	}
	return e.in
}

func (e *Exchange) Out() *Message {
	if e.out == nil {
		e.out = NewMessage()
	}
	return e.out
}

func (e *Exchange) SetError(err error) {
	e.err = err
}

func (e *Exchange) Error() error {
	return e.err
}

func (e *Exchange) IsError() bool {
	return e.err != nil
}

func (e *Exchange) SetProp(name string, value any) {
	e.props[name] = value
}

func (e *Exchange) SetProps(props map[string]any) {
	e.props = props
}

func (e *Exchange) Prop(name string) (any, bool) {
	if v, ok := e.props[name]; ok {
		return v, true
	}
	return nil, false
}

func (e *Exchange) MustProp(name string) any {
	v, ok := e.Prop(name)
	if ok {
		return v
	}
	panic(fmt.Sprintf("exchange: prop not found '%s'", name))
}

func (e *Exchange) Props() map[string]any {
	return e.props
}
