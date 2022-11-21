package lambda_test

import (
	"testing"

	"gfx.cafe/util/go/lambda"
	"github.com/stretchr/testify/assert"
)

func TestMapAdd(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	exp := []int{5, 6, 7, 8, 10, 9}
	ans := lambda.Map(arr, func(x int) int { return x + 4 })
	assert.EqualValues(t, ans, exp)
}

func TestFanAdd(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	exp := []int{5, 6, 7, 8, 10, 9}
	ans := lambda.Fan(arr, func(x int) int { return x + 4 })
	assert.EqualValues(t, ans, exp)
}
