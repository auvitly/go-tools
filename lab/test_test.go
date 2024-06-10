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
		lab.Data[Arguments],
		lab.Data[bool],
	]{
		{
			Title: "#1 Equal",
			Request: lab.Data[Arguments]{
				Data: Arguments{
					Arg1: 1,
					Arg2: 1,
				},
			},
			Expected: lab.Data[bool]{
				Data: true,
			},
		},
		{
			Title: "#2 Not Equal",
			Request: lab.Data[Arguments]{
				Data: Arguments{
					Arg1: 0,
					Arg2: 1,
				},
			},
			Expected: lab.Data[bool]{
				Data: false,
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
				test.Expected.Data,
				fn(test.Request.Data.Arg1, test.Request.Data.Arg2),
			)
		})
	}
}

func TestDataError(t *testing.T) {
	t.Parallel()

	type Arguments struct {
		Arg1 int
		Arg2 int
	}

	var target = errors.New("div by zero")

	var tests = []lab.Test[
		lab.Data[Arguments],
		lab.Result[float64, error],
	]{
		{
			Title: "#1 Equal",
			Request: lab.Data[Arguments]{
				Data: Arguments{
					Arg1: 2,
					Arg2: 1,
				},
			},
			Expected: lab.Result[float64, error]{
				Data: 2,
			},
		},
		{
			Title: "#2 Not Equal",
			Request: lab.Data[Arguments]{
				Data: Arguments{
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

			result, err := fn(test.Request.Data.Arg1, test.Request.Data.Arg2)
			if err != nil {
				assert.ErrorIs(tt, err, test.Expected.Error)

				return
			}

			assert.Equal(tt, test.Expected.Data, result)
		})
	}
}
