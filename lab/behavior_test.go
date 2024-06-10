package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/stretchr/testify/assert"
	"testing"
)

type KV struct {
	Key   string
	Value string
}

var behavior = lab.Calls[KV, map[string]string]{
	Data: []KV{
		{
			Key:   "key_1",
			Value: "value_1",
		},
		{
			Key:   "key_2",
			Value: "value_2",
		},
	},
	Setter: func(ctrl map[string]string, data KV) {
		ctrl[data.Key] = data.Value
	},
}

var testsBehavior = []lab.Test[
	lab.Payload[string],
	lab.Payload[string],
]{
	{
		Title: "#1 TestKV",
		Behavior: []lab.Behavior{
			&behavior,
		},
		Request: lab.Payload[string]{
			Payload: "key_1",
		},
		Expected: lab.Payload[string]{
			Payload: "value_1",
		},
	},
	{
		Title: "#2 TestKV",
		Behavior: []lab.Behavior{
			&behavior,
		},
		Request: lab.Payload[string]{
			Payload: "key_2",
		},
		Expected: lab.Payload[string]{
			Payload: "value_2",
		},
	},
}

func TestBehavior(t *testing.T) {
	t.Parallel()

	for i := range testsBehavior {
		var test = testsBehavior[i]

		t.Run(test.Title, func(tt *testing.T) {
			var ctrl = make(map[string]string)

			test.ApplyBehavior(ctrl)

			assert.Equal(tt, test.Expected.Payload, ctrl[test.Request.Payload])
		})
	}
}
