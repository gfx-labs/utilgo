package lambda

import "sync"

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

func FanNil[T any](xs []T, fx func(T)) {
	wg := sync.WaitGroup{}
	wg.Add(len(xs))
	for _, vv := range xs {
		v := vv
		go func() {
			fx(v)
			wg.Done()
		}()
	}
	wg.Wait()
}

func Fan[T any](xs []T, fx func(T) T) []T {
	wg := sync.WaitGroup{}
	wg.Add(len(xs))
	for ii, vv := range xs {
		i, v := ii, vv
		go func() {
			xs[i] = fx(v)
			wg.Done()
		}()
	}
	wg.Wait()
	return xs
}

func FanV[T, V any](xs []T, fx func(T) V) []V {
	wg := sync.WaitGroup{}
	wg.Add(len(xs))
	ov := make([]V, len(xs))
	for ii, vv := range xs {
		i, v := ii, vv
		go func() {
			ov[i] = fx(v)
			wg.Done()
		}()
	}
	wg.Wait()
	return ov
}

func FanError[T any](xs []T, fx func(T) (T, error)) ([]T, []error) {
	wg := sync.WaitGroup{}
	wg.Add(len(xs))
	oe := make([]error, len(xs))
	for ii, vv := range xs {
		i, v := ii, vv
		go func() {
			xs[i], oe[i] = fx(v)
			wg.Done()
		}()
	}
	wg.Wait()
	return xs, oe
}

func FanErrorV[T, V any](xs []T, fx func(T) (V, error)) ([]V, []error) {
	wg := sync.WaitGroup{}
	wg.Add(len(xs))
	ov := make([]V, len(xs))
	oe := make([]error, len(xs))
	for ii, vv := range xs {
		i, v := ii, vv
		go func() {
			ov[i], oe[i] = fx(v)
			wg.Done()
		}()
	}
	wg.Wait()
	return ov, oe
}
