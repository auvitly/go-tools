package lab

import (
	"testing"
)

type Behavior[S any] func(t *testing.T, suite S)

type In[P any, Suite any] struct {
	Payload  P
	Behavior Behavior[Suite]
}
