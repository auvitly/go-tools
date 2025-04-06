package reflector_test

import (
	"net/http"
	"testing"

	"github.com/auvitly/go-tools/standard/utils/reflector"
	"github.com/auvitly/go-tools/test/lab"
	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	var tests = []lab.Test[lab.Any, lab.Any]{
		{
			Name: "int",
			In:   0,
			Out:  reflector.Nil[int](),
		},
		{
			Name: "*int",
			In:   (*int)(nil),
			Out:  reflector.Nil[*int](),
		},
		{
			Name: "string",
			In:   "",
			Out:  reflector.Nil[string](),
		},
		{
			Name: "http.Client",
			In:   http.Client{},
			Out:  reflector.Nil[http.Client](),
		},
	}

	for i := range tests {
		var test = tests[i]

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.In, test.Out)
		})
	}

}
