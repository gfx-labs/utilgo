package lambda

func Entries[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	vals := make([]V, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return keys, vals
}

func Flatten[T any](xss [][]T) []T {
	return Foldl1(xss, Merge[T])
}

func Merge[T any](a, b []T) []T {
	return append(a, b...)
}
