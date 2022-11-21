package lambda

func Filter[T any](fx func(T) bool) func([]T) []T {
	return func(xs []T) []T {
		n := 0
		for _, x := range xs {
			if fx(x) {
				xs[n] = x
				n++
			}
		}
		return xs[:n]
	}
}
func FilterIdx[T any](fx func(T, int) bool) func([]T) []T {
	return func(xs []T) []T {
		n := 0
		for idx, x := range xs {
			if fx(x, idx) {
				xs[n] = x
				n++
			}
		}
		return xs[:n]
	}
}
func FilterIdxXs[T any](fx func(T, int, []T) bool) func([]T) []T {
	return func(xs []T) []T {
		n := 0
		for idx, x := range xs {
			if fx(x, idx, xs) {
				xs[n] = x
				n++
			}
		}
		return xs[:n]
	}
}

func Every[T any](fx func(T) bool) func([]T) bool {
	return func(xs []T) bool {
		for _, x := range xs {
			if !fx(x) {
				return false
			}
		}
		return true
	}
}
func EveryIdx[T any](fx func(T, int) bool) func([]T) bool {
	return func(xs []T) bool {
		for idx, x := range xs {
			if !fx(x, idx) {
				return false
			}
		}
		return true
	}
}
func EveryIdxXs[T any](fx func(T, int, []T) bool) func([]T) bool {
	return func(xs []T) bool {
		for idx, x := range xs {
			if !fx(x, idx, xs) {
				return false
			}
		}
		return true
	}
}
