package ifx

type Pool[T any] interface {
	Get() T
	TryGet() (T, bool)

	Put(T)
}
