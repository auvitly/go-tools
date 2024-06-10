package lab_test

import (
	"errors"
	"github.com/auvitly/go-tools/lab"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestData(t *testing.T) {
	t.Parallel()

	type Arguments struct {
		Arg1 int
		Arg2 int
	}

	var tests = []lab.Test[
		lab.Payload[Arguments],
		lab.Payload[bool],
	]{
		{
			Title: "#1 Equal",
			Request: lab.Payload[Arguments]{
				Payload: Arguments{
					Arg1: 1,
					Arg2: 1,
				},
			},
			Expected: lab.Payload[bool]{
				Payload: true,
			},
		},
		{
			Title: "#2 Not Equal",
			Request: lab.Payload[Arguments]{
				Payload: Arguments{
					Arg1: 0,
					Arg2: 1,
				},
			},
			Expected: lab.Payload[bool]{
				Payload: false,
			},
		},
	}

	var fn = func(a, b int) bool { return a == b }

	for i := range tests {
		var test = tests[i]

		t.Run(test.Title, func(tt *testing.T) {
			tt.Parallel()

			assert.Equal(
				tt,
				test.Expected.Payload,
				fn(test.Request.Payload.Arg1, test.Request.Payload.Arg2),
			)
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	type Arguments struct {
		Arg1 int
		Arg2 int
	}

	var target = errors.New("div by zero")

	var tests = []lab.Test[
		lab.Payload[Arguments],
		lab.Result[float64, error],
	]{
		{
			Title: "#1 Equal",
			Request: lab.Payload[Arguments]{
				Payload: Arguments{
					Arg1: 2,
					Arg2: 1,
				},
			},
			Expected: lab.Result[float64, error]{
				Payload: 2,
			},
		},
		{
			Title: "#2 Not Equal",
			Request: lab.Payload[Arguments]{
				Payload: Arguments{
					Arg1: 1,
					Arg2: 0,
				},
			},
			Expected: lab.Result[float64, error]{
				Error: target,
			},
		},
	}

	var fn = func(a, b int) (float64, error) {
		if b == 0 {
			return 0, target
		}

		return float64(a) / float64(b), nil
	}

	for i := range tests {
		var test = tests[i]

		t.Run(test.Title, func(tt *testing.T) {
			tt.Parallel()

			result, err := fn(test.Request.Payload.Arg1, test.Request.Payload.Arg2)
			if err != nil {
				assert.ErrorIs(tt, err, test.Expected.Error)

				return
			}

			assert.Equal(tt, test.Expected.Payload, result)
		})
	}
}
