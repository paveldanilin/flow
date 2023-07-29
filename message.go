package flow

import (
	"fmt"
	"reflect"
)

type Message struct {
	headers map[string]any
	body    interface{}
}

func NewMessage() *Message {
	return &Message{
		headers: make(map[string]any),
		body:    nil,
	}
}

func (m *Message) Header(name string) (any, bool) {
	val, ok := m.headers[name]
	if ok {
		return val, true
	}
	return nil, false
}

func (m *Message) MustHeader(name string) any {
	v, ok := m.Header(name)
	if ok {
		return v
	}
	panic(fmt.Sprintf("message: header not found '%s'", name))
}

func (m *Message) SetHeader(name string, value any) {
	m.headers[name] = value
}

func (m *Message) SetHeaders(headers map[string]any) {
	if headers == nil {
		headers = map[string]any{}
	}
	m.headers = headers
}

func (m *Message) HasHeader(name string) bool {
	_, ok := m.headers[name]
	return ok
}

func (m *Message) Headers() map[string]any {
	return m.headers
}

func (m *Message) HeaderKeys() []string {
	var keys []string
	for key := range m.headers {
		keys = append(keys, key)
	}
	return keys
}

func (m *Message) Body() any {
	return m.body
}

func (m *Message) SetBody(body any) {
	m.body = body
}

// String returns a Message body as a string
func (m *Message) String() string {
	if s, ok := m.body.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", m.body)
}

func (m *Message) BodyType() string {
	// TODO: move to prop
	t := reflect.TypeOf(m.body).Elem()
	return t.PkgPath() + "." + t.Name()
}
