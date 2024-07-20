package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/auvitly/go-tools/lab/kit"
	"testing"
	"time"
)

func TestBehavior_Case1(t *testing.T) {
	t.Parallel()

	var test = lab.Test[
		kit.In[kit.Pair[string, int], kit.TODO],
		kit.Out[*time.Time],
	]{
		{
			Name: "#1",
			In: kit.In[kit.Pair[string, int], kit.TODO]{
				Behavior: func(t *testing.T, payload kit.Pair[string, int], suite kit.TODO) {},
			},
			Out: kit.Out[*time.Time]{},
		},
	}

	test.Run(t, func(t *testing.T, test lab.Experiment[kit.In[kit.Pair[string, int], kit.TODO], kit.Out[*time.Time]]) {

	})
}
