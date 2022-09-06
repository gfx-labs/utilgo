package lambda

import "golang.org/x/exp/constraints"

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
