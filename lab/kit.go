package lab

import (
	"encoding/json"
	"time"
)

// Now - application start timestamp. Can be used in functions, models.
var Now = time.Now()

// Pointer - return a pointer to copy value.
func Pointer[V any](value V) *V {
	return &value
}

// MustMarshalJSON - returns the result of executing the json.Marshal function, panics in case of an error.
func MustMarshalJSON(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return data
}

// MustUnmarshalJSON - returns the result of executing the json.Unmarshal function, panics in case of an error.
func MustUnmarshalJSON[R any](data []byte) (result R) {
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}

	return result
}
