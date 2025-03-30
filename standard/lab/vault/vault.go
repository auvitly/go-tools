package vault

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

var vault = struct {
	mu      sync.Mutex
	storage map[*testing.T]map[string]any
}{
	storage: make(map[*testing.T]map[string]any),
}

// Store - stores the specified value within the active pointer to the testing object *testing.T.
// After the test is completed, the data will be deleted.
func Store[V any](t *testing.T, key string, value V) (result V) {
	if t == nil {
		panic(fmt.Sprintf("*testing.T is nil for key=%s, value=%v", key, value))
	}

	vault.mu.Lock()
	defer vault.mu.Unlock()

	if vault.storage[t] == nil {
		vault.storage[t] = make(map[string]any)
	}

	stored, ok := vault.storage[t][key].(V)

	switch {
	case ok && !reflect.DeepEqual(stored, value):
		panic(fmt.Sprintf("value with key '%s' already stored with value %#v",
			key,
			fmt.Sprintf("%v", stored),
		))
	case ok && reflect.DeepEqual(stored, value):
		value = stored
	case !ok:
		vault.storage[t][key] = value
	}

	return value
}

// Load - loading an object of type V from the testing object storage *testing.T.
func Load[V any](t *testing.T, key string) (value V) {
	vault.mu.Lock()
	defer vault.mu.Unlock()

	if vault.storage[t] == nil || vault.storage[t][key] == nil {
		panic(fmt.Sprintf(
			"not found value with key '%s'",
			key,
		))
	}

	stored, ok := vault.storage[t][key].(V)
	if !ok {
		panic(fmt.Sprintf(
			"stored value with key '%s' error: expected type %T, actual type %T",
			key, value, stored,
		))
	}

	return stored
}
