package lab_test

import (
	"context"
	"github.com/auvitly/go-tools/lab"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Data struct {
	Input  string
	Output *time.Time
}

func TestBehavior_Case1(t *testing.T) {
	t.Parallel()

	var preparations = lab.Preparations{
		&lab.Behavior[Data, map[string]*time.Time]{
			Data: []Data{
				{Input: "key_1", Output: lab.Pointer(lab.Now)},
				{Input: "key_2", Output: lab.Pointer(lab.Now.Add(time.Hour))},
			},
			Setter: func(t *testing.T, s map[string]*time.Time, data Data) {
				s[data.Input] = data.Output

				t.Cleanup(func() { delete(s, data.Input) })
			},
		},
	}

	var test = lab.Test[string, *time.Time]{
		{
			Name:   "#1 TestKV",
			Input:  "key_1",
			Output: lab.Pointer(lab.Now),
		},
		{
			Name:   "#2 TestKV",
			Input:  "key_2",
			Output: lab.Pointer(lab.Now.Add(time.Hour)),
		},
		{
			Name:   "#3 TestKV",
			Input:  "key_3",
			Output: nil,
		},
	}

	var s = make(map[string]*time.Time)

	preparations.Init(t, s)

	test.Run(t, func(tt *testing.T, exp lab.Experiment[string, *time.Time]) {
		tt.Parallel()

		assert.Equal(tt, exp.Output, s[exp.Input])
	})
}

func TestBehavior_Case2(t *testing.T) {
	t.Parallel()

	var preparation = func(data []Data) lab.Preparation {
		return &lab.Behavior[Data, *context.Context]{
			Data: data,
			Setter: func(t *testing.T, ctx *context.Context, data Data) {
				*ctx = lab.NewContext().WithValue(data.Input, data.Output).Context
			},
		}
	}

	var test = lab.Test[string, any]{
		{
			Preparations: lab.Preparations{
				preparation([]Data{{Input: "key_1", Output: lab.Pointer(lab.Now)}}),
			},
			Name:   "#1 TestKV",
			Input:  "key_1",
			Output: lab.Pointer(lab.Now),
		},
		{
			Preparations: lab.Preparations{
				preparation([]Data{{Input: "key_2", Output: lab.Pointer(lab.Now.Add(time.Hour))}}),
			},
			Name:   "#2 TestKV",
			Input:  "key_2",
			Output: lab.Pointer(lab.Now.Add(time.Hour)),
		},
		{
			Name:   "#3 TestKV",
			Input:  "key_3",
			Output: nil,
		},
	}

	test.Run(t, func(tt *testing.T, exp lab.Experiment[string, any]) {
		tt.Parallel()

		var ctx = context.Background()

		exp.Init(tt, &ctx)

		assert.Equal(tt, exp.Output, ctx.Value(exp.Input))
	})
}
