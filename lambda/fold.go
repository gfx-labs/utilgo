package lambda

func Foldl[T any](x T, xs []T, fx func(T, T) T) T {
	if len(xs) == 0 {
		return x
	}
	return Foldl(fx(x, xs[0]), xs[1:], fx)
}
func Foldl1[T any](xs []T, fx func(T, T) T) T {
	if len(xs) < 1 {
		return *new(T)
	}
	return Foldl(xs[0], xs[1:], fx)
}

func Foldr[T any](x T, xs []T, fx func(T, T) T) T {
	if len(xs) == 0 {
		return x
	}
	return Foldr(fx(xs[len(xs)-1], x), xs[:len(xs)-1], fx)
}

func Foldr1[T any](xs []T, fx func(T, T) T) T {
	if len(xs) < 1 {
		return *new(T)
	}
	return Foldr(xs[len(xs)-1], xs[:len(xs)-1], fx)
}
