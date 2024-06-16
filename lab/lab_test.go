package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockData struct {
	Input  string
	Output time.Time
}

var behavior = &lab.Behavior[MockData, map[string]*time.Time]{
	Data: []MockData{
		{Input: "key_1", Output: lab.Now},
		{Input: "key_2", Output: lab.Now.Add(time.Hour)},
	},
	Setter: func(t *testing.T, ctrl map[string]*time.Time, data MockData) {
		ctrl[data.Input] = lab.Pointer(data.Output)

		t.Cleanup(func() { delete(ctrl, data.Input) })
	},
}

func TestBehavior(t *testing.T) {
	t.Parallel()

	var tests = []lab.Test[string, *time.Time]{
		{
			Title:       "#1 TestKV",
			Preparation: []lab.Preparation{behavior},
			Input:       "key_1",
			Output:      lab.Pointer(lab.Now),
		},
		{
			Title:       "#2 TestKV",
			Preparation: []lab.Preparation{behavior},
			Input:       "key_2",
			Output:      lab.Pointer(lab.Now.Add(time.Hour)),
		},
		{
			Title:       "#3 TestKV",
			Preparation: []lab.Preparation{behavior},
			Input:       "key_3",
			Output:      nil,
		},
	}

	for i := range tests {
		var test = tests[i]

		t.Run(test.Title, func(tt *testing.T) {
			var ctrl = make(map[string]*time.Time)

			test.Prepare(tt, ctrl)

			assert.Equal(tt, test.Output, ctrl[test.Input])
		})
	}
}
