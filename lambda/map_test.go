package lambda_test

import (
	"testing"

	"gfx.cafe/util/go/lambda"
	"github.com/stretchr/testify/assert"
)

func TestMapError(t *testing.T) {
	_, err := lambda.MapError(func(int) (int, error) {
		return 0, nil
	})([]int{1, 2, 3})
	assert.NoError(t, err)
}

func TestMapAdd(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	exp := []int{5, 6, 7, 8, 10, 9}
	ans := lambda.Map(func(x int) int { return x + 4 })(arr)
	assert.EqualValues(t, ans, exp)
}

func TestFanAdd(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	exp := []int{5, 6, 7, 8, 10, 9}
	ans := lambda.Fan(func(x int) int { return x + 4 })(arr)
	assert.EqualValues(t, ans, exp)
}
