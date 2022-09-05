package ifx

type Cloner[T any] interface {
	Clone() T
}

type CloneFunc[T any] func() T

func (c CloneFunc[T]) Clone() T {
	return c()
}
