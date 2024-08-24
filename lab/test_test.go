package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"testing"
	"time"
)

func TestBehavior_Case1(t *testing.T) {
	t.Parallel()

	var tests = []lab.Test[
		lab.In[
			string,
			lab.TODO,
		],
		lab.Out[*time.Time],
	]{
		{
			Name: "#1",
			In: lab.In[string, lab.TODO]{
				Behavior: func(t *testing.T, payload string, suite lab.TODO) {
					t.Helper()
				},
			},
			Out: lab.Out[*time.Time]{},
		},
	}

	for i := range tests {
		var test = tests[i]

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
		})
	}
}
