package random

import (
	crand "crypto/rand"
	"math/rand"
)

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func Integer[V interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}]() (v V) {
	return V(rand.Uint64())
}

func Float[V interface {
	~float32 | ~float64
}]() (v V) {
	return V(rand.Float64())
}

func Bytes(n int) []byte {
	var buf = make([]byte, n)

	_, err := crand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}

func String(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	str := string(b)

	return str
}
