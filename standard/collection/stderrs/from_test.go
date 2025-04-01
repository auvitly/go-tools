package stderrs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/stretchr/testify/require"
)

type ForRegistry struct {
	Code    int
	Message string
}

func (e ForRegistry) Error() string {
	return fmt.Sprintf("error with code %d, message %s", e.Code, e.Message)
}

func TestRegistry(t *testing.T) {
	stderrs.RegistryFrom(func(err error) *stderrs.Error {
		var my ForRegistry

		if errors.As(err, &my) {
			switch my.Code {
			case 1:
				return stderrs.Internal.SetMessage(my.Message)
			default:
				return stderrs.Unknown.SetMessage(my.Message)
			}
		}

		return nil
	})

	stderr, ok := stderrs.From(ForRegistry{Code: 1, Message: "message"})
	require.True(t, ok)
	require.True(t, stderr.Is(stderrs.Internal))

	stderr, ok = stderrs.From(ForRegistry{Code: 0, Message: "message"})
	require.True(t, ok)
	require.True(t, stderr.Is(stderrs.Unknown))
}

type FromImpl struct {
	Code    int
	Message string
}

func (e FromImpl) Error() string {
	return fmt.Sprintf("error with code %d, message %s", e.Code, e.Message)
}

func (e FromImpl) StandardFrom(err error) *stderrs.Error {
	var my FromImpl

	if errors.As(err, &my) {
		switch my.Code {
		case 1:
			return stderrs.Internal.SetMessage(my.Message)
		default:
			return stderrs.Unknown.SetMessage(my.Message)
		}
	}

	return nil
}

func TestFromImpl(t *testing.T) {
	stderr, ok := stderrs.From(FromImpl{Code: 1, Message: "message"})
	require.True(t, ok)
	require.True(t, stderr.Is(stderrs.Internal))

	stderr, ok = stderrs.From(FromImpl{Code: 0, Message: "message"})
	require.True(t, ok)
	require.True(t, stderr.Is(stderrs.Unknown))
}
