package lambda

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Complex | constraints.Integer | constraints.Float
}

func Add[T Numeric](a, b T) T {
	return a + b
}

func Sub[T Numeric](a, b T) T {
	return a - b
}

func Mul[T Numeric](a, b T) T {
	return a * b
}

func Div[T Numeric](a, b T) T {
	return a / b
}

type Sequencable interface {
	constraints.Integer | constraints.Float
}

func Gte[T Sequencable](a, b T) bool {
	return a >= b
}

func Gt[T Sequencable](a, b T) bool {
	return a > b
}
func Lte[T Sequencable](a, b T) bool {
	return a <= b
}

func Lt[T Sequencable](a, b T) bool {
	return a < b
}

func NumEq[T Sequencable](a, b T) bool {
	return a == b
}

func Seq[T Sequencable](from, to T, by T) []T {
	sz := int(Div(Sub(to, from), by))
	if sz <= 0 {
		return nil
	}
	out := make([]T, 0, sz)
	for i := from; Lte(i, to); i = Add(i, by) {
		out = append(out, i)
	}
	return out
}
