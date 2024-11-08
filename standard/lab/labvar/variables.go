package labvar

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var (
	// Timestamp - timestamp within the current locale.
	// The method for generating timestamps is assumed relative to this object.
	// If necessary, you can use the Add method to define time boundaries.
	Timestamp = time.Now().Local()
	// Error - a random error returned dependency that does not need to be checked.
	Error = errors.New(fmt.Sprintf("unexpected error with id '%s'", uuid.New()))
)

// Pointer - returns pointer on copy value.
func Pointer[T any](v T) *T {
	return &v
}
