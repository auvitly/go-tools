package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/auvitly/go-tools/lab/labvar"
	"testing"
	"time"
)

func TestBehavior_Case1(t *testing.T) {
	t.Parallel()

	var tests = lab.Tests[
		lab.In[string, lab.TODO],
		lab.Out[*time.Time, lab.Empty],
	]{
		{
			Name: "#1",
			In: lab.In[string, lab.TODO]{
				Request: labvar.Timestamp.Format(time.RFC3339),
				Behavior: func(t *testing.T) {
					t.Helper()
				},
			},
			Out: lab.Out[*time.Time, lab.Empty]{
				Response: labvar.Pointer(labvar.Timestamp),
			},
		},
	}

	tests.Run(t, func(t *testing.T, test lab.Test[lab.In[string, lab.TODO], lab.Out[*time.Time, lab.Empty]]) {
		t.Parallel()
	})
}
