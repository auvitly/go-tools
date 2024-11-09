package vault

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"sync"

	"testing"
)

var (
	mu    sync.Mutex
	vault = make(map[*testing.T]map[string]any)
)

// Store - stores the specified value within the active pointer to the testing object *testing.T.
// After the test is completed, the data will be deleted.
func Store[V any](t *testing.T, key string, value V) (result V) {
	t.Helper()

	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()

		delete(vault, t)
	})

	mu.Lock()
	defer mu.Unlock()

	if vault[t] == nil {
		vault[t] = make(map[string]any)
	}

	stored, ok := vault[t][key].(V)

	switch {
	case ok && !reflect.DeepEqual(stored, value):
		t.Fatalf("value with key '%s' already stored with value %#v",
			key,
			spew.Sprintf("%v", stored),
		)

		return result

	case ok && reflect.DeepEqual(stored, value):
		value = stored
	case !ok:
		vault[t][key] = value
	}

	return value
}

// Load - loading an object of type V from the testing object storage *testing.T.
func Load[V any](t *testing.T, key string) (value V) {
	t.Helper()

	mu.Lock()
	defer mu.Unlock()

	if vault[t] == nil || vault[t][key] == nil {
		t.Fatalf("not found value with key '%s'", key)
	}

	stored, ok := vault[t][key].(V)
	if !ok {
		t.Fatalf("stored value with key '%s' error: expected type %T, actual type %T", key, value, stored)
	}

	return stored
}
