package generic

func Swap[T any](xs []T) func(i, j int) {
	return func(i, j int) {
		xs[i], xs[j] = xs[j], xs[i]
	}
}
