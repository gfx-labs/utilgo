package lambda

import (
	"reflect"
	"testing"
)

type test_case[T any] struct {
	name   string
	arr    []T
	fn     func(T, T) T
	expect T
}

func TestFoldlInt(t *testing.T) {
	cases := []test_case[int]{
		{
			name:   "sum",
			arr:    []int{1, 2, 3, 4, 5, 6, 7},
			fn:     func(a, b int) int { return a + b },
			expect: 28,
		},
	}

	for _, v := range cases {
		ans := Foldl(v.fn)(0, v.arr)
		if !reflect.DeepEqual(ans, v.expect) {
			t.Errorf("failed test %s, expected %v not %v", v.name, ans, v.expect)
		}
	}
}

func TestFoldlFloat64(t *testing.T) {
	cases := []test_case[float64]{
		{
			name:   "div",
			arr:    []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0},
			fn:     func(a, b float64) float64 { return a / b },
			expect: ((((((float64(1.0) / 2.0) / 3.0) / 4.0) / 5.0) / 6.0) / 7.0),
		},
	}
	for _, v := range cases {
		ans := Foldl1(v.fn)(v.arr)
		if !reflect.DeepEqual(ans, v.expect) {
			t.Errorf("failed test %s, expected %v not %v", v.name, v.expect, ans)
		}
	}
}

func TestFoldr1Float64(t *testing.T) {
	cases := []test_case[float64]{
		{
			name:   "div",
			arr:    []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0},
			fn:     func(a, b float64) float64 { return a / b },
			expect: (float64(1.0) / (2.0 / (3.0 / (4.0 / (5.0 / float64(6.0/7.0)))))),
		},
	}
	for _, v := range cases {
		ans := Foldr1(v.fn)(v.arr)
		if !reflect.DeepEqual(ans, v.expect) {
			t.Errorf("failed test %s, expected %v not %v", v.name, v.expect, ans)
		}
	}
}
