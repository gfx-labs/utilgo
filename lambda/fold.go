package lambda

func Foldl[T any](fx func(T, T) T) func(x T, xs []T) T {
	return func(x T, xs []T) T {
		if len(xs) == 0 {
			return x
		}
		return Foldl(fx)(fx(x, xs[0]), xs[1:])
	}
}
func Foldl1[T any](fx func(T, T) T) func(xs []T) T {
	return func(xs []T) T {
		if len(xs) < 1 {
			return *new(T)
		}
		return Foldl(fx)(xs[0], xs[1:])
	}
}

func Foldr[T any](fx func(T, T) T) func(x T, xs []T) T {
	return func(x T, xs []T) T {
		if len(xs) == 0 {
			return x
		}
		return Foldr(fx)(fx(xs[len(xs)-1], x), xs[:len(xs)-1])
	}
}

func Foldr1[T any](fx func(T, T) T) func(xs []T) T {
	return func(xs []T) T {
		if len(xs) < 1 {
			return *new(T)
		}
		return Foldr(fx)(xs[len(xs)-1], xs[:len(xs)-1])
	}
}
