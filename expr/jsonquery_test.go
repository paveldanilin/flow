package expr

import (
	"github.com/paveldanilin/flow"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonQuery_Evaluate(t *testing.T) {
	jq := MustJsonQuery("person/age#int")

	exchange := flow.NewExchange()
	exchange.In().SetBody(`{
            "person":{
               "name":"John",
               "age":31,
               "female":false,
               "city":null,
               "hobbies":[
                  "coding",
                  "eating",
                  "football"
               ]
            }
         }`)

	ret, err := jq.Evaluate(exchange)

	assert.NoError(t, err)
	assert.Equal(t, 31, ret)
}
