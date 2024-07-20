package kit

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Behavior[P, S any] func(t *testing.T, payload P, suite S)

type In[P any, Suite any] struct {
	Payload  P
	Behavior Behavior[P, Suite]
}

func (i *In[Payload, Suite]) Init(t *testing.T, suite Suite) {
	require.NotNil(t, i.Behavior, "behavior not exists")

	i.Behavior(t, i.Payload, suite)
}
