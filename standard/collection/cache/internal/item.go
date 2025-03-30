package internal

type Item[V any] struct {
	Expirations []Expiration
	Value       V
}
