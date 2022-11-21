package lambda

import "math"

// returns []K and []V in separate response args.
func Entries[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	vals := make([]V, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return keys, vals
}

// returns keys of map as a slice
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// returns values of map as a slice
func Values[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// adds all keys in map m2 to m1, then returns m1
func MergeMap[K comparable, V any](m1 map[K]V, m2 map[K]V) (m map[K]V) {
	m = m1
	for k, v := range m2 {
		m1[k] = v
	}
	return
}

// transforms [][]T to []T. it's really just Foldl1(xss, Merge)
func Flatten[T any](xss [][]T) []T {
	return Foldl1(xss, Merge[T])
}

// concats slice a and b and returns such
func Merge[T any](a, b []T) []T {
	return append(a, b...)
}

// concat slice a with each slice in xs, in order.
func MergeN[T any](a []T, xs ...[]T) []T {
	for _, b := range xs {
		a = append(a, b...)
	}
	return a
}

// split xs into n groups
// extra elements will be put in the last array, to ensure that only N groups will ever exist
// that is, if len(xs) == 10 && n == 3, then the outputs will have lengths of [3, 3, 4], respectively
// if len(xs) <= n, result will be [][]T{xs}
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

// segment xs into groups of n, each group being no larger than n
// that is, if len(xs) == 10 && n == 3, then the outputs will have lengths of [3, 3, 3, 1], respectively
func Segment[T any](xs []T, n int) [][]T {
	if n == 0 {
		return nil
	}
	idx := 0
	out := make([][]T, int(math.Ceil(float64(len(xs))/float64(n))))
	for i := 0; i < len(xs); i = i + n {
		from := i
		to := i + n
		if to > len(xs) {
			to = len(xs)
			out[idx] = xs[from:to]
			return out
		}
		out[idx] = xs[from:to]
		idx = idx + 1
	}
	return out
}

// copies slice xs into map in which the key is the slice in dex
func MapSlice[T any](xs []T) map[int]T {
	m := make(map[int]T, len(xs))
	for idx, v := range xs {
		m[idx] = v
	}
	return m
}
