package lambda

func Foldr[T any](xs []T, fx func(T, T) T) T {
	if len(xs) < 1 {
		return *new(T)
	}
	if len(xs) == 1 {
		return xs[0]
	}
	for i := 1; i < len(xs); i++ {
		xs[0] = fx(xs[0], xs[i])
	}
	return xs[0]
}

func Foldl[T any](xs []T, fx func(T, T) T) T {
	if len(xs) < 1 {
		return *new(T)
	}
	if len(xs) == 1 {
		return xs[0]
	}
	li := len(xs) - 1
	for i := li - 1; i >= 0; i-- {
		xs[li] = fx(xs[li], xs[i])
	}
	return xs[li]
}
