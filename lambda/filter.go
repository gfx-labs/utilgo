package lambda

func Filter[T any](xs []T, fx func(T) bool) []T {
	n := 0
	for _, x := range xs {
		if fx(x) {
			xs[n] = x
			n++
		}
	}
	return xs[:n]
}
