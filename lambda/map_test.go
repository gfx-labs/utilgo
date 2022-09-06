package lambda_test

import (
	"reflect"
	"testing"

	"git.tuxpa.in/a/lambda"
)

func TestMapAdd(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Map(arr, func(x int) int { return x + 4 })
	exp := []int{5, 6, 7, 8, 10, 9}
	if !reflect.DeepEqual(ans, exp) {
		t.Fatalf("expect %v to be %v", ans, exp)
	}
}
