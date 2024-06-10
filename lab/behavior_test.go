package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBehavior(t *testing.T) {
	t.Parallel()

	type KV struct {
		Key   string
		Value string
	}

	var setter = func(ctrl map[string]string, data KV) {
		ctrl[data.Key] = data.Value
	}

	var tests = []lab.Test[
		lab.Payload[string],
		lab.Payload[string],
	]{
		{
			Title: "TestKV",
			Behavior: []lab.Behavior{
				&lab.Call[KV, map[string]string]{
					Data: KV{
						Key:   "key",
						Value: "value",
					},
					Setter: setter,
				},
			},
			Request: lab.Payload[string]{
				Payload: "key",
			},
			Expected: lab.Payload[string]{
				Payload: "value",
			},
		},
	}

	for i := range tests {
		var test = tests[i]

		t.Run(test.Title, func(tt *testing.T) {
			var storage = make(map[string]string)

			test.SetBehavior(storage)

			assert.Equal(tt, test.Expected.Payload, storage[test.Request.Payload])
		})
	}
}
