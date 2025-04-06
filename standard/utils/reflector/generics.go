package reflector

func Nil[T any]() T { return *new(T) }
