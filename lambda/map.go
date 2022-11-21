package lambda

import "sync"

type ErrorSlice []error

func (e ErrorSlice) First() error {
	for _, v := range e {
		if v != nil {
			return v
		}
	}
	return nil
}
func (e ErrorSlice) IfError() ErrorSlice {
	for _, v := range e {
		if v != nil {
			return e
		}
	}
	return nil
}

func (e ErrorSlice) Error() string {
	for _, v := range e {
		if v != nil {
			return v.Error()
		}
	}
	return "Error() call on Errorslice With No Errors"
}

func MapNil[T any](fx func(T)) func([]T) {
	return func(xs []T) {
		for _, v := range xs {
			fx(v)
		}
	}
}

func Map[T any](fx func(T) T) func([]T) []T {
	return func(xs []T) []T {
		for i, v := range xs {
			xs[i] = fx(v)
		}
		return xs
	}
}

func MapV[T, V any](fx func(T) V) func([]T) []V {
	return func(xs []T) []V {
		ov := make([]V, len(xs))
		for i, v := range xs {
			ov[i] = fx(v)
		}
		return ov
	}
}

func MapError[T any](fx func(T) (T, error)) func(xs []T) ([]T, error) {
	return func(xs []T) ([]T, error) {
		oe := make(ErrorSlice, len(xs))
		for i, v := range xs {
			xs[i], oe[i] = fx(v)
		}
		return xs, oe.IfError()
	}
}

func MapErrorV[T, V any](fx func(T) (V, error)) func(xs []T) ([]V, error) {
	return func(xs []T) ([]V, error) {
		ov := make([]V, len(xs))
		oe := make(ErrorSlice, len(xs))
		for i, v := range xs {
			ov[i], oe[i] = fx(v)
		}
		return ov, oe.IfError()
	}
}

func FanNil[T any](fx func(T)) func(xs []T) {
	return func(xs []T) {
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
}

func Fan[T any](fx func(T) T) func(xs []T) []T {
	return func(xs []T) []T {
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
}

func FanV[T, V any](fx func(T) V) func(xs []T) []V {
	return func(xs []T) []V {
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
}

func FanError[T any](fx func(T) (T, error)) func(xs []T) ([]T, error) {
	return func(xs []T) ([]T, error) {
		wg := sync.WaitGroup{}
		wg.Add(len(xs))
		oe := make(ErrorSlice, len(xs))
		for ii, vv := range xs {
			i, v := ii, vv
			go func() {
				xs[i], oe[i] = fx(v)
				wg.Done()
			}()
		}
		wg.Wait()
		return xs, oe.IfError()
	}
}

func FanErrorV[T, V any](fx func(T) (V, error)) func(xs []T) ([]V, error) {
	return func(xs []T) ([]V, error) {
		wg := sync.WaitGroup{}
		wg.Add(len(xs))
		ov := make([]V, len(xs))
		oe := make(ErrorSlice, len(xs))
		for ii, vv := range xs {
			i, v := ii, vv
			go func() {
				ov[i], oe[i] = fx(v)
				wg.Done()
			}()
		}
		wg.Wait()
		return ov, oe.IfError()
	}
}
