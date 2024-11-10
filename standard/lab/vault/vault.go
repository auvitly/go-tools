package vault

import (
	"fmt"
	"github.com/auvitly/go-tools/afmt"
	"reflect"
	"sync"
)

type Vault struct {
	mu      sync.Mutex
	storage map[string]any
}

func New() *Vault {
	return &Vault{
		storage: make(map[string]any),
	}
}

// Store - stores the specified value within the active pointer to the testing object *testing.T.
// After the test is completed, the data will be deleted.
func Store[V any](vault *Vault, key string, value V) (result V) {
	vault.mu.Lock()
	defer vault.mu.Unlock()

	stored, ok := vault.storage[key].(V)

	switch {
	case ok && !reflect.DeepEqual(stored, value):
		panic(fmt.Sprintf("value with key '%s' already stored with value %#v",
			key,
			afmt.Sprintf("%v", stored),
		))
	case ok && reflect.DeepEqual(stored, value):
		value = stored
	case !ok:
		vault.storage[key] = value
	}

	return value
}

// Load - loading an object of type V from the testing object storage *testing.T.
func Load[V any](vault *Vault, key string) (value V) {
	vault.mu.Lock()
	defer vault.mu.Unlock()

	if vault == nil || vault.storage[key] == nil {
		panic(fmt.Sprintf(
			"not found value with key '%s'",
			key,
		))
	}

	stored, ok := vault.storage[key].(V)
	if !ok {
		panic(fmt.Sprintf(
			"stored value with key '%s' error: expected type %T, actual type %T",
			key, value, stored,
		))
	}

	return stored
}
