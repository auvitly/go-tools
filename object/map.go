package object

type Map interface {
	implMap()
}

type MapOf[K comparable, V any] map[K]V

func (m MapOf[K, V]) implMap() {}
func (m MapOf[K, V]) Len() int { return len(m) }
