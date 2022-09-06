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
	return Foldr(xss, Merge[T])
}

func Merge[T any](a, b []T) []T {
	return append(a, b...)
}

//split xs into n groups
//extra elements will be put in the last array to keep order
//if len(xs) < n, result will be [][]T{xs}
func Split[T any](xs []T, n int) [][]T {
	if n == 0 {
		return nil
	}
	groupSize := len(xs) / n
	if groupSize == 0 {
		return [][]T{xs}
	}
	idx := 0
	out := make([][]T, n)
	for i := 0; i < len(xs); i = i + groupSize {
		from := i
		to := i + groupSize
		if idx == (n - 1) {
			to = len(xs)
		}
		if to >= len(xs) {
			to = len(xs)
			out[idx] = append(out[idx], xs[from:to]...)
			return out
		}
		out[idx] = append(out[idx], xs[from:to]...)
		idx = idx + 1
	}
	return out
}
