package lambda

func MapNil[T any](xs []T, fx func(T)) {
	for _, v := range xs {
		fx(v)
	}
}

func Map[T any](xs []T, fx func(T) T) []T {
	for i, v := range xs {
		xs[i] = fx(v)
	}
	return xs
}

func MapV[T, V any](xs []T, fx func(T) V) []V {
	ov := make([]V, len(xs))
	for i, v := range xs {
		ov[i] = fx(v)
	}
	return ov
}

func MapError[T any](xs []T, fx func(T) (T, error)) ([]T, []error) {
	oe := make([]error, len(xs))
	for i, v := range xs {
		xs[i], oe[i] = fx(v)
	}
	return xs, oe
}

func MapErrorV[T, V any](xs []T, fx func(T) (V, error)) ([]V, []error) {
	ov := make([]V, len(xs))
	oe := make([]error, len(xs))
	for i, v := range xs {
		ov[i], oe[i] = fx(v)
	}
	return ov, oe
}
