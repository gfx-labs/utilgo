package lambda

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

// copies slice xs into map in which the key is the slice in dex
func MapSlice[T any](xs []T) map[int]T {
	m := make(map[int]T, len(xs))
	for idx, v := range xs {
		m[idx] = v
	}
	return m
}
